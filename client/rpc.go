// Author: xufei
// Date: 2019-09-06 10:07

package client

import "google.golang.org/grpc"

type rpcServer struct {
	conn *grpc.ClientConn
}

func newRpcServer(conn *grpc.ClientConn) *rpcServer {
	return &rpcServer{
		conn: conn,
	}
}
