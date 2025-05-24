package ioc

import (
	"context"
	"fmt"

	intrv1 "github.com/huangyul/go-webook/api/proto/gen/intr/v1"
	"github.com/huangyul/go-webook/interactive/service"
	interactiveClient "github.com/huangyul/go-webook/internal/client"
	"github.com/spf13/viper"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/naming/resolver"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func InitInteractiveClient(svc service.InteractiveService, client *clientv3.Client) intrv1.InteractiveServiceClient {

	kvs, errr := client.Get(context.Background(), "/service/interactive/", clientv3.WithPrefix())
	if errr != nil {
		panic(errr)
	}

	for key, val := range kvs.Kvs {
		fmt.Printf("key: %d, val: %s", key, val.String())
	}

	name := viper.GetString("grpc.client.interactive")

	reso, err := resolver.NewBuilder(client)
	if err != nil {
		panic(err)
	}

	cc, err := grpc.NewClient("etcd:///service/"+name, grpc.WithResolvers(reso), grpc.WithDefaultServiceConfig(`{"loadBalancingConfig":[{"round_robin":{}}]}`), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}

	remote := intrv1.NewInteractiveServiceClient(cc)
	local := interactiveClient.NewInteractiveServiceAdapter(svc)

	c := interactiveClient.NewInteractiveClient(remote, local)

	return c
}
