package dao

import (
	"context"
	"errors"
	"sync/atomic"
)

const (
	PatternSrcFirst = "src_first"
	PatternSrcOnly  = "src_only"
	PatternDstFirst = "dst_first"
	PatternDstOnly  = "dst_only"
)

var errUnknownPattern = errors.New("unknown pattern")

type DoubleWriteDao struct {
	src     InteractiveDao
	dst     InteractiveDao
	pattern atomic.Value
}

func NewDoubleWriteDao(src InteractiveDao, dst InteractiveDao) *DoubleWriteDao {
	return &DoubleWriteDao{src: src, dst: dst}
}

func (dao *DoubleWriteDao) UpdatePattern(p string) {
	dao.pattern.Store(p)
}

func (dao *DoubleWriteDao) IncrReadCnt(ctx context.Context, bizID int64, biz string) error {
	switch dao.pattern.Load() {
	case PatternSrcOnly:
		return dao.src.IncrReadCnt(ctx, bizID, biz)
	case PatternSrcFirst:
		err := dao.src.IncrReadCnt(ctx, bizID, biz)
		if err != nil {
			return err
		}
		return dao.dst.IncrReadCnt(ctx, bizID, biz)
	case PatternDstFirst:
		err := dao.dst.IncrReadCnt(ctx, bizID, biz)
		if err != nil {
			return err
		}
		return dao.src.IncrReadCnt(ctx, bizID, biz)
	case PatternDstOnly:
		return dao.dst.IncrReadCnt(ctx, bizID, biz)
	default:
		return errUnknownPattern
	}
}

func (dao *DoubleWriteDao) InsertLikeBiz(ctx context.Context, userId int64, bizID int64, biz string) error {
	//TODO implement me
	panic("implement me")
}

func (dao *DoubleWriteDao) DeleteLikeBiz(ctx context.Context, userId int64, bizID int64, biz string) error {
	//TODO implement me
	panic("implement me")
}

func (dao *DoubleWriteDao) InsertCollectBiz(ctx context.Context, uid int64, id int64, cid int64, biz string) error {
	//TODO implement me
	panic("implement me")
}

func (dao *DoubleWriteDao) Get(ctx context.Context, id int64, biz string) (Interactive, error) {
	//TODO implement me
	panic("implement me")
}

func (dao *DoubleWriteDao) GetLikedInfo(ctx context.Context, uid int64, id int64, biz string) (UserLikeBiz, error) {
	//TODO implement me
	panic("implement me")
}

func (dao *DoubleWriteDao) GetCollectInfo(ctx context.Context, uid int64, id int64, biz string) (UserCollectionBiz, error) {
	//TODO implement me
	panic("implement me")
}
