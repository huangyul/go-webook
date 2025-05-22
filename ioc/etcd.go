package ioc

import (
	"github.com/spf13/viper"
	clientv3 "go.etcd.io/etcd/client/v3"
)

func InitEtcd() *clientv3.Client {
	addr := viper.GetString("etcd.addr")

	client, err := clientv3.NewFromURL(addr)
	if err != nil {
		panic(err)
	}

	return client

}
