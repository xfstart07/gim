// Author: xufei
// Date: 2019-09-04 14:55

package client

import (
	"gim/internal/lg"
	"gim/internal/util"
)

type Client struct {
	waitGroup util.WaitGroupWrapper
}

func New() *Client {
	return &Client{}
}

func (c *Client) Main() {
	if err := InitConfig(); err != nil {
		panic(err)
	}

	// 设置日志的级别
	lg.SetLevel(GetConfig().LogLevel)

	server := newHTTPServer()
	c.waitGroup.Wrap(func() {
		server.Run()
	})

	c.waitGroup.Wait()
	lg.Logger().Info("Client: done!")
}
