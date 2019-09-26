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

func (s *Server) Publish(channelName, msg string) error {
	return s.redisClient.Publish(channelName, msg).Err()
}

func (s *Server) PublishGroup(msg model.MsgReq) error {
	var errs error

	channels := s.userCache.GetAllServerChannelInfo()
	user := userSessionMap.getSessionByUserID(msg.UserID)
	formatMsg := user.FormatMsg(msg.Msg)

	for _, channel := range channels {
		if msg.UserID == channel.UserID {
			continue
		}

		pushMsg := model.PushMsg{
			UserID:  channel.UserID,
			Msg:     formatMsg,
			MsgType: constant.ChatMsg,
		}
		msgBody, _ := json.Marshal(pushMsg)

		err := s.Publish(channel.ChannelName, string(msgBody))
		if err != nil {
			errs = errors.Wrap(errs, err.Error())
		}
	}
	if errs != nil {
		errs = errors.WithStack(errs)
	}

	return errs
}

func (s *Server) SubscribeMessageByUser(channelInfo model.UserChannelInfo) {
	subscribe := s.redisClient.Subscribe(channelInfo.ChannelName)

	go func() {
		// go chan 缓存是 100，当消息超过 100 时 30 秒后消息会丢失
		for ch := range subscribe.Channel() {
			go func(payLoad string) {
				lg.Logger().Info("接收到消息处理中...")

				payload := payLoad
				var msg model.PushMsg
				_ = json.Unmarshal([]byte(payload), &msg)

				stream := userSessionMap.get(msg.UserID)
				if err := s.pushMsg(stream, msg); err != nil {
					lg.Logger().Error("发送消息失败", zap.Error(err))
				}
			}(ch.Payload)
		}
	}()
}
