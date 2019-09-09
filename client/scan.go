// Author: xufei
// Date: 2019-09-09 09:30

package client

import (
	"fmt"
	"gim/internal/lg"
	"gim/model"
	"strconv"
	"strings"

	"go.uber.org/zap"
)

// message format: userID;;message
func (c *Client) Scan() {
	for {
		var msg string
		lg.Logger().Info("请输入: ")
		_, err := fmt.Scanf("%s\n", &msg)
		if err != nil {
			lg.Logger().Info("不能输入空消息!")
			continue
		}
		lg.Logger().Debug("用户输入消息: " + msg)

		if strings.Contains(msg, ";;") {
			msgs := strings.Split(msg, ";;")
			userID, err := strconv.ParseInt(msgs[0], 10, 64)
			if err != nil {
				lg.Logger().Error("用户ID输入错误", zap.Error(err))
				continue
			}

			// p2p chat
			_ = c.sendP2PMsg(model.P2PReq{
				ReceiverID: userID,
				UserID:     GetConfig().UserID,
				Msg:        msgs[1],
			})
		} else {
			// group chat
			_ = c.sendGroupMsg(model.MsgReq{
				UserID: GetConfig().UserID,
				Msg:    msg,
			})
		}
	}
}
