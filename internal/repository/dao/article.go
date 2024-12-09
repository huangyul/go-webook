package dao

import (
	"context"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

type ArticleDao interface {
	Insert(ctx context.Context, art Article) (int64, error)
	UpdateByID(ctx context.Context, art Article) error
	Sync(ctx context.Context, art Article) (int64, error)
}

var _ ArticleDao = (*GormArticleDao)(nil)

type GormArticleDao struct {
	db *gorm.DB
}

func NewArticleDao(db *gorm.DB) ArticleDao {
	return &GormArticleDao{db: db}
}

func (dao *GormArticleDao) Sync(ctx context.Context, art Article) (int64, error) {
	id := art.ID
	err := dao.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var err error
		dao := NewArticleDao(tx)
		if art.ID > 0 {
			err = dao.UpdateByID(ctx, art)
		} else {
			id, err = dao.Insert(ctx, art)
		}
		if err != nil {
			return err
		}
		art.ID = id
		pubArt := PublishedArticle(art)
		pubArt.UpdatedAt = time.Now().UnixMilli()
		pubArt.CreatedAt = time.Now().UnixMilli()
		return tx.Clauses(clause.OnConflict{
			Columns: []clause.Column{
				{Name: "id"},
			},
			DoUpdates: clause.Assignments(map[string]interface{}{
				"title":      pubArt.Title,
				"content":    pubArt.Content,
				"updated_at": pubArt.UpdatedAt,
			}),
		}).Create(pubArt).Error
	})
	return id, err
}

func (dao *GormArticleDao) UpdateByID(ctx context.Context, art Article) error {
	now := time.Now().UnixMilli()
	res := dao.db.WithContext(ctx).Model(&art).Where("id = ? AND author_id = ?", art.ID, art.AuthorID).Updates(map[string]interface{}{
		"title":      art.Title,
		"content":    art.Content,
		"updated_at": now,
	})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (dao *GormArticleDao) Insert(ctx context.Context, art Article) (int64, error) {
	now := time.Now().UnixMilli()
	art.CreatedAt = now
	art.UpdatedAt = now
	err := dao.db.WithContext(ctx).Create(&art).Error
	return art.ID, err
}

type Article struct {
	ID        int64  `gorm:"primary_key;AUTO_INCREMENT"`
	Title     string `gorm:"type:varchar(4096);"`
	Content   string `gorm:"type:text;"`
	AuthorID  int64  `gorm:"index"`
	CreatedAt int64
	UpdatedAt int64
}

type PublishedArticle Article
