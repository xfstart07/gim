// Author: xufei
// Date: 2019-09-06 10:23

package client

import (
	context2 "context"
	"fmt"
	"gim/client/handler"
	"gim/client/service"
	"gim/internal/ciface"
	"gim/internal/constant"
	"gim/internal/lg"
	"gim/model"
	"gim/pkg/rpc_service"
	"io"
	"time"

	"google.golang.org/grpc"

	"go.uber.org/zap"
)

// 用户客户端
type userClient struct {
	ctx      *context
	config   *model.ClientConfig
	userInfo model.User
	isExit   bool // 标记客户端是否已经退出，避免多次退出
	errCount int

	msgHander ciface.MessageHandler
	msgWriter service.MsgLogger

	rpc       *rpcServer
	channel   rpc_service.GIMService_ChannelClient
	sendCh    chan *rpc_service.GIMRequest
	receiveCh chan *rpc_service.GIMResponse
	closeCh   chan int
}

func newUserClient(ctx *context, cfg *model.ClientConfig) (uc *userClient) {
	uc = &userClient{
		ctx:       ctx,
		config:    cfg,
		msgWriter: service.NewWriter(cfg),
		userInfo:  cfg.User,
	}
	uc.msgHander = handler.NewMessageHandler(uc, GetConfig())

	return
}

func (c *userClient) Start() (err error) {
	c.reset()

	rpcDiaUrl := fmt.Sprintf("%s://authority/%s", c.ctx.client.etcdResolver.Scheme(), c.config.EtcdServerName)
	ctx, cancel := context2.WithTimeout(context2.Background(), 10*time.Second)
	defer cancel()

	// https://github.com/grpc/grpc/blob/master/doc/service_config.md
	// conn, err := grpc.DialContext(ctx, rpcDiaUrl, grpc.WithBalancerName(roundrobin.Name), grpc.WithInsecure())
	conn, err := grpc.DialContext(ctx, rpcDiaUrl, grpc.WithInsecure(), grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`))
	if err != nil {
		panic(err)
	}
	c.rpc = newRpcServer(conn)

	c.channel, err = rpc_service.NewGIMServiceClient(c.rpc.conn).Channel(context2.Background())
	if err != nil {
		return
	}

	c.ctx.client.waitGroup.Wrap(func() {
		c.recvPump()
	})
	c.ctx.client.waitGroup.Wrap(func() {
		c.dispatch()
	})

	c.ctx.client.waitGroup.Wrap(func() {
		c.Login()
	})

	return nil
}

func (c *userClient) reset() {
	c.sendCh = make(chan *rpc_service.GIMRequest)
	c.receiveCh = make(chan *rpc_service.GIMResponse)
	c.closeCh = make(chan int)
}

func (c *userClient) reconnect() error {
	c.errCount++

	return c.Start()
}

func (c *userClient) dispatch() {
	heartbeatTime := time.NewTicker(time.Duration(GetConfig().HeartbeatTime) * time.Second)
	for {
		select {
		case req := <-c.sendCh:
			err := c.sendChannel(req)
			if err != nil {
				lg.Logger().Error("发送消息失败", zap.Error(err), zap.Any("req", req))
			}
		case res := <-c.receiveCh:
			// TODO: 判断登录状态
			c.writerMsg(res.ResponseID, res.ResMsg, res.MsgType)
		case <-heartbeatTime.C:
			err := c.sendHeartbeat()
			if err != nil {
				lg.Logger().Error("心跳发送失败, 服务端无法连接", zap.Error(err))

				// reconnect
				for c.errCount < c.config.ReconnectCount {
					if err := c.reconnect(); err != nil {
						lg.Logger().Error("用户重连失败", zap.Error(err))
					} else {
						lg.Logger().Info("用户重连成功")
						return
					}
				}

				return
			}
		case <-c.closeCh:
			for c.errCount < c.config.ReconnectCount {
				if err := c.reconnect(); err != nil {
					lg.Logger().Error("用户重连失败", zap.Error(err))
				} else {
					lg.Logger().Info("用户重连成功")
					return
				}
			}
			return
		}
	}
}

func (c *userClient) Shutdown() {
	if c.isExit {
		return
	}

	c.isExit = true

	_ = c.channel.CloseSend()
	_ = c.rpc.conn.Close()

	if err := c.msgHander.OfflineUser(model.MsgReq{UserID: c.userInfo.UserID}); err != nil {
		lg.Logger().Error("发送离线消息失败", zap.Error(err))
	}
	c.msgWriter.Close()

	close(c.closeCh)

	lg.Logger().Fatal("用户下线！")
}

func (c *userClient) Login() {
	lg.Logger().Debug("用户: " + c.userInfo.UserName + " 登录中...")
	c.sendMsg(c.userInfo.UserName, constant.LoginMsg)
}

func (c *userClient) sendMsg(msg string, msgType int) {
	req := &rpc_service.GIMRequest{
		RequestID: c.userInfo.UserID,
		ReqMsg:    msg,
		MsgType:   int32(msgType),
	}

	c.sendCh <- req
}

func (c *userClient) sendHeartbeat() error {
	req := &rpc_service.GIMRequest{
		RequestID: c.userInfo.UserID,
		ReqMsg:    "Ping",
		MsgType:   constant.PingMsg,
	}

	return c.sendChannel(req)
}

func (c *userClient) sendChannel(req *rpc_service.GIMRequest) error {
	return c.channel.Send(req)
}

func (c *userClient) recvPump() {
	for {
		res, err := c.channel.Recv()
		if err != nil {
			if err == io.EOF {
				return
			}

			if err == grpc.ErrClientConnClosing {
				return
			}

			lg.Logger().Error("消息接收失败", zap.Error(err))
			c.Shutdown()
			return
		}

		c.receiveCh <- res
	}
}

func (c *userClient) writerMsg(userID int64, msg string, msgType int32) {
	lg.Logger().Info(msg)

	if msgType == constant.LoginMsg {
		if msg != "OK" {
			// Shutdown, because user logined
			c.Shutdown()
		}
	}

	// asyncWriter log
	if msgType == constant.ChatMsg {
		c.msgWriter.Log(msg)
	}
}
