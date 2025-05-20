package grpc

import (
	"context"

	intrv1 "github.com/huangyul/go-webook/api/proto/gen/intr/v1"
	"github.com/huangyul/go-webook/interactive/service"
	"google.golang.org/grpc"
)

var _ intrv1.InteractiveServiceServer = (*InteractiveService)(nil)

type InteractiveService struct {
	intrv1.UnimplementedInteractiveServiceServer
	svc service.InteractiveService
}

func NewInteractiveService(svc service.InteractiveService) *InteractiveService {
	return &InteractiveService{
		svc: svc,
	}
}

func (s *InteractiveService) Register(grpcS *grpc.Server) {
	intrv1.RegisterInteractiveServiceServer(grpcS, s)
}

func (s *InteractiveService) CancelCollect(ctx context.Context, req *intrv1.CancelCollectRequest) (*intrv1.CancelCollectResponse, error) {
	err := s.svc.CancelCollect(ctx, req.Biz, req.BizId, req.UserId)
	return &intrv1.CancelCollectResponse{}, err
}
func (s *InteractiveService) CancelLike(ctx context.Context, req *intrv1.CancelLikeRequest) (*intrv1.CancelLikeResponse, error) {
	err := s.svc.CancelLike(ctx, req.Biz, req.BizId, req.UserId)
	return &intrv1.CancelLikeResponse{}, err
}
func (s *InteractiveService) Collect(ctx context.Context, req *intrv1.CollectRequest) (*intrv1.CollectResponse, error) {
	err := s.svc.Collect(ctx, req.Biz, req.BizId, req.UserId)
	return &intrv1.CollectResponse{}, err
}
func (s *InteractiveService) Get(ctx context.Context, req *intrv1.GetRequest) (*intrv1.GetResponse, error) {
	res, err := s.svc.Get(ctx, req.Biz, req.BizId, req.UserId)
	if err != nil {
		return &intrv1.GetResponse{}, err
	}
	return &intrv1.GetResponse{
		Intr: &intrv1.Interactive{
			Id:         res.Id,
			CollectCnt: res.CollectCnt,
			LikeCnt:    res.LikeCnt,
			ReadCnt:    res.ReadCnt,
			Collected:  res.Collectd,
			Liked:      res.Liked,
		},
	}, nil
}
func (s *InteractiveService) GetByIds(ctx context.Context, req *intrv1.GetByIdsRequest) (*intrv1.GetByIdsResponse, error) {
	res, err := s.svc.GetByIds(ctx, req.Biz, req.Ids)
	if err != nil {
		return &intrv1.GetByIdsResponse{}, err
	}
	intrs := make(map[int64]*intrv1.Interactive)
	for _, r := range res {
		intrs[r.Id] = &intrv1.Interactive{
			Id:         r.Id,
			CollectCnt: r.CollectCnt,
			LikeCnt:    r.LikeCnt,
			ReadCnt:    r.ReadCnt,
			Collected:  r.Collectd,
			Liked:      r.Liked,
		}
	}
	return &intrv1.GetByIdsResponse{
		Intrs: intrs,
	}, nil
}
func (s *InteractiveService) IncrReadCnt(ctx context.Context, req *intrv1.IncrReadCntRequest) (*intrv1.IncrReadCntResponse, error) {
	err := s.svc.IncrReadCnt(ctx, req.Biz, req.BizId)
	return &intrv1.IncrReadCntResponse{}, err
}
func (s *InteractiveService) Like(ctx context.Context, req *intrv1.LikeRequest) (*intrv1.LikeResponse, error) {
	err := s.svc.Like(ctx, req.Biz, req.BizId, req.UserId)
	return &intrv1.LikeResponse{}, err
}
