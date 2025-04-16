package dao

import (
	"context"
	"gorm.io/gorm"
	"time"
)

type HistoryDAO interface {
	Insert(ctx context.Context, data *History) error
	ListByUserId(ctx context.Context, userId int64) (history []*History, err error)
}

type GormHistoryDAO struct {
	db *gorm.DB
}

func NewHistoryDAO(db *gorm.DB) HistoryDAO {
	return &GormHistoryDAO{db: db}
}

func (dao *GormHistoryDAO) Insert(ctx context.Context, data *History) error {
	now := time.Now()
	data.CreatedAt = now
	data.UpdatedAt = now
	return dao.db.WithContext(ctx).Create(data).Error
}

func (dao *GormHistoryDAO) ListByUserId(ctx context.Context, userId int64) (history []*History, err error) {
	var result []*History
	err = dao.db.WithContext(ctx).Where("user_id = ?", userId).Find(&result).Error
	return result, err
}

type History struct {
	Id        int64 `gorm:"primary_key;AUTO_INCREMENT"`
	ArticleId int64
	UserId    int64
	CreatedAt time.Time
	UpdatedAt time.Time
}
