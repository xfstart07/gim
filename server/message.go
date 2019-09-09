// Author: xufei
// Date: 2019-09-06 14:29

package server

import (
	"fmt"
	"gim/internal/lg"
	"gim/internal/rpc_service"
	"gim/model"

	"github.com/pkg/errors"
	"go.uber.org/zap"
)

var errMessageSendFailed = errors.New("message send failed")

func (s *Server) sendMsg(msg model.MsgReq) error {
	stream := userSessionMap.get(msg.UserID)

	res := &rpc_service.GIMResponse{
		ResponseID: msg.UserID,
		ResMsg:     msg.Msg,
	}

	err := stream.Send(res)
	if err != nil {
		lg.Logger().Error("消息发送失败", zap.Error(err))
		return errMessageSendFailed
	}

	return nil
}

func (s *Server) sendP2PMsg(msg model.P2PReq) error {
	stream := userSessionMap.get(msg.ReceiverID)

	res := &rpc_service.GIMResponse{
		ResponseID: msg.UserID,
		ResMsg:     msg.Msg,
	}

	err := stream.Send(res)
	if err != nil {
		lg.Logger().Error("消息发送失败", zap.Error(err))
		return errMessageSendFailed
	}

	return nil
}

func (s *Server) sendGroupMsg(msg model.MsgReq) error {
	var errs error

	userSessionMap.rangStreams(func(key, value interface{}) {
		userID := key.(int64)
		if userID == msg.UserID {
			// 自己不需要发送消息
			return
		}

		stream := value.(rpc_service.GIMService_ChannelServer)
		res := &rpc_service.GIMResponse{
			ResponseID: msg.UserID,
			ResMsg:     msg.Msg,
		}

		err := stream.Send(res)
		if err != nil {
			lg.Logger().Error("消息发送失败", zap.Error(err))
			errs = errors.Wrap(err, fmt.Sprintf("%d 消息发送失败", userID))
		}
	})

	return errs
}
