package ioc

import (
	"github.com/fsnotify/fsnotify"
	interv1 "github.com/huangyul/go-blog/api/proto/gen/inter/v1"
	"github.com/huangyul/go-blog/interactive/service"
	"github.com/huangyul/go-blog/internal/client/interactive"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func InitInteractiveGrpcClient(svc service.InteractiveService) interv1.InteractiveServiceClient {

	type Config struct {
		Addr      string `yaml:"addr"`
		Secure    bool   `yaml:"secure"`
		Threshold int    `yaml:"threshold"`
	}
	var cfg Config
	if err := viper.UnmarshalKey("grpc.client.inter", &cfg); err != nil {
		panic(err)
	}
	var opts []grpc.DialOption

	if !cfg.Secure {
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}

	cc, err := grpc.NewClient(cfg.Addr, opts...)
	if err != nil {
		panic(err)
	}

	remote := interv1.NewInteractiveServiceClient(cc)
	local := interactive.NewLocalInteractiveClient(svc)
	res := interactive.NewClient(remote, local)

	viper.OnConfigChange(func(in fsnotify.Event) {
		cfg = Config{}
		if err := viper.UnmarshalKey("grpc.client.inter", &cfg); err != nil {
			panic(err)
		}
		res.SetThreshold(cfg.Threshold)
	})

	return res

}
