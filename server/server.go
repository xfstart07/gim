// Author: xufei
// Date: 2019-09-04 17:23

package server

import (
	"fmt"
	"gim/internal/lg"
	"gim/internal/util"
	"net"

	"github.com/go-redis/redis"
)

type Server struct {
	redisClient *redis.Client
	waitGroup   util.WaitGroupWrapper
}

func New() *Server {
	return &Server{}
}

func (s *Server) Main() {
	if err := InitConfig(); err != nil {
		panic(err)
	}

	// 设置日志的级别
	lg.SetLevel(GetConfig().LogLevel)

	// set external service, redis
	s.initRedis()

	ctx := &context{s}
	server := newHTTPServer(ctx)
	s.waitGroup.Wrap(func() {
		server.Run()
	})

	rpcSrv := NewRpcServer(ctx)
	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", GetConfig().RpcPort))
	if err != nil {
		panic(err)
	}
	s.waitGroup.Wrap(func() {
		rpcSrv.Run(listener)
	})

	s.waitGroup.Wait()
	lg.Logger().Info("Client: done!")
}

func (s *Server) initRedis() {
	s.redisClient = redis.NewClient(&redis.Options{
		Addr:     GetConfig().RedisURL,
		Password: GetConfig().RedisPass,
		DB:       GetConfig().RedisDB,
	})

	_, err := s.redisClient.Ping().Result()
	if err != nil {
		panic(err)
	}
	lg.Logger().Info("redis connected...")
}
