// Author: xufei
// Date: 2019-09-04 14:55

package client

import (
	"fmt"
	"gim/internal/lg"
	"gim/internal/util"
	"gim/pkg/etcdkit"

	"google.golang.org/grpc/balancer/roundrobin"

	"google.golang.org/grpc/resolver"

	"google.golang.org/grpc"
)

type AppClient struct {
	rpc       *rpcServer
	waitGroup util.WaitGroupWrapper
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

	server := newHTTPServer()
	c.waitGroup.Wrap(func() {
		server.Run()
	})

	// 服务发现注册
	etcdResolver := etcdkit.NewResolver(GetConfig().EtcdUrl, GetConfig().EtcdServerName)
	resolver.Register(etcdResolver)

	rpcDiaUrl := fmt.Sprintf("%s://authority/%s", etcdResolver.Scheme(), GetConfig().EtcdServerName)
	conn, err := grpc.Dial(rpcDiaUrl, grpc.WithBalancerName(roundrobin.Name), grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	c.rpc = newRpcServer(conn)

	uClient, err := newUserClient(ctx, GetConfig())
	if err != nil {
		panic(err)
	}
	c.waitGroup.Wrap(func() {
		uClient.Login()
	})

	scanner := NewScan(ctx)
	c.waitGroup.Wrap(func() {
		scanner.Scan()
	})

	c.waitGroup.Wait()
	lg.Logger().Info("AppClient: done!")
}
