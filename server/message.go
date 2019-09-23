// Author: xufei
// Date: 2019-09-06 14:29

package server

import (
	"encoding/json"
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
		lg.Logger().Error("消息发送失败", zap.Error(err))
		return err
	}

	return nil
}

func (s *Server) sendP2PMsg(msg model.P2PReq) error {
	channelInfo := s.userCache.ServerChannelInfo(msg.ReceiverIDToString())

	user := userSessionMap.getSessionByUserID(msg.UserID)
	pushMsg := model.PushMsg{
		UserID:  msg.ReceiverID,
		Msg:     user.FormatMsg(msg.Msg),
		MsgType: constant.ChatMsg,
	}
	msgBody, _ := json.Marshal(pushMsg)

	return s.Publish(channelInfo.ChannelName, string(msgBody))
}

func (s *Server) sendGroupMsg(msg model.MsgReq) error {
	return s.PublishGroup(msg)
}
