package dao

import (
	"context"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

type InteractiveDao interface {
	IncrReadCnt(ctx context.Context, bizID int64, biz string) error
	InsertLikeBiz(ctx context.Context, userId int64, bizID int64, biz string) error
	DeleteLikeBiz(ctx context.Context, userId int64, bizID int64, biz string) error
}

var _ InteractiveDao = (*GormInteractiveDao)(nil)

type GormInteractiveDao struct {
	db *gorm.DB
}

func NewInteractiveDao(db *gorm.DB) InteractiveDao {
	return &GormInteractiveDao{db: db}
}

func (dao *GormInteractiveDao) InsertLikeBiz(ctx context.Context, userId int64, bizID int64, biz string) error {
	now := time.Now().UnixMilli()
	return dao.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		err := tx.Clauses(
			clause.OnConflict{
				DoUpdates: clause.Assignments(map[string]interface{}{
					"status":     1,
					"updated_at": now,
				}),
			}).Create(&UserLikeBiz{
			BizID:     bizID,
			UserID:    userId,
			Biz:       biz,
			Status:    1,
			CreatedAt: now,
			UpdatedAt: now,
		}).Error
		if err != nil {
			return err
		}
		return tx.Clauses(clause.OnConflict{
			DoUpdates: clause.Assignments(map[string]interface{}{
				"like_cnt":   gorm.Expr("like_cnt + ?", 1),
				"updated_at": now,
			}),
		}).Create(&Interactive{
			BizID:     bizID,
			Biz:       biz,
			CreatedAt: now,
			UpdatedAt: now,
			LikeCnt:   1,
		}).Error
	})
}

func (dao *GormInteractiveDao) DeleteLikeBiz(ctx context.Context, userId int64, bizID int64, biz string) error {
	now := time.Now().UnixMilli()
	return dao.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		err := tx.Model(&UserLikeBiz{}).Where("user_id = ? AND biz_id = ? AND biz = ?", userId, bizID, biz).Updates(map[string]interface{}{
			"status":     0,
			"updated_at": now,
		}).Error
		if err != nil {
			return err
		}
		return tx.Model(&Interactive{}).Where("biz = ? AND biz_id = ?", biz, bizID).Updates(map[string]interface{}{
			"updated_at": now,
			"like_cnt":   gorm.Expr("like_cnt - ?", 1),
		}).Error
	})
}

func (dao *GormInteractiveDao) IncrReadCnt(ctx context.Context, bizID int64, biz string) error {
	now := time.Now().UnixMilli()
	return dao.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns: []clause.Column{
			{Name: "id"},
		},
		DoUpdates: clause.Assignments(map[string]interface{}{
			"updated_at": now,
			"read_cnt":   gorm.Expr("read_cnt + ?", 1),
		}),
	}).Create(&Interactive{
		BizID:     bizID,
		Biz:       biz,
		ReadCnt:   1,
		CreatedAt: now,
		UpdatedAt: now,
	}).Error
}

type Interactive struct {
	ID         int64  `gorm:"primary_key;AUTO_INCREMENT"`
	BizID      int64  `gorm:"uniqueIndex:biz_type_id"`
	Biz        string `gorm:"type:varchar(128);uniqueIndex:biz_type_id"`
	ReadCnt    int
	LikeCnt    int
	CollectCnt int
	CreatedAt  int64
	UpdatedAt  int64
}

type UserLikeBiz struct {
	ID        int64  `gorm:"primary_key;AUTO_INCREMENT"`
	UserID    int64  `gorm:"uniqueIndex:user_biz_type_id"`
	BizID     int64  `gorm:"uniqueIndex:user_biz_type_id"`
	Biz       string `gorm:"type:varchar(128);uniqueIndex:user_biz_type_id"`
	Status    int    `gorm:"comment:0-取消 1-喜爱"`
	CreatedAt int64
	UpdatedAt int64
}
