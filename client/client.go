// Author: xufei
// Date: 2019-09-04 14:55

package client

import (
	"fmt"
	"gim/internal/lg"
	"gim/internal/util"

	"google.golang.org/grpc"
)

type Client struct {
	rpc       *rpcServer
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

	ctx := &context{c}

	server := newHTTPServer()
	c.waitGroup.Wrap(func() {
		server.Run()
	})

	conn, err := grpc.Dial(fmt.Sprintf("%s:%s", GetConfig().ServerIP, GetConfig().ServerRPCPort), grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	c.rpc = newRpcServer(conn)

	uClient, err := newUserClient(ctx)
	if err != nil {
		panic(err)
	}
	c.waitGroup.Wrap(func() {
		uClient.dispatch()
	})
	c.waitGroup.Wrap(func() {
		uClient.Login()
	})

	scanner := NewScan(ctx)
	c.waitGroup.Wrap(func() {
		scanner.Scan()
	})

	c.waitGroup.Wait()
	lg.Logger().Info("Client: done!")
}
