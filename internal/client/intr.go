package client

import (
	"context"
	"math/rand"
	"sync/atomic"

	intrv1 "github.com/huangyul/go-webook/api/proto/gen/intr/v1"
	"google.golang.org/grpc"
)

var _ intrv1.InteractiveServiceClient = (*InteractiveClient)(nil)

type InteractiveClient struct {
	threshold atomic.Value
	remote    intrv1.InteractiveServiceClient
	local     intrv1.InteractiveServiceClient
}

// CancelCollect
func (i *InteractiveClient) CancelCollect(ctx context.Context, in *intrv1.CancelCollectRequest, opts ...grpc.CallOption) (*intrv1.CancelCollectResponse, error) {
	return i.selectClient().CancelCollect(ctx, in, opts...)
}

// CancelLike
func (i *InteractiveClient) CancelLike(ctx context.Context, in *intrv1.CancelLikeRequest, opts ...grpc.CallOption) (*intrv1.CancelLikeResponse, error) {
	return i.selectClient().CancelLike(ctx, in, opts...)
}

// Collect
func (i *InteractiveClient) Collect(ctx context.Context, in *intrv1.CollectRequest, opts ...grpc.CallOption) (*intrv1.CollectResponse, error) {
	return i.selectClient().Collect(ctx, in, opts...)
}

// Get
func (i *InteractiveClient) Get(ctx context.Context, in *intrv1.GetRequest, opts ...grpc.CallOption) (*intrv1.GetResponse, error) {
	return i.selectClient().Get(ctx, in, opts...)
}

// GetByIds
func (i *InteractiveClient) GetByIds(ctx context.Context, in *intrv1.GetByIdsRequest, opts ...grpc.CallOption) (*intrv1.GetByIdsResponse, error) {
	return i.selectClient().GetByIds(ctx, in, opts...)
}

// IncrReadCnt
func (i *InteractiveClient) IncrReadCnt(ctx context.Context, in *intrv1.IncrReadCntRequest, opts ...grpc.CallOption) (*intrv1.IncrReadCntResponse, error) {
	return i.selectClient().IncrReadCnt(ctx, in, opts...)
}

// Like
func (i *InteractiveClient) Like(ctx context.Context, in *intrv1.LikeRequest, opts ...grpc.CallOption) (*intrv1.LikeResponse, error) {
	return i.selectClient().Like(ctx, in, opts...)
}

func NewInteractiveClient(remote intrv1.InteractiveServiceClient, local intrv1.InteractiveServiceClient) intrv1.InteractiveServiceClient {
	client := &InteractiveClient{
		remote: remote,
		local:  local,
	}
	client.threshold.Store(int32(100))
	return client
}

func (i *InteractiveClient) selectClient() intrv1.InteractiveServiceClient {
	// [0, 100) 的随机数
	num := rand.Int31n(100)
	if num < i.threshold.Load().(int32) {
		return i.remote
	}
	return i.local
}
