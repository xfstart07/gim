// Author: xufei
// Date: 2019-09-11 17:23

package etcdkit

import (
	"context"
	"fmt"
	"gim/internal/lg"
	"strings"

	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/mvcc/mvccpb"
	"google.golang.org/grpc/resolver"
)

// resolver is the implementaion of grpc.resolve.Builder, grpc.resolver.Resolver
type Resolver struct {
	target  string
	service string
	cli     *clientv3.Client
	cc      resolver.ClientConn
}

// NewResolver return resolver builder
// target example: "http://127.0.0.1:2379,http://127.0.0.1:12379,http://127.0.0.1:22379"
// service is service name
func NewResolver(target string, service string) resolver.Builder {
	return &Resolver{target: target, service: service}
}

// Scheme return etcdv3 schema
func (r *Resolver) Scheme() string {
	return schema
}

// ResolveNow
func (r *Resolver) ResolveNow(rn resolver.ResolveNowOption) {
}

// Close
func (r *Resolver) Close() {
}

// Build to resolver.Resolver
func (r *Resolver) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOption) (resolver.Resolver, error) {
	var err error

	r.cli, err = clientv3.New(clientv3.Config{
		Endpoints: strings.Split(r.target, ","),
	})
	if err != nil {
		return nil, fmt.Errorf("grpclb: create clientv3 client failed: %v", err)
	}

	r.cc = cc

	go r.watch(fmt.Sprintf("/%s/%s/", schema, r.service))

	return r, nil
}

func (r *Resolver) watch(prefix string) {
	lg.Logger().Info("etcd resolver watch...")

	addrDict := make(map[string]resolver.Address)

	// 更新新发现的地址
	update := func() {
		addrList := make([]resolver.Address, 0, len(addrDict))
		for _, v := range addrDict {
			addrList = append(addrList, v)
		}

		r.cc.UpdateState(resolver.State{
			Addresses: addrList,
		})
		//r.cc.NewAddress(addrList)
	}

	resp, err := r.cli.Get(context.Background(), prefix, clientv3.WithPrefix())
	if err == nil {
		for i := range resp.Kvs {
			addrDict[string(resp.Kvs[i].Value)] = resolver.Address{Addr: string(resp.Kvs[i].Value)}
		}
	}

	update()

	rch := r.cli.Watch(context.Background(), prefix, clientv3.WithPrefix(), clientv3.WithPrevKV())
	for n := range rch {
		for _, ev := range n.Events {
			switch ev.Type {
			case mvccpb.PUT:
				addrDict[string(ev.Kv.Key)] = resolver.Address{Addr: string(ev.Kv.Value)}
			case mvccpb.DELETE:
				lg.Logger().Info("Delete:=" + string(ev.PrevKv.Key))
				delete(addrDict, string(ev.PrevKv.Key))
			}
		}

		update()
	}
}
