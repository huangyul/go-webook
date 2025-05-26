package interceptor

import (
	"context"
	"net"
	"strings"

	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

type Builder struct{}

func (b *Builder) PeerName(ctx context.Context) string {
	return b.grpcHeaderValue(ctx, "app")
}

func (b *Builder) PeerIP(ctx context.Context) string {
	clientIP := b.grpcHeaderValue(ctx, "client-ip")
	if clientIP != "" {
		return clientIP
	}
	pe, ok := peer.FromContext(ctx)
	if !ok {
		return ""
	}
	if pe.Addr == net.Addr(nil) {
		return ""
	}
	if ipStr := strings.Split(pe.Addr.String(), ":"); len(ipStr) > 1 {
		return ipStr[1]
	}
	return ""

}

func (b *Builder) grpcHeaderValue(ctx context.Context, key string) string {
	if key == "" {
		return ""
	}
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ""
	}
	return strings.Join(md.Get(key), ";")
}
