package grpc

import (
	"context"
	interv1 "github.com/huangyul/go-blog/api/proto/gen/inter/v1"
	"github.com/huangyul/go-blog/interactive/service"
	"google.golang.org/grpc"
)

type InteractiveServiceServer struct {
	interv1.UnimplementedInteractiveServiceServer
	svc service.InteractiveService
}

func NewInteractiveServiceServer(svc service.InteractiveService) *InteractiveServiceServer {
	return &InteractiveServiceServer{svc: svc}
}

func (i *InteractiveServiceServer) Register(srv *grpc.Server) {
	interv1.RegisterInteractiveServiceServer(srv, i)
}

func (i *InteractiveServiceServer) IncrReadCnt(ctx context.Context, request *interv1.IncrReadCntRequest) (*interv1.EmptyResponse, error) {
	err := i.svc.IncrReadCnt(ctx, request.GetBizId(), request.GetBiz())
	return &interv1.EmptyResponse{}, err
}

func (i *InteractiveServiceServer) Like(ctx context.Context, request *interv1.LikeRequest) (*interv1.EmptyResponse, error) {
	err := i.svc.Like(ctx, request.GetUserId(), request.GetBizId(), request.GetBiz())
	return &interv1.EmptyResponse{}, err
}

func (i *InteractiveServiceServer) CancelLike(ctx context.Context, request *interv1.CancelLikeRequest) (*interv1.EmptyResponse, error) {
	err := i.svc.CancelLike(ctx, request.GetUserId(), request.GetBizId(), request.GetBiz())
	return &interv1.EmptyResponse{}, err
}

func (i *InteractiveServiceServer) Collect(ctx context.Context, request *interv1.CollectRequest) (*interv1.EmptyResponse, error) {
	err := i.svc.Collect(ctx, request.GetUid(), request.GetId(), request.GetCid(), request.GetBiz())
	return &interv1.EmptyResponse{}, err
}

func (i *InteractiveServiceServer) Get(ctx context.Context, request *interv1.GetRequest) (*interv1.InteractiveResponse, error) {
	res, err := i.svc.Get(ctx, request.GetUid(), request.GetId(), request.GetBiz())
	if err != nil {
		return &interv1.InteractiveResponse{}, err
	}
	return &interv1.InteractiveResponse{
		ReadCnt:    int32(res.ReadCnt),
		Liked:      res.Liked,
		Collected:  res.Collected,
		LikeCnt:    int32(res.LikeCnt),
		CollectCnt: int32(res.CollectCnt),
	}, nil
}
