// Author: xufei
// Date: 2019-09-06 11:31

package server

import (
	"gim/internal/rpc_service"
	"gim/model"
	"sync"
)

var userSessionMap = newUserSession()

type userSession struct {
	sessions sync.Map
	streams  sync.Map
}

func newUserSession() *userSession {
	return &userSession{
		streams: sync.Map{},
	}
}

// Sessions

func (s *userSession) saveSession(userID int64, userName string) {
	s.sessions.Store(userID, userName)
}

func (s *userSession) removeSession(userID int64) {
	s.sessions.Delete(userID)
}

func (s *userSession) getSession(stream rpc_service.GIMService_ChannelServer) model.User {
	user := model.User{}
	s.rangStreams(func(key, value interface{}) {
		if stream == value.(rpc_service.GIMService_ChannelServer) {
			name, _ := s.sessions.Load(key)
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
	stream, _ := s.streams.Load(userID)
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
