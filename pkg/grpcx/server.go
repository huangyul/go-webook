package grpcx

import (
	"fmt"
	"google.golang.org/grpc"
	"net"
)

type Server struct {
	Server *grpc.Server
	Addr   string
}

func (s *Server) Serve() error {
	l, err := net.Listen("tcp", s.Addr)
	if err != nil {
		return err
	}
	fmt.Printf("server listening at %s\n", s.Addr)
	return s.Server.Serve(l)
}
