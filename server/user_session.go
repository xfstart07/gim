// Author: xufei
// Date: 2019-09-06 11:31

package server

import (
	"fmt"
	"gim/internal/constant"
	"gim/model"
	"gim/pkg/rpc_service"
	"sync"

	"github.com/go-redis/redis"
)

var userSessionMap *userSession

type userSession struct {
	client   *redis.Client
	sessions sync.Map
	streams  sync.Map
}

func newUserSession(cli *redis.Client) *userSession {
	return &userSession{
		client: cli,
	}
}

// Sessions

func (s *userSession) saveSession(userID int64, userName string) {
	s.sessions.Store(userID, userName)
	key := fmt.Sprintf("%s%d", constant.ServerSessionPrefixName, userID)
	s.client.Set(key, userName, -1)
}

func (s *userSession) removeSession(userID int64) {
	s.sessions.Delete(userID)
	key := fmt.Sprintf("%s%d", constant.ServerSessionPrefixName, userID)
	s.client.Del(key)
}

func (s *userSession) getSessionByUserID(userID int64) model.User {
	name, ok := s.sessions.Load(userID)
	if !ok {
		key := fmt.Sprintf("%s%d", constant.ServerSessionPrefixName, userID)
		name = s.client.Get(key).Val()
	}
	return model.User{
		UserID:   userID,
		UserName: name.(string),
	}
}

func (s *userSession) getSessionByStream(stream rpc_service.GIMService_ChannelServer) model.User {
	user := model.User{}
	s.rangStreams(func(key, value interface{}) {
		if stream == value.(rpc_service.GIMService_ChannelServer) {
			name, ok := s.sessions.Load(key)
			if !ok {
				key := fmt.Sprintf("%s%v", constant.ServerSessionPrefixName, key)
				name = s.client.Get(key).Val()
			}

			user.UserID = key.(int64)
			user.UserName = name.(string)
		}
	})

	return user
}

// Streams

func (s *userSession) put(userID int64, stream rpc_service.GIMService_ChannelServer) {
	s.streams.Store(userID, stream)
}

func (s *userSession) get(userID int64) rpc_service.GIMService_ChannelServer {
	stream, ok := s.streams.Load(userID)
	if !ok {
		return nil
	}
	return stream.(rpc_service.GIMService_ChannelServer)
}

func (s *userSession) remove(userID int64) {
	s.streams.Delete(userID)
}

func (s *userSession) rangStreams(f func(key, value interface{})) {
	s.streams.Range(func(key, value interface{}) bool {
		f(key, value)
		return true
	})
}
