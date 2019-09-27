// Author: xufei
// Date: 2019-09-05 11:30

package server

import (
	"gim/model"
	"time"
)

// 注册用户，存入 redis
func (s *Server) accountRegister(user model.User) (model.User, error) {
	// 生成用户 ID，目前使用时间戳，但是当高并发用户同时注册时，可能会有冲突，生成出相同值
	user.UserID = time.Now().UnixNano()
	return s.accountSrv.Register(user)
}

func (s *Server) userOffline(user model.User) {
	userSessionMap.ctx.server.accountSrv.RemoveSession(user.UserID)
	userSessionMap.remove(user.UserID)
}

func (s *Server) getOnlineUsers() []model.User {
	return s.accountSrv.GetAllOnlineUsers()
}
