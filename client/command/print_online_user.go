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

type printOnlineUserCommand struct {
	config *model.ClientConfig
}

func (c *printOnlineUserCommand) Process(msg string) {
	users := c.reqOnlineUsers()
	c.print(users)
}

func (c *printOnlineUserCommand) print(users []model.User) {
	lg.Logger().Info("~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~")

	for _, user := range users {
		lg.Logger().Info(fmt.Sprintf("userID=%d, userName=%s", user.UserID, user.UserName))
	}

	lg.Logger().Info("~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~")
}

type onlineUsersResp struct {
	model.CodeMessage
	Data []model.User `json:"data"`
}

func (c *printOnlineUserCommand) reqOnlineUsers() []model.User {
	url := fmt.Sprintf("%s/onlineUsers", c.config.ServerURL)

	resp, err := http.Get(url)
	if err != nil {
		lg.Logger().Error("获取在线用户列表失败", zap.Error(err))
		return nil
	}
	defer resp.Body.Close()

	respBody, _ := ioutil.ReadAll(resp.Body)
	lg.Logger().Info("获取结果" + string(respBody))

	var result onlineUsersResp
	if err := json.Unmarshal(respBody, &result); err != nil {
		lg.Logger().Error(err.Error())
	}

	return result.Data
}

func NewPrintOnlineUserCommand(config *model.ClientConfig) *printOnlineUserCommand {
	return &printOnlineUserCommand{
		config: config,
	}
}
