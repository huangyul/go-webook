package dao

import (
	"context"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

type InteractiveDAO interface {
	IncrReadCnt(ctx context.Context, biz string, bizId int64) error
	InsertLikeInfo(ctx context.Context, biz string, bizId int64, userId int64) error
	DeleteLikeInfo(ctx context.Context, biz string, bizId int64, userId int64) error
}

func NewInteractiveDAO(db *gorm.DB) InteractiveDAO {
	return &GormInteractiveDAO{
		db: db,
	}
}

type GormInteractiveDAO struct {
	db *gorm.DB
}

func (dao *GormInteractiveDAO) InsertLikeInfo(ctx context.Context, biz string, bizId int64, userId int64) error {
	now := time.Now()
	return dao.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		err := tx.Clauses(clause.OnConflict{
			DoUpdates: clause.Assignments(map[string]interface{}{
				"updated_at": now,
				"status":     1,
			}),
		}).Create(&UserLikeBiz{
			Biz:       biz,
			BizId:     bizId,
			UserId:    userId,
			Status:    1,
			CreatedAt: now,
			UpdatedAt: now,
		}).Error
		if err != nil {
			return err
		}
		return tx.Clauses(clause.OnConflict{
			DoUpdates: clause.Assignments(map[string]interface{}{
				"updated_at": now,
				"like_cnt":   gorm.Expr("like_cnt + ?", 1),
			}),
		}).Create(&Interactive{
			Biz:       biz,
			BizId:     bizId,
			LikeCnt:   1,
			CreatedAt: now,
			UpdatedAt: now,
		}).Error
	})
}

func (dao *GormInteractiveDAO) DeleteLikeInfo(ctx context.Context, biz string, bizId int64, userId int64) error {
	now := time.Now()
	return dao.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		err := tx.Model(&UserLikeBiz{}).Where("biz_id = ? and user_id = ? and biz = ?", bizId, userId, biz).Updates(map[string]interface{}{
			"updated_at": now,
			"status":     0,
		}).Error
		if err != nil {
			return err
		}
		return tx.Model(&Interactive{}).Where("biz_id = ? and biz = ?", bizId, biz).Updates(map[string]any{
			"updated_at": now,
			"like_cnt":   gorm.Expr("like_cnt - ?", 1),
		}).Error
	})
}

func (dao *GormInteractiveDAO) IncrReadCnt(ctx context.Context, biz string, bizId int64) error {
	now := time.Now()
	return dao.db.WithContext(ctx).Clauses(clause.OnConflict{
		DoUpdates: clause.Assignments(map[string]interface{}{
			"updated_at": now,
			"read_cnt":   gorm.Expr("read_cnt + ?", 1),
		}),
	}).Create(&Interactive{
		Biz:       biz,
		BizId:     bizId,
		ReadCnt:   1,
		CreatedAt: now,
		UpdatedAt: now,
	}).Error
}

type Interactive struct {
	Id         int64     `gorm:"primary_key;auto_increment"`
	BizId      int64     `gorm:"column:biz_id;uniqueIndex:idx_biz_id_biz;"`
	Biz        string    `gorm:"column:biz;uniqueIndex:idx_biz_id_biz;type:varchar(255)"`
	ReadCnt    int64     `gorm:"column:read_cnt"`
	LikeCnt    int64     `gorm:"column:like_cnt"`
	CollectCnt int64     `gorm:"column:collect_cnt"`
	CreatedAt  time.Time `gorm:"column:created_at"`
	UpdatedAt  time.Time `gorm:"column:updated_at"`
}

type UserLikeBiz struct {
	Id        int64     `gorm:"primary_key;auto_increment"`
	BizId     int64     `gorm:"column:biz_id;uniqueIndex:idx_biz_id_user_id_biz;"`
	UserId    int64     `gorm:"column:user_id;uniqueIndex:idx_biz_id_user_id_biz;"`
	Biz       string    `gorm:"column:biz;uniqueIndex:idx_biz_id_user_id_biz;type:varchar(128)"`
	Status    int       `gorm:"column:status;common:'0/取消点赞；1/点赞'"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}
