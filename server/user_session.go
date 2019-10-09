// Author: xufei
// Date: 2019-09-06 11:31

package server

import (
	"gim/model"
	"gim/pkg/rpc_service"
	"gim/server/service"
	"sync"
)

var userSessionMap *userSession

type userSession struct {
	ctx        *context
	streams    sync.Map
	subscribes sync.Map
}

func newUserSession(ctx *context) *userSession {
	return &userSession{
		ctx: ctx,
	}
}

// Sessions

func (s *userSession) getSessionByStream(stream rpc_service.GIMService_ChannelServer) model.User {
	var user model.User
	s.rangStreams(func(key, value interface{}) bool {
		if stream == value.(rpc_service.GIMService_ChannelServer) {
			user = s.ctx.server.accountSrv.GetSessionByUserID(key.(int64))
			return true
		} else {
			return false
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

func (s *userSession) rangStreams(f func(key, value interface{}) bool) {
	s.streams.Range(func(key, value interface{}) bool {
		return f(key, value)
	})
}

func (s *userSession) putSubscribe(userID int64, sub service.PubSub) {
	s.subscribes.Store(userID, sub)
}

func (s *userSession) getSubscribe(userID int64) service.PubSub {
	sub, ok := s.subscribes.Load(userID)
	if !ok {
		return nil
	}
	return sub.(service.PubSub)
}

func (s *userSession) removeSubscribe(userID int64) {
	s.subscribes.Delete(userID)
}
