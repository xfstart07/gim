// Author: xufei
// Date: 2019-09-27 10:55

package command

import (
	"encoding/json"
	"fmt"
	"gim/internal/lg"
	"gim/model"
	"io/ioutil"
	"net/http"

	"go.uber.org/zap"
)

type PrintOnlineUserCommand struct {
	config *model.ClientConfig
}

func (c *PrintOnlineUserCommand) Process(msg string) {
	users := c.reqOnlineUsers()
	c.print(users)
}

func (c *PrintOnlineUserCommand) print(users []model.User) {
	lg.Logger().Info("~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~")

	for _, user := range users {
		lg.Logger().Info(fmt.Sprintf("userID=%d, userName=%s", user.UserID, user.UserName))
	}

	lg.Logger().Info("~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~")
}

func (c *PrintOnlineUserCommand) reqOnlineUsers() []model.User {
	url := fmt.Sprintf("http://%s:%s/onlineUsers", c.config.ServerIP, c.config.ServerPort)

	resp, err := http.Get(url)
	if err != nil {
		lg.Logger().Error("获取在线用户列表失败", zap.Error(err))
		return nil
	}
	defer resp.Body.Close()

	respBody, _ := ioutil.ReadAll(resp.Body)
	lg.Logger().Info("获取结果" + string(respBody))

	var users []model.User
	_ = json.Unmarshal(respBody, &users)

	return users
}

func NewPrintOnlineUserCommand(config *model.ClientConfig) *PrintOnlineUserCommand {
	return &PrintOnlineUserCommand{
		config: config,
	}
}
