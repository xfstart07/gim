// Author: xufei
// Date: 2019-09-06 09:57

package server

import (
	"gim/internal/constant"
	"gim/internal/lg"
	"gim/pkg/rpc_service"
	"io"

	"go.uber.org/zap"
)

// rpc 注册，消息处理
type channelService struct {
	ctx *context
}

func NewChannelService(ctx *context) *channelService {
	return &channelService{
		ctx: ctx,
	}
}

func (c *channelService) Channel(stream rpc_service.GIMService_ChannelServer) error {
	for {
		request, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				break
			}

			// 连接异常，用户下线
			user := userSessionMap.getSessionByStream(stream)
			c.ctx.server.userOffline(user)

			lg.Logger().Info("连接异常", zap.Error(err))
			return err
		}

		res := channelHandler(stream, request)

		err = stream.Send(res)
		if err != nil {
			lg.Logger().Error("响应失败", zap.Error(err))
		}
	}

	return nil
}

func channelHandler(stream rpc_service.GIMService_ChannelServer, req *rpc_service.GIMRequest) *rpc_service.GIMResponse {
	if req == nil {
		lg.Logger().Info("无消息", zap.Any("req", req))
		return nil
	}

	res := &rpc_service.GIMResponse{
		ResponseID: req.RequestID,
		MsgType:    req.MsgType,
	}

	if req.MsgType == constant.LoginMsg {
		// 登录处理
		res.ResMsg = "OK"

		// 存储用户的连接，用户信息
		userSessionMap.saveSession(req.RequestID, req.ReqMsg)
		userSessionMap.put(req.RequestID, stream)

		lg.Logger().Info(req.ReqMsg + " 用户登录成功")
	}

	if req.MsgType == constant.PingMsg {
		res.ResMsg = "心跳信息: Pong"

		// 心跳消息处理
		lg.Logger().Info("心跳信息: " + req.ReqMsg)
	}

	return res
}
