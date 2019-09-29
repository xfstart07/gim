// Author: xufei
// Date: 2019-09-09 09:40

package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"gim/client/service"
	"gim/internal/lg"
	"gim/model"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"go.uber.org/zap"
)

var (
	errP2PSendFail = errors.New("私聊消息发送失败")
)

type messageHandler struct {
	config          *model.ClientConfig
	innerCommandCtx *service.InnerCommandContext
}

func NewMessageHandler(cfg *model.ClientConfig) *messageHandler {
	return &messageHandler{
		config:          cfg,
		innerCommandCtx: service.NewInnerCommandContext(cfg),
	}
}

func (h *messageHandler) CheckMsg(msg string) bool {
	if msg == "" {
		lg.Logger().Info("不能输入空消息!")
		return false
	}
	return true
}

func (h *messageHandler) SendMsg(msg string) {
	var err error

	msgStrings := strings.Split(msg, ";;")
	if len(msgStrings) > 1 {
		userID, _ := strconv.ParseInt(msgStrings[0], 10, 64)

		// p2p chat
		err = h.sendP2PMsg(model.P2PReq{
			ReceiverID: userID,
			UserID:     h.config.UserID,
			Msg:        msgStrings[1],
		})
	} else {
		// group chat
		err = h.sendGroupMsg(model.MsgReq{
			UserID: h.config.UserID,
			Msg:    msg,
		})
	}

	if err != nil {
		lg.Logger().Error(err.Error())
	}
}

func (h *messageHandler) sendP2PMsg(req model.P2PReq) error {
	url := fmt.Sprintf("http://%s:%s/sendP2PMsg", h.config.ServerIP, h.config.ServerPort)

	msgBody, _ := json.Marshal(req)
	lg.Logger().Debug(fmt.Sprintf("post = %s, %s", url, string(msgBody)))

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(msgBody))
	if err != nil {
		lg.Logger().Error("私聊消息发送失败", zap.Error(err))
		return errP2PSendFail
	}
	defer resp.Body.Close()

	respBody, _ := ioutil.ReadAll(resp.Body)
	lg.Logger().Info("发送结果" + string(respBody))

	return nil
}

func (h *messageHandler) sendGroupMsg(req model.MsgReq) error {
	url := fmt.Sprintf("http://%s:%s/sendGroupMsg", h.config.ServerIP, h.config.ServerPort)

	msgBody, _ := json.Marshal(req)
	lg.Logger().Debug(fmt.Sprintf("post = %s, %s", url, string(msgBody)))

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(msgBody))
	if err != nil {
		lg.Logger().Error("群聊消息发送失败", zap.Error(err))
		return errP2PSendFail
	}
	defer resp.Body.Close()

	respBody, _ := ioutil.ReadAll(resp.Body)
	lg.Logger().Info("发送结果" + string(respBody))

	return nil
}

func (h *messageHandler) InnerCommand(msg string) bool {
	if strings.HasPrefix(msg, ":") {
		commander := h.innerCommandCtx.CreateCommander(msg)
		if commander != nil {
			commander.Process(msg)
		} else {
			lg.Logger().Warn("命令错误!")
		}

		return true
	}

	return false
}
