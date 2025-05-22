package ioc

import (
	intrGrpc "github.com/huangyul/go-webook/interactive/grpc"
	"github.com/huangyul/go-webook/pkg/grpcx"
	"github.com/spf13/viper"
	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc"
)

func InitGrpcServer(s1 *intrGrpc.InteractiveService, client *clientv3.Client) *grpcx.Server {
	server := grpc.NewServer()

	s1.Register(server)

	port := viper.GetInt("grpc.port")
	return &grpcx.Server{
		Server: server,
		Name:   "interactive",
		Client: client,
		Prot:   port,
	}
}
