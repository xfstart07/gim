// Author: xufei
// Date: 2019-09-04

package client

import (
	"fmt"
	"gim/internal/http_helper"
	"gim/internal/lg"

	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
)

type httpServer struct {
	router *gin.Engine
}

func newHTTPServer() *httpServer {
	server := &httpServer{
		router: gin.Default(),
	}

	server.setRouter()

	return server
}

// 注册接口
func (s *httpServer) setRouter() {
	s.router.POST("/sendMsg", sendMsg)
}

func (s *httpServer) Run() {
	err := s.router.Run(fmt.Sprintf(":%s", GetConfig().WebPort))
	if err != nil {
		lg.Logger().Error("API 服务启动失败", zap.Error(err))
	}
}

// example: curl -X POST --header 'Content-Type: application/json' -d '{"user_id": 1567750270024892000, "msg": "你好"}' http://localhost:8082/sendMsg
func sendMsg(ctx *gin.Context) {
	// TODO: send message，to server

	http_helper.RenderCreated(ctx, nil)
	lg.Logger().Info("成功")
}
