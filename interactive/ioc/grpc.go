package ioc

import (
	interGrpc "github.com/huangyul/go-blog/interactive/grpc"
	"github.com/huangyul/go-blog/pkg/grpcx"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

func InitGrpc(s1 *interGrpc.InteractiveServiceServer) *grpcx.Server {
	type Config struct {
		EtcdAddr string `yaml:"etcdAddr"`
		Port     int    `yaml:"port"`
		Name     string `yaml:"name"`
	}
	var cfg Config
	if err := viper.UnmarshalKey("grpc.server", &cfg); err != nil {
		panic(err)
	}
	s := grpc.NewServer()
	s1.Register(s)
	return &grpcx.Server{
		Port:     cfg.Port,
		Name:     cfg.Name,
		EtcdAddr: cfg.EtcdAddr,
		Server:   s,
	}
}
