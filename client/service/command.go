// Author: xufei
// Date: 2019-09-27 11:15

package service

import (
	"gim/client/command"
	"gim/internal/ciface"
	"gim/model"
)

type InnerCommandContext struct {
	config *model.ClientConfig
}

func NewInnerCommandContext(cfg *model.ClientConfig) *InnerCommandContext {
	return &InnerCommandContext{
		config: cfg,
	}
}

func (c *InnerCommandContext) CreateCommander(userClient ciface.UserClient, cmd string) ciface.InnerCommander {
	switch cmd {
	case command.SystemCommandPrintOnlineUser:
		return command.NewPrintOnlineUserCommand(c.config)
	case command.SystemCommandShutdown:
		return command.NewShutDownCommand(userClient)
	}

	return nil
}
