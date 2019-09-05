// Author: xufei
// Date: 2019-09-04 17:23

package server

import (
	"gim/internal/lg"
	"gim/internal/util"
)

type Server struct {
	waitGroup util.WaitGroupWrapper
}

func New() *Server {
	return &Server{}
}

func (s *Server) Main() {
	if err := InitConfig(); err != nil {
		panic(err)
	}

	// 设置日志的级别
	lg.SetLevel(GetConfig().LogLevel)

	server := newHTTPServer()
	s.waitGroup.Wrap(func() {
		server.Run()
	})

	s.waitGroup.Wait()
	lg.Logger().Info("Client: done!")
}
