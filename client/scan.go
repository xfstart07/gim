// Author: xufei
// Date: 2019-09-09 09:30

package client

import (
	"bufio"
	"fmt"
	"gim/client/handler"
	"gim/internal/lg"
	"os"
)

type scanner struct {
	ctx        *context
	buffer     *bufio.Scanner
	msgHandler handler.MessageHandleInterface
}

func NewScan(ctx *context) *scanner {
	return &scanner{
		ctx:        ctx,
		buffer:     bufio.NewScanner(os.Stdin),
		msgHandler: handler.NewMessageHandler(GetConfig()),
	}
}

func (s *scanner) Scan() {
	for {
		var msg string
		lg.Logger().Info("请输入: ")

		s.buffer.Scan()
		msg = s.buffer.Text()
		if msg == "" {
			lg.Logger().Info("不能输入空消息!")
			continue
		}

		if !s.msgHandler.CheckMsg(msg) {
			continue
		}
		s.msgHandler.SendMsg(msg)

		lg.Logger().Info(fmt.Sprintf("%s: 【%s】", GetConfig().Username, msg))
	}
}
