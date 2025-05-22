package grpcx

import (
	"context"
	"fmt"
	"net"
	"strconv"

	"github.com/huangyul/go-webook/pkg/netx"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/naming/endpoints"
	"google.golang.org/grpc"
)

type Server struct {
	*grpc.Server
	Prot   int
	Name   string
	Client *clientv3.Client
}

func (s *Server) Serve() error {
	lis, err := net.Listen("tcp", "localhost:"+strconv.Itoa(s.Prot))
	if err != nil {
		panic(err)
	}

	s.registerEtcd()

	return s.Server.Serve(lis)
}

func (s *Server) registerEtcd() error {
	em, err := endpoints.NewManager(s.Client, "service/"+s.Name)
	if err != nil {
		panic(err)
	}

	leasResp, err := s.Client.Grant(context.Background(), 5)
	if err != nil {
		panic(err)
	}

	key := "service/" + s.Name + "/" + netx.GetOutboundIP() + strconv.Itoa(s.Prot)

	err = em.AddEndpoint(context.Background(), key, endpoints.Endpoint{
		Addr: netx.GetOutboundIP() + ":" + strconv.Itoa(s.Prot),
	}, clientv3.WithLease(leasResp.ID))
	if err != nil {
		panic(err)
	}

	keepAliveCh, er := s.Client.KeepAlive(context.Background(), leasResp.ID)
	go func() {
		if er != nil {
			fmt.Printf("failure to renew contract")
		}
		for ch := range keepAliveCh {
			fmt.Println(ch)
			if ch == nil {
				fmt.Println("service deleted from etcd")
			}
		}
	}()
	return err
}
