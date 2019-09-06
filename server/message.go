// Author: xufei
// Date: 2019-09-06 14:29

package server

import (
	"errors"
	"gim/internal/lg"
	"gim/internal/rpc_service"
	"gim/model"

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
