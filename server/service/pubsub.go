// Author: xufei
// Date: 2019-10-08 15:49

package service

import (
	"encoding/json"
	"gim/internal/lg"
	"gim/model"

	"github.com/go-redis/redis"
)

type SubscribeFunc func(payload string)

type PubSub interface {
	Publish(channelName string, msg model.PushMsg) error
	Subscribe(channelInfo model.UserChannelInfo, callback SubscribeFunc)
	UnSubscribe(channelInfo model.UserChannelInfo) error
}

type pubSubRedisService struct {
	client    *redis.Client
	subscribe *redis.PubSub
}

func (s *pubSubRedisService) Publish(channelName string, msg model.PushMsg) error {
	msgBody, _ := json.Marshal(msg)

	lg.Logger().Sugar().Debug("推送消息", msg)
	return s.client.Publish(channelName, string(msgBody)).Err()
}

func (s *pubSubRedisService) Subscribe(channelInfo model.UserChannelInfo, callback SubscribeFunc) {
	s.subscribe = s.client.Subscribe(channelInfo.ChannelName)

	go func() {
		// go chan 缓存是 100，当消息超过 100 时 30 秒后消息会丢失
		for ch := range s.subscribe.Channel() {
			go func(payLoad string) {
				callback(payLoad)
			}(ch.Payload)
		}
	}()
}

func (s *pubSubRedisService) UnSubscribe(channelInfo model.UserChannelInfo) error {
	return s.subscribe.Unsubscribe(channelInfo.ChannelName)
}

func NewPubSubRedisService(cli *redis.Client) *pubSubRedisService {
	return &pubSubRedisService{
		client: cli,
	}
}
