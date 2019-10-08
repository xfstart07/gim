// Author: xufei
// Date: 2019-09-17 17:04

package server

import (
	"encoding/json"
	"gim/internal/constant"
	"gim/internal/lg"
	"gim/model"

	"github.com/pkg/errors"

	"go.uber.org/zap"
)

func (s *Server) PublishMessage(channelName string, msg model.PushMsg) error {
	var errs error

	if channelName != "" {
		// a single message
		return s.pubsubSrv.Publish(channelName, msg)
	}

	// group send message

	channels := s.accountSrv.GetAllServerChannelInfo()
	formatMsg := s.accountSrv.GetSessionByUserID(msg.UserID).FormatMsg(msg.Msg)

	for _, channel := range channels {
		if msg.UserID == channel.UserID {
			continue
		}

		err := s.pubsubSrv.Publish(channel.ChannelName, model.PushMsg{
			UserID:  channel.UserID,
			Msg:     formatMsg,
			MsgType: constant.ChatMsg,
		})
		if err != nil {
			errs = errors.Wrap(errs, err.Error())
		}
	}
	if errs != nil {
		errs = errors.WithStack(errs)
	}

	return errs
}

func (s *Server) SubscribeMessage(channelInfo model.UserChannelInfo) {
	s.pubsubSrv.Subscribe(channelInfo, func(payload string) {
		lg.Logger().Info("接收到消息处理中...")

		var msg model.PushMsg
		_ = json.Unmarshal([]byte(payload), &msg)

		stream := userSessionMap.get(msg.UserID)
		if err := s.pushMsg(stream, msg); err != nil {
			lg.Logger().Error("发送消息失败", zap.Error(err))
		}
	})
}
