// Author: xufei
// Date: 2019-09-17 16:48

package service

import (
	"encoding/json"
	"gim/internal/constant"
	"gim/model"

	"github.com/go-redis/redis"
)

type UserCache interface {
	StoreServerChannelInfo(userID string, channelInfo model.UserChannelInfo) error
	ServerChannelInfo(userID string) model.UserChannelInfo
	GetAllServerChannelInfo() []model.UserChannelInfo
}

type userCacheService struct {
	client *redis.Client
}

func (c *userCacheService) GetAllServerChannelInfo() []model.UserChannelInfo {
	keys := c.client.Keys(constant.ServerChannelPrefixName + "*").Val()

	values := c.client.MGet(keys...).Val()

	channels := make([]model.UserChannelInfo, 0, len(values))
	for _, value := range values {
		var info model.UserChannelInfo
		b := value.(string)
		_ = json.Unmarshal([]byte(b), &info)

		channels = append(channels, info)
	}

	return channels
}

func (c *userCacheService) StoreServerChannelInfo(userID string, channelInfo model.UserChannelInfo) error {
	infoByte, _ := json.Marshal(channelInfo)
	return c.client.Set(constant.ServerChannelPrefixName+userID, string(infoByte), -1).Err()
}

func (c *userCacheService) ServerChannelInfo(userID string) model.UserChannelInfo {
	var channelInfo model.UserChannelInfo
	infoByte := c.client.Get(constant.ServerChannelPrefixName + userID).Val()

	_ = json.Unmarshal([]byte(infoByte), &channelInfo)

	return channelInfo
}

func NewUserCache(client *redis.Client) *userCacheService {
	return &userCacheService{
		client: client,
	}
}
