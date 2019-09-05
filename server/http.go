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
	ctx    *context
	router *gin.Engine
}

func newHTTPServer(ctx *context) *httpServer {
	server := &httpServer{
		ctx:    ctx,
		router: gin.Default(),
	}

	server.setRouter()

	return server
}

// 注册接口
func (s *httpServer) setRouter() {
	s.router.POST("/registerAccount", s.registerAccount)
}

func (s *httpServer) Run() {
	err := s.router.Run(fmt.Sprintf(":%s", GetConfig().ServerPort))
	if err != nil {
		lg.Logger().Error("API 服务启动失败", zap.Error(err))
	}
}

// example: curl -X POST --header 'Content-Type: application/json' -d '{"user_name": "leon"}' http://localhost:8081/registerAccount
func (s *httpServer) registerAccount(ctx *gin.Context) {
	user := model.User{}
	err := ctx.Bind(&user)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, model.ErrResult{Message: err.Error()})
		return
	}

	// register account
	register, err := s.ctx.server.accountRegister(user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, model.ErrResult{Message: err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, model.CodeResult{
		Code:    "0",
		Message: "注册成功",
		Data:    register,
	})
	lg.Logger().Info("注册用户成功", zap.Any("account", register))
}
