// Author: xufei
// Date: 2019-09-27 15:03

package service

import (
	"fmt"
	"gim/internal/constant"
	"gim/model"
)

func (s *accountService) SaveSession(userID int64, userName string) {
	s.sessions.Store(userID, userName)
	key := sessionUserKey(userID)
	s.client.Set(key, userName, -1)
}

func (s *accountService) GetSessionByUserID(userID int64) model.User {
	name, ok := s.sessions.Load(userID)

	if !ok {
		key := sessionUserKey(userID)
		name = s.client.Get(key).Val()
	}

	return model.User{
		UserID:   userID,
		UserName: name.(string),
	}
}

func (s *accountService) RemoveSession(userID int64) {
	s.sessions.Delete(userID)
	key := sessionUserKey(userID)
	s.client.Del(key)
}

func (s *accountService) SaveAndCheckLogin(userID int64) bool {
	add := s.client.SAdd(constant.LoginStatusSetKey, userID).Val()

	// 0 表示已存在
	return add == 0
}

func sessionUserKey(userID int64) string {
	return fmt.Sprintf("%s%d", constant.ServerSessionPrefixName, userID)
}
