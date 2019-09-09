// Author: xufei
// Date: 2019-09-06 10:23

package client

import (
	context2 "context"
	"fmt"
	"gim/internal/constant"
	"gim/internal/lg"
	"gim/internal/rpc_service"
	"gim/model"
	"io"
	"time"

	"go.uber.org/zap"
)

type userClient struct {
	ctx      *context
	userInfo model.User

	channel   rpc_service.GIMService_ChannelClient
	sendCh    chan *rpc_service.GIMRequest
	receiveCh chan *rpc_service.GIMResponse
}

func newUserClient(ctx *context) (uc *userClient, err error) {
	uc = &userClient{
		sendCh:    make(chan *rpc_service.GIMRequest),
		receiveCh: make(chan *rpc_service.GIMResponse),
		ctx:       ctx,
	}

	uc.userInfo = model.User{
		UserID:   GetConfig().UserID,
		UserName: GetConfig().Username,
	}

	uc.channel, err = rpc_service.NewGIMServiceClient(ctx.client.rpc.conn).
		Channel(context2.Background())
	if err != nil {
		return
	}

	uc.ctx.client.waitGroup.Wrap(func() {
		uc.recvStream()
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
			// TODO: writer log
			lg.Logger().Info(fmt.Sprintf("接收到的消息: [%d]: %s", res.ResponseID, res.ResMsg))
		case <-heartbeatTime.C:
			err := c.sendHeartbeat()
			if err != nil {
				lg.Logger().Error("心跳发送失败", zap.Error(err))
			}
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

func (c *userClient) recvStream() {
	for {
		res, err := c.channel.Recv()
		if err != nil {
			if err == io.EOF {
				lg.Logger().Error("无消息可以接收")
				return
			}

			// TODO: 连接失败，应该关闭用户客户端
			lg.Logger().Error("消息接收失败", zap.Error(err))
			return
		}

		c.receiveCh <- res
	}
}
