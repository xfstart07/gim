// Author: xufei
// Date: 2019-09-04

package server

import (
	"fmt"
	"gim/internal/http_helper"
	"gim/internal/lg"
	"gim/model"
	"gim/server/constant"

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
	s.router.POST("/sendMsg", s.sendMsg)
	s.router.POST("/sendP2PMsg", s.sendP2PMsg)
	s.router.POST("/sendGroupMsg", s.sendGroupMsg)
	s.router.GET("/onlineUsers", s.onlineUsers)
	s.router.POST("/offlineUser", s.offlineUser)
}

func (s *httpServer) Run() {
	err := s.router.Run(fmt.Sprintf(":%s", GetConfig().ServerPort))
	if err != nil {
		lg.Logger().Fatal(err.Error())
	}
}

// example: curl -X POST --header 'Content-Type: application/json' -d '{"user_name": "leon"}' http://localhost:8081/registerAccount
func (s *httpServer) registerAccount(ctx *gin.Context) {
	user := model.User{}
	err := ctx.Bind(&user)
	if err != nil {
		http_helper.Render400(ctx, err)
		return
	}

	// register account
	register, err := s.ctx.server.accountRegister(user)
	if err != nil {
		lg.Logger().Error("注册失败", zap.Error(err))
		// TODO: 如何判断 err
		//if errors.Cause(err) == constant.ErrAccountRegistered {
		//	err = constant.ErrAccountRegistered
		//}
		http_helper.Render500(ctx, constant.ErrServerFail)
		return
	}

	http_helper.RenderCreated(ctx, register)
	lg.Logger().Info("注册用户成功", zap.Any("account", register))
}

// sendMsg 想指定用户发送消息
// POST
// example: curl -X POST --header 'Content-Type: application/json' -d '{"user_id": 1567750270024892000, "msg": "你好"}' http://localhost:8082/sendMsg
func (s *httpServer) sendMsg(ctx *gin.Context) {
	msg := model.MsgReq{}
	err := ctx.Bind(&msg)
	if err != nil {
		http_helper.Render400(ctx, err)
		return
	}

	err = s.ctx.server.sendMsg(msg)
	if err != nil {
		lg.Logger().Error("发送消息失败", zap.Error(err))
		http_helper.Render500(ctx, errMessageSendFailed)
		return
	}

	http_helper.RenderCreated(ctx, nil)
}

// sendP2PMsg 用户私聊
// POST
// example: curl -X POST --header 'Content-Type: application/json' -d '{"user_id": 1567750270024892000, "msg": "你好", "receiver_id": 1567750270024892000,}' http://localhost:8082/sendP2PMsg
func (s *httpServer) sendP2PMsg(ctx *gin.Context) {
	msg := model.P2PReq{}
	err := ctx.Bind(&msg)
	if err != nil {
		http_helper.Render400(ctx, err)
		return
	}

	err = s.ctx.server.sendP2PMsg(msg)
	if err != nil {
		lg.Logger().Error("发送消息失败", zap.Error(err))
		http_helper.Render500(ctx, errMessageSendFailed)
		return
	}

	http_helper.RenderCreated(ctx, nil)
}

// sendGroupMsg 用户群聊
// POST
// example: curl -X POST --header 'Content-Type: application/json' -d '{"user_id": 1567750270024892000, "msg": "你好"}' http://localhost:8082/sendGroupMsg
func (s *httpServer) sendGroupMsg(ctx *gin.Context) {
	msg := model.MsgReq{}
	err := ctx.Bind(&msg)
	if err != nil {
		http_helper.Render400(ctx, err)
		return
	}

	err = s.ctx.server.sendGroupMsg(msg)
	if err != nil {
		lg.Logger().Error("发送消息失败", zap.Error(err))
		http_helper.Render500(ctx, errMessageSendFailed)
		return
	}

	http_helper.RenderCreated(ctx, nil)
}

func (s *httpServer) onlineUsers(ctx *gin.Context) {
	users := s.ctx.server.getOnlineUsers()

	http_helper.RenderOK(ctx, users)
}

func (s *httpServer) offlineUser(ctx *gin.Context) {
	msg := model.MsgReq{}
	err := ctx.Bind(&msg)
	if err != nil {
		http_helper.Render400(ctx, err)
		return
	}

	// 用户下线
	user := s.ctx.server.accountSrv.GetSessionByUserID(msg.UserID)
	s.ctx.server.userOffline(user)

	http_helper.RenderCreated(ctx, nil)
}
