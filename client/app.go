// Author: xufei
// Date: 2019-09-04 14:55

package client

import (
	"gim/internal/ciface"
	"gim/internal/lg"
	"gim/internal/util"
	"gim/pkg/etcdkit"

	"google.golang.org/grpc/resolver"
)

type AppClient struct {
	etcdResolver resolver.Builder
	waitGroup    util.WaitGroupWrapper
	userClient   ciface.UserClient
}

func New() *AppClient {
	return &AppClient{}
}

func (c *AppClient) Main() {
	if err := InitConfig(); err != nil {
		panic(err)
	}

	// 设置日志的级别
	lg.SetLevel(GetConfig().LogLevel)

	ctx := &context{c}

	server := newHTTPServer(ctx)
	c.waitGroup.Wrap(func() {
		server.Run()
	})

	// 服务发现注册
	c.etcdResolver = etcdkit.NewResolver(GetConfig().EtcdUrl, GetConfig().EtcdServerName, lg.Logger())
	resolver.Register(c.etcdResolver)

	c.userClient = newUserClient(ctx, GetConfig())
	if err := c.userClient.Start(); err != nil {
		panic(err)
	}

	scanner := NewScan(ctx)
	c.waitGroup.Wrap(func() {
		scanner.Scan()
	})

	c.waitGroup.Wait()
	lg.Logger().Info("AppClient: done!")
}
