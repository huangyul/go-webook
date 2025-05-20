package grpcx

import (
	"net"

	"google.golang.org/grpc"
)

type Server struct {
	Addr string
	*grpc.Server
}

func (s *Server) Serve() error {
	lis, err := net.Listen("tcp", s.Addr)
	if err != nil {
		panic(err)
	}
	return s.Server.Serve(lis)
}
