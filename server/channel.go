// Author: xufei
// Date: 2019-09-06 09:57

package server

import (
	"fmt"
	"gim/internal/constant"
	"gim/internal/lg"
	"gim/model"
	"gim/pkg/rpc_service"
	"gim/server/service"
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
			lg.Logger().Info("连接异常", zap.Error(err))

			// 连接异常，用户下线
			user := userSessionMap.getSessionByStream(stream)
			if user.UserID != 0 {
				lg.Logger().Info(user.UserName + "下线成功!")
				c.ctx.server.userOffline(user)
			} else {
				lg.Logger().Warn("未查询到用户信息!")
			}

			return err
		}

		res := c.channelHandler(stream, request)

		err = stream.Send(res)
		if err != nil {
			lg.Logger().Error("响应失败", zap.Error(err))
		}
	}

	return nil
}

func (c *channelService) channelHandler(stream rpc_service.GIMService_ChannelServer, req *rpc_service.GIMRequest) *rpc_service.GIMResponse {
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

		if c.ctx.server.accountSrv.SaveAndCheckLogin(req.RequestID) {
			res.ResMsg = "Logined"

			return res
		}

		// 存储用户的连接，用户信息
		c.ctx.server.accountSrv.SaveSession(req.RequestID, req.ReqMsg)
		userSessionMap.put(req.RequestID, stream)

		// 保存用户登录的服务器信息和订阅名称
		channelInfo := model.UserChannelInfo{
			UserID:      req.RequestID,
			ChannelName: fmt.Sprintf("%s-%d", GetConfig().RpcPort, req.RequestID),
		}
		_ = c.ctx.server.accountSrv.StoreServerChannelInfo(req.RequestID, channelInfo)

		pb := service.NewPubSubRedisService(c.ctx.server.redisClient)
		userSessionMap.putSubscribe(channelInfo.UserID, pb)

		// 用户消息接收订阅
		c.ctx.server.SubscribeMessage(pb, channelInfo)

		lg.Logger().Info(req.ReqMsg + " 用户登录成功")
	}

	if req.MsgType == constant.PingMsg {
		res.ResMsg = "心跳信息: Pong"

		// 心跳消息处理
		lg.Logger().Info("心跳信息: " + req.ReqMsg)
	}

	return res
}
