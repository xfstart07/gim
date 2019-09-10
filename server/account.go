// Author: xufei
// Date: 2019-09-05 11:30

package server

import (
	"gim/internal/lg"
	"gim/model"
	"gim/server/service"
	"time"
)

// 注册用户，存入 redis
func (s *Server) accountRegister(user model.User) (model.User, error) {
	srv := service.NewAccountService(s.redisClient)

	// TODO: 生成用户 ID，目前使用 时间戳，但是当用户同时注册时，可能会有冲突，生成出相同值
	user.UserID = time.Now().UnixNano()
	return srv.Register(user)
}

func (s *Server) userOffline(user model.User) {
	userSessionMap.remove(user.UserID)
	lg.Logger().Info(user.UserName + "下线成功!")
}
