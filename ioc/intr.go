package ioc

import (
	intrv1 "github.com/huangyul/go-webook/api/proto/gen/intr/v1"
	"github.com/huangyul/go-webook/interactive/service"
	"github.com/huangyul/go-webook/internal/client"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func InitInteractiveClient(svc service.InteractiveService) intrv1.InteractiveServiceClient {

	intrAddr := viper.GetString("grpc.client.interactive")

	cc, err := grpc.NewClient(intrAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}

	remote := intrv1.NewInteractiveServiceClient(cc)
	local := client.NewInteractiveServiceAdapter(svc)

	c := client.NewInteractiveClient(remote, local)

	return c
}
