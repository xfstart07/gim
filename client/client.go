// Author: xufei
// Date: 2019-09-06 10:23

package client

import (
	context2 "context"
	"gim/internal/constant"
	"gim/internal/lg"
	"gim/model"
	"gim/pkg/rpc_service"
	"io"
	"time"

	"go.uber.org/zap"
)

// 用户客户端
type userClient struct {
	ctx      *context
	config   *model.ClientConfig
	userInfo model.User
	errCount int

	msgLog MsgLogger

	channel   rpc_service.GIMService_ChannelClient
	sendCh    chan *rpc_service.GIMRequest
	receiveCh chan *rpc_service.GIMResponse
	closeCh   chan int
}

func newUserClient(ctx *context, cfg *model.ClientConfig) (uc *userClient, err error) {
	uc = &userClient{
		ctx:       ctx,
		config:    cfg,
		msgLog:    NewWriter(ctx, cfg),
		sendCh:    make(chan *rpc_service.GIMRequest),
		receiveCh: make(chan *rpc_service.GIMResponse),
		closeCh:   make(chan int),
	}

	uc.userInfo = model.User{
		UserID:   GetConfig().UserID,
		UserName: GetConfig().Username,
	}

	uc.channel, err = rpc_service.NewGIMServiceClient(ctx.client.rpc.conn).Channel(context2.Background())
	if err != nil {
		return
	}

	uc.ctx.client.waitGroup.Wrap(func() {
		uc.recvPump()
	})
	uc.ctx.client.waitGroup.Wrap(func() {
		uc.dispatch()
	})

	return
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
			c.writerMsg(res.ResponseID, res.ResMsg, res.MsgType)
		case <-heartbeatTime.C:
			err := c.sendHeartbeat()
			if err != nil {
				lg.Logger().Error("心跳发送失败, 服务端无法连接", zap.Error(err))
				// reconnect
				c.errCount++
				if c.errCount >= c.config.ReconnectCount {
					return
				}
			}
		case <-c.closeCh:
			// receive user client shutdown signal
			lg.Logger().Info("用户下线！")
			return
		}
	}
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

			c.shutdown()
			lg.Logger().Error("消息接收失败", zap.Error(err))
			return
		}

		c.receiveCh <- res
	}
}

func (c *userClient) shutdown() {
	close(c.closeCh)
}

func (c *userClient) writerMsg(userID int64, msg string, msgType int32) {
	lg.Logger().Info(msg)

	// asyncWriter log
	if msgType == constant.ChatMsg {
		c.msgLog.Log(msg)
	}
}
