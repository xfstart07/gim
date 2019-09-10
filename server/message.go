// Author: xufei
// Date: 2019-09-06 14:29

package server

import (
	"fmt"
	"gim/internal/constant"
	"gim/internal/lg"
	"gim/model"
	"gim/pkg/rpc_service"

	"github.com/pkg/errors"
	"go.uber.org/zap"
)

var errMessageSendFailed = errors.New("message send failed")

func (s *Server) sendMsg(msg model.MsgReq) error {
	stream := userSessionMap.get(msg.UserID)

	return s.pushMsg(stream, model.PushMsg{
		UserID:  msg.UserID,
		MsgType: constant.ChatMsg,
		Msg:     msg.Msg,
	})
}

func (s *Server) pushMsg(stream rpc_service.GIMService_ChannelServer, msg model.PushMsg) error {
	user := userSessionMap.getSessionByUserID(msg.UserID)

	res := &rpc_service.GIMResponse{
		ResponseID: msg.UserID,
		ResMsg:     fmt.Sprintf("%s: 【%s】", user.UserName, msg.Msg),
		MsgType:    msg.MsgType,
	}

	err := stream.Send(res)
	if err != nil {
		lg.Logger().Error("消息发送失败", zap.Error(err))
		return err
	}

	return nil
}

func (s *Server) sendP2PMsg(msg model.P2PReq) error {
	stream := userSessionMap.get(msg.ReceiverID)

	return s.pushMsg(stream, model.PushMsg{
		UserID:  msg.UserID,
		Msg:     msg.Msg,
		MsgType: constant.ChatMsg,
	})
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
		err := s.pushMsg(stream, model.PushMsg{
			UserID:  msg.UserID,
			Msg:     msg.Msg,
			MsgType: constant.ChatMsg,
		})
		if err != nil {
			errs = errors.Wrap(err, fmt.Sprintf("%d 消息发送失败", msg.UserID))
		}
	})

	return errs
}
