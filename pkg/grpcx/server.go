package grpcx

import (
	"context"
	"fmt"
	"github.com/huangyul/go-blog/pkg/netx"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/naming/endpoints"
	"google.golang.org/grpc"
	"net"
	"strconv"
	"time"
)

type Server struct {
	*grpc.Server
	EtcdAddr string
	Port     int
	Name     string
	client   *clientv3.Client
	KaCancel context.CancelFunc
}

func (s *Server) Serve() error {
	add := ":" + strconv.Itoa(s.Port)
	lis, err := net.Listen("tcp", add)
	if err != nil {
		return err
	}
	s.Register()
	return s.Server.Serve(lis)
}

func (s *Server) Register() error {
	client, err := clientv3.NewFromURL(s.EtcdAddr)
	if err != nil {
		return err
	}
	s.client = client
	em, err := endpoints.NewManager(client, "service/"+s.Name)

	addr := netx.GetOutboundIP() + ":" + strconv.Itoa(s.Port)

	key := "service/" + s.Name + "/" + addr

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// 租期
	var ttl int64 = 5
	leaseResp, err := s.client.Grant(ctx, ttl)
	if err != nil {
		return err
	}

	err = em.AddEndpoint(ctx, key, endpoints.Endpoint{
		Addr: addr,
	}, clientv3.WithLease(leaseResp.ID))
	if err != nil {
		return err
	}

	kaCtx, cancel := context.WithCancel(context.Background())
	s.KaCancel = cancel
	// 保持心跳和续约
	keepResp, err := s.client.KeepAlive(kaCtx, leaseResp.ID)
	go func() {
		for resp := range keepResp {
			// 记录日志
			fmt.Print(resp.String())
		}
	}()

	return err
}

func (s *Server) Close() error {
	if s.KaCancel != nil {
		s.KaCancel()
	}
	if s.client != nil {
		return s.client.Close()
	}
	s.GracefulStop()
	return nil
}
