package client

import (
	"context"

	intrv1 "github.com/huangyul/go-webook/api/proto/gen/intr/v1"
	"github.com/huangyul/go-webook/interactive/domain"
	"github.com/huangyul/go-webook/interactive/service"
	"google.golang.org/grpc"
)

type InteractiveSeviceAdapter struct {
	svc service.InteractiveService
}

// CancelCollect
func (i *InteractiveSeviceAdapter) CancelCollect(ctx context.Context, in *intrv1.CancelCollectRequest, opts ...grpc.CallOption) (*intrv1.CancelCollectResponse, error) {
	err := i.svc.CancelCollect(ctx, in.Biz, in.BizId, in.UserId)
	return &intrv1.CancelCollectResponse{}, err
}

// CancelLike
func (i *InteractiveSeviceAdapter) CancelLike(ctx context.Context, in *intrv1.CancelLikeRequest, opts ...grpc.CallOption) (*intrv1.CancelLikeResponse, error) {
	err := i.svc.CancelCollect(ctx, in.Biz, in.BizId, in.UserId)
	return &intrv1.CancelLikeResponse{}, err
}

// Collect
func (i *InteractiveSeviceAdapter) Collect(ctx context.Context, in *intrv1.CollectRequest, opts ...grpc.CallOption) (*intrv1.CollectResponse, error) {
	err := i.svc.Collect(ctx, in.Biz, in.BizId, in.UserId)
	return &intrv1.CollectResponse{}, err
}

// Get
func (i *InteractiveSeviceAdapter) Get(ctx context.Context, in *intrv1.GetRequest, opts ...grpc.CallOption) (*intrv1.GetResponse, error) {
	res, err := i.svc.Get(ctx, in.Biz, in.BizId, in.UserId)
	if err != nil {
		return nil, err
	}
	return &intrv1.GetResponse{
		Intr: i.toDTO(res),
	}, nil
}

// GetByIds
func (i *InteractiveSeviceAdapter) GetByIds(ctx context.Context, in *intrv1.GetByIdsRequest, opts ...grpc.CallOption) (*intrv1.GetByIdsResponse, error) {
	res, err := i.svc.GetByIds(ctx, in.Biz, in.Ids)
	if err != nil {
		return nil, err
	}
	intrs := make(map[int64]*intrv1.Interactive)
	for k, r := range res {
		intrs[k] = i.toDTO(&r)
	}
	return &intrv1.GetByIdsResponse{
		Intrs: intrs,
	}, nil
}

// IncrReadCnt
func (i *InteractiveSeviceAdapter) IncrReadCnt(ctx context.Context, in *intrv1.IncrReadCntRequest, opts ...grpc.CallOption) (*intrv1.IncrReadCntResponse, error) {
	err := i.svc.IncrReadCnt(ctx, in.Biz, in.BizId)
	return &intrv1.IncrReadCntResponse{}, err
}

// Like
func (i *InteractiveSeviceAdapter) Like(ctx context.Context, in *intrv1.LikeRequest, opts ...grpc.CallOption) (*intrv1.LikeResponse, error) {
	err := i.svc.Like(ctx, in.Biz, in.BizId, in.UserId)
	return &intrv1.LikeResponse{}, err
}

var _ intrv1.InteractiveServiceClient = (*InteractiveSeviceAdapter)(nil)

func NewInteractiveServiceAdapter(svc service.InteractiveService) *InteractiveSeviceAdapter {
	return &InteractiveSeviceAdapter{
		svc: svc,
	}
}

func (i *InteractiveSeviceAdapter) toDTO(intr *domain.Interactive) *intrv1.Interactive {
	return &intrv1.Interactive{
		Id:         intr.Id,
		ReadCnt:    intr.ReadCnt,
		CollectCnt: intr.CollectCnt,
		LikeCnt:    intr.LikeCnt,
		Liked:      intr.Liked,
		Collected:  intr.Collectd,
	}
}
