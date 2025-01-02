package ioc

import (
	"github.com/spf13/viper"
	etcdv3 "go.etcd.io/etcd/client/v3"
)

func InitEtcd() *etcdv3.Client {
	addr := viper.GetString("etcd.addr")
	if addr == "" {
		panic("etcd addr is empty")
	}
	client, err := etcdv3.NewFromURL(addr)
	if err != nil {
		panic(err)
	}
	return client
}
