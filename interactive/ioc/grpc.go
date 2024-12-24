package ioc

import (
	interGrpc "github.com/huangyul/go-blog/interactive/grpc"
	"github.com/huangyul/go-blog/pkg/grpcx"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

func InitGrpc(s1 *interGrpc.InteractiveServiceServer) *grpcx.Server {
	s := grpc.NewServer()
	s1.Register(s)
	return &grpcx.Server{
		Server: s,
		Addr:   viper.GetString("grpc.server.addr"),
	}
}
