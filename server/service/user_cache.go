// Author: xufei
// Date: 2019-09-17 16:48

package service

import (
	"encoding/json"
	"fmt"
	"gim/internal/constant"
	"gim/model"
	"strconv"
)

func (s *accountService) GetAllOnlineUsers() []model.User {
	keys := s.store.SMembers(constant.LoginStatusSetKey).Val()

	users := make([]model.User, 0, len(keys))
	for _, value := range keys {
		userID, _ := strconv.ParseInt(value, 10, 64)

		user := s.GetSessionByUserID(userID)
		users = append(users, user)
	}

	return users
}

func (s *accountService) GetAllServerChannelInfo() []model.UserChannelInfo {
	keys := s.store.Keys(constant.ServerChannelPrefixName + "*").Val()

	values := s.store.MGet(keys...).Val()

	channels := make([]model.UserChannelInfo, 0, len(values))
	for _, value := range values {
		var info model.UserChannelInfo
		b := value.(string)
		_ = json.Unmarshal([]byte(b), &info)

		channels = append(channels, info)
	}

	return channels
}

func (s *accountService) StoreServerChannelInfo(userID int64, channelInfo model.UserChannelInfo) error {
	infoByte, _ := json.Marshal(channelInfo)
	return s.store.Set(channelInfoKey(userID), string(infoByte), -1).Err()
}

func (s *accountService) ServerChannelInfo(userID int64) model.UserChannelInfo {
	var channelInfo model.UserChannelInfo
	infoByte := s.store.Get(channelInfoKey(userID)).Val()

	_ = json.Unmarshal([]byte(infoByte), &channelInfo)

	return channelInfo
}

func channelInfoKey(userID int64) string {
	return fmt.Sprintf("%s%d", constant.ServerChannelPrefixName, userID)
}
