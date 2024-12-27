package interactive

import (
	"context"
	interv1 "github.com/huangyul/go-blog/api/proto/gen/inter/v1"
	"github.com/huangyul/go-blog/interactive/service"
	"google.golang.org/grpc"
)

type LocalClient struct {
	svc service.InteractiveService
}

func NewLocalInteractiveClient(svc service.InteractiveService) interv1.InteractiveServiceClient {
	return &LocalClient{svc: svc}
}

func (i *LocalClient) IncrReadCnt(ctx context.Context, in *interv1.IncrReadCntRequest, opts ...grpc.CallOption) (*interv1.EmptyResponse, error) {
	err := i.svc.IncrReadCnt(ctx, in.GetBizId(), in.GetBiz())
	return &interv1.EmptyResponse{}, err
}

func (i *LocalClient) Like(ctx context.Context, in *interv1.LikeRequest, opts ...grpc.CallOption) (*interv1.EmptyResponse, error) {
	err := i.svc.Like(ctx, in.GetUserId(), in.GetBizId(), in.GetBiz())
	return &interv1.EmptyResponse{}, err
}

func (i *LocalClient) CancelLike(ctx context.Context, in *interv1.CancelLikeRequest, opts ...grpc.CallOption) (*interv1.EmptyResponse, error) {
	err := i.svc.CancelLike(ctx, in.GetUserId(), in.GetBizId(), in.GetBiz())
	return &interv1.EmptyResponse{}, err
}

func (i *LocalClient) Collect(ctx context.Context, in *interv1.CollectRequest, opts ...grpc.CallOption) (*interv1.EmptyResponse, error) {
	err := i.svc.Collect(ctx, in.GetUid(), in.GetId(), in.GetCid(), in.GetBiz())
	return &interv1.EmptyResponse{}, err
}

func (i *LocalClient) Get(ctx context.Context, in *interv1.GetRequest, opts ...grpc.CallOption) (*interv1.InteractiveResponse, error) {
	inter, err := i.svc.Get(ctx, in.GetUid(), in.GetId(), in.GetBiz())
	if err != nil {
		return &interv1.InteractiveResponse{}, err
	}
	return &interv1.InteractiveResponse{
		CollectCnt: int32(inter.CollectCnt),
		ReadCnt:    int32(inter.ReadCnt),
		LikeCnt:    int32(inter.LikeCnt),
		Collected:  inter.Collected,
		Liked:      inter.Liked,
	}, nil
}
