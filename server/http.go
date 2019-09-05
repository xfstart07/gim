// Author: xufei
// Date: 2019-09-04

package server

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
	s.router.POST("/registerAccount", registerAccount)
}

func (s *httpServer) Run() {
	err := s.router.Run(fmt.Sprintf(":%s", GetConfig().ServerPort))
	if err != nil {
		lg.Logger().Error("API 服务启动失败", zap.Error(err))
	}
}

// example: curl -X POST --header 'Content-Type: application/json' -d '{"user_name": "leon"}' http://localhost:8081/registerAccount
func registerAccount(ctx *gin.Context) {
	user := model.User{}
	err := ctx.Bind(&user)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, model.ErrResult{Message: err.Error()})
		return
	}

	// TODO: register account，store redis

	ctx.JSON(http.StatusCreated, model.CodeResult{
		Code:    "0",
		Message: "success",
		Data:    user,
	})
	lg.Logger().Info("注册用户成功", zap.Any("params", user))
}
