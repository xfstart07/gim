// Author: xufei
// Date: 2019-09-06 14:29

package server

import (
	"gim/internal/constant"
	"gim/model"
	"gim/pkg/rpc_service"

	"github.com/pkg/errors"
)

var errMessageSendFailed = errors.New("message send failed")

func (s *Server) sendMsg(msg model.MsgReq) error {
	stream := userSessionMap.get(msg.UserID)

	return s.pushMsg(stream, model.PushMsg{
		UserID: msg.UserID,
		Msg:    msg.Msg,
	})
}

func (s *Server) pushMsg(stream rpc_service.GIMService_ChannelServer, msg model.PushMsg) error {
	res := &rpc_service.GIMResponse{
		ResponseID: msg.UserID,
		ResMsg:     msg.Msg,
		MsgType:    constant.ChatMsg,
	}

	err := stream.Send(res)
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (s *Server) sendP2PMsg(msg model.P2PReq) error {
	channelInfo := s.accountSrv.ServerChannelInfo(msg.ReceiverID)

	user := s.accountSrv.GetSessionByUserID(msg.UserID)
	err := s.PublishMessage(channelInfo.ChannelName, model.PushMsg{
		UserID:  msg.ReceiverID,
		Msg:     user.FormatMsg(msg.Msg),
		MsgType: constant.ChatMsg,
	})
	if err != nil {
		err = errors.WithStack(err)
	}

	return err
}

func (s *Server) sendGroupMsg(msg model.MsgReq) error {
	return s.PublishMessage("", model.PushMsg{
		UserID: msg.UserID,
		Msg:    msg.Msg,
	})
}
