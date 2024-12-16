package dao

import (
	"context"
	"gorm.io/gorm"
	"time"
)

type HistoryDao interface {
	GetByUser(ctx context.Context, userId int64) ([]History, error)
	Insert(ctx context.Context, history History) error
}

type GormHistoryDao struct {
	db *gorm.DB
}

func NewHistoryDao(db *gorm.DB) HistoryDao {
	return &GormHistoryDao{db: db}
}

func (dao *GormHistoryDao) Insert(ctx context.Context, history History) error {
	now := time.Now().UnixMilli()
	history.CreatedAt = now
	history.UpdatedAt = now
	return dao.db.WithContext(ctx).Create(&history).Error
}

func (dao *GormHistoryDao) GetByUser(ctx context.Context, userId int64) ([]History, error) {
	var history []History
	err := dao.db.WithContext(ctx).Where("user_id = ?", userId).Find(&history).Error
	return history, err
}

type History struct {
	ID        int64 `gorm:"primary_key;AUTO_INCREMENT"`
	UserID    int64
	BizID     int64
	CreatedAt int64
	UpdatedAt int64
}
