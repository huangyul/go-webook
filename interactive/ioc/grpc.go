package ioc

import (
	intrGrpc "github.com/huangyul/go-webook/interactive/grpc"
	"github.com/huangyul/go-webook/pkg/grpcx"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

func InitGrpcServer(s1 *intrGrpc.InteractiveService) *grpcx.Server {
	server := grpc.NewServer()

	s1.Register(server)

	addr := viper.GetString("grpc.addr")
	return &grpcx.Server{
		Server: server,
		Addr:   addr,
	}
}
