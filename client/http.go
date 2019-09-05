// Author: xufei
// Date: 2019-09-04

package client

import (
	"fmt"
	"gim/internal/lg"
	"gim/model"
	"net/http"

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
	err := s.router.Run(fmt.Sprintf(":%s", GetConfig().ServerPort))
	if err != nil {
		lg.Logger().Error("API 服务启动失败", zap.Error(err))
	}
}

// example: curl -X POST --header 'Content-Type: application/json' -d '{"message": "hello"}' http://localhost:8082/sendMsg
func sendMsg(ctx *gin.Context) {
	// TODO: send message，through grpc

	ctx.JSON(http.StatusCreated, model.CodeResult{
		Code:    "0",
		Message: "success",
	})
	lg.Logger().Info("成功")
}
