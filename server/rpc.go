// Author: xufei
// Date: 2019-09-06 09:40

package server

import (
	"gim/internal/lg"
	"gim/internal/rpc_service"
	"net"

	"google.golang.org/grpc"
)

type rpcServer struct {
	ctx    *context
	server *grpc.Server
}

func NewRpcServer(ctx *context) *rpcServer {
	s := grpc.NewServer()
	rpc_service.RegisterGIMServiceServer(s, NewChannelService(ctx))

	return &rpcServer{
		ctx:    ctx,
		server: s,
	}
}

func (s *rpcServer) Run(listener net.Listener) {
	lg.Logger().Info("rpc server...")
	err := s.server.Serve(listener)
	if err != nil {
		panic(err)
	}
}
