package circuitbreaker

import (
	"context"

	"github.com/go-kratos/aegis/circuitbreaker"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Builder struct {
	breaker circuitbreaker.CircuitBreaker
}

func (b *Builder) Build() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		er := b.breaker.Allow()
		if er != nil {
			b.breaker.MarkFailed()
			return nil, status.Errorf(codes.Unavailable, "circuitbreak")
		} else {
			resp, err := handler(ctx, req)
			if err != nil {
				b.breaker.MarkFailed()
			} else {
				b.breaker.MarkSuccess()
			}
			return resp, err
		}
	}
}
