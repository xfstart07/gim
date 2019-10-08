// Author: xufei
// Date: 2019-09-06 10:07

package client

import (
	context2 "context"
	"fmt"
	"gim/model"
	"gim/pkg/rpc_service"
	"time"

	"google.golang.org/grpc"
)

type rpcServer struct {
	ctx    *context
	config *model.ClientConfig
	conn   *grpc.ClientConn
}

func newRpcServer(ctx *context, config *model.ClientConfig) *rpcServer {
	rpcDiaUrl := fmt.Sprintf("%s://authority/%s", ctx.client.etcdResolver.Scheme(), config.EtcdServerName)
	ctxTime, cancel := context2.WithTimeout(context2.Background(), 10*time.Second)
	defer cancel()

	// https://github.com/grpc/grpc/blob/master/doc/service_config.md
	// conn, err := grpc.DialContext(ctx, rpcDiaUrl, grpc.WithBalancerName(roundrobin.Name), grpc.WithInsecure())
	conn, err := grpc.DialContext(ctxTime, rpcDiaUrl, grpc.WithInsecure(),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`))
	if err != nil {
		panic(err)
	}

	return &rpcServer{
		ctx:    ctx,
		config: config,
		conn:   conn,
	}
}

func (s *rpcServer) GetChannel() (rpc_service.GIMService_ChannelClient, error) {
	return rpc_service.NewGIMServiceClient(s.conn).Channel(context2.Background())
}
