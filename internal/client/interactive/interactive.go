package interactive

import (
	"context"
	interv1 "github.com/huangyul/go-blog/api/proto/gen/inter/v1"
	"google.golang.org/grpc"
	"math/rand"
	"sync/atomic"
)

type Client struct {
	remote    interv1.InteractiveServiceClient
	local     interv1.InteractiveServiceClient
	threshold atomic.Value
}

var _ interv1.InteractiveServiceClient = (*Client)(nil)

func NewClient(remote interv1.InteractiveServiceClient, local interv1.InteractiveServiceClient) *Client {
	c := &Client{
		remote:    remote,
		local:     local,
		threshold: atomic.Value{},
	}
	c.threshold.Store(100)
	return c
}

func (c *Client) SetThreshold(threshold int) {
	c.threshold.Store(threshold)
}

func (c *Client) selectClient() interv1.InteractiveServiceClient {
	r := rand.Intn(100)
	if r < c.threshold.Load().(int) {
		return c.remote
	}
	return c.local
}

func (c *Client) IncrReadCnt(ctx context.Context, in *interv1.IncrReadCntRequest, opts ...grpc.CallOption) (*interv1.EmptyResponse, error) {
	return c.selectClient().IncrReadCnt(ctx, in, opts...)
}

func (c *Client) Like(ctx context.Context, in *interv1.LikeRequest, opts ...grpc.CallOption) (*interv1.EmptyResponse, error) {
	return c.selectClient().Like(ctx, in, opts...)
}

func (c *Client) CancelLike(ctx context.Context, in *interv1.CancelLikeRequest, opts ...grpc.CallOption) (*interv1.EmptyResponse, error) {
	return c.selectClient().CancelLike(ctx, in, opts...)
}

func (c *Client) Collect(ctx context.Context, in *interv1.CollectRequest, opts ...grpc.CallOption) (*interv1.EmptyResponse, error) {
	return c.selectClient().Collect(ctx, in, opts...)
}

func (c *Client) Get(ctx context.Context, in *interv1.GetRequest, opts ...grpc.CallOption) (*interv1.InteractiveResponse, error) {
	return c.selectClient().Get(ctx, in, opts...)
}
