// Author: xufei
// Date: 2019-09-09 09:30

package client

import (
	"fmt"
	"gim/client/handler"
	"gim/internal/lg"

	"go.uber.org/zap"
)

type scanner struct {
	ctx            *context
	messageHandler handler.MessageHandleInterface
}

func NewScan(ctx *context) *scanner {
	return &scanner{
		ctx:            ctx,
		messageHandler: handler.NewMessageHandler(GetConfig()),
	}
}

func (s *scanner) Scan() {
	for {
		var msg string
		lg.Logger().Info("请输入: ")
		_, err := fmt.Scanf("%s\n", &msg)
		if err != nil {
			lg.Logger().Info("不能输入空消息!")
			continue
		}
		lg.Logger().Debug("用户输入消息: " + msg)

		if !s.messageHandler.CheckMsg(msg) {
			continue
		}

		err = s.messageHandler.SendMsg(msg)
		if err != nil {
			lg.Logger().Error("消息发送失败", zap.Error(err))
		}
	}
}
