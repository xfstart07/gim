// Author: xufei
// Date: 2019-09-04

package client

import (
	"fmt"
	"gim/client/handler"
	"gim/internal/http_helper"
	"gim/internal/lg"
	"gim/model"

	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
)

type httpServer struct {
	router     *gin.Engine
	msgHandler handler.MessageHandleInterface
}

func newHTTPServer() *httpServer {
	server := &httpServer{
		router:     gin.Default(),
		msgHandler: handler.NewMessageHandler(GetConfig()),
	}

	server.setRouter()

	return server
}

// 注册接口
func (s *httpServer) setRouter() {
	s.router.POST("/sendMsg", s.sendMsg)
}

func (s *httpServer) Run() {
	err := s.router.Run(fmt.Sprintf(":%s", GetConfig().WebPort))
	if err != nil {
		lg.Logger().Error("API 服务启动失败", zap.Error(err))
	}
}

// example: curl -X POST --header 'Content-Type: application/json' -d '{"user_id": 1567750270024892000, "msg": "你好"}' http://localhost:8082/sendMsg
func (s *httpServer) sendMsg(ctx *gin.Context) {
	params := model.MsgReq{}
	if err := ctx.BindJSON(&params); err != nil {
		http_helper.Render400(ctx, err)
		return
	}

	s.msgHandler.SendMsg(fmt.Sprintf("%d;;%s", params.UserID, params.Msg))

	http_helper.RenderCreated(ctx, model.CodeResult{Code: model.CodeSuccessed, Message: "成功"})
}
