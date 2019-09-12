// Author: xufei
// Date: 2019-09-11 15:56

package etcdkit

import (
	"context"
	"fmt"
	"net"
	"strings"

	"github.com/coreos/etcd/clientv3"
)

const (
	schema = "gim_resolver"
)

var deregister = make(chan struct{})

// Register 服务端向 etcd 注册
// target etcd 服务地址，service 服务名称，host，port 注册服务端的地址和端口
// ttl 服务端存活时间
func Register(target, service, host, port string, ttl int) error {
	serverValue := net.JoinHostPort(host, port)
	serverKey := fmt.Sprintf("/%s/%s/%s", schema, service, serverValue)

	// get etcd client
	client, err := clientv3.New(clientv3.Config{
		Endpoints: strings.Split(target, ","),
	})
	if err != nil {
		return err
	}

	// 获得一个 lease grant
	resp, err := client.Grant(context.TODO(), int64(ttl))
	if err != nil {
		return err
	}

	// TODO: 应该先获取看是否已经存在，不存在然后在 put
	if _, err := client.Put(context.TODO(), serverKey, serverValue, clientv3.WithLease(resp.ID)); err != nil {
		return err
	}

	// 设置与 etcd 的长连接存活
	if _, err = client.KeepAlive(context.TODO(), resp.ID); err != nil {
		return err
	}

	// 创建注销监测服务
	go func() {
		<-deregister
		client.Delete(context.TODO(), serverKey)
		deregister <- struct{}{}
	}()

	return err
}

// UnRegister 删除注册
func UnRegister() {
	deregister <- struct{}{}
	<-deregister
}
