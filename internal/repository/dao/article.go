package dao

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

type ArticleDAO interface {
	Insert(ctx context.Context, art *Article) (int64, error)
	UpdateById(ctx context.Context, art *Article) error
	Sync(ctx context.Context, art *Article) error
}

var (
	ErrArticleNotFound = errors.New("article not found")
)

func NewArticleDAO(db *gorm.DB) ArticleDAO {
	return &GormArticleDAO{
		db: db,
	}
}

type GormArticleDAO struct {
	db *gorm.DB
}

func (dao *GormArticleDAO) Sync(ctx context.Context, art *Article) error {
	err := dao.db.Transaction(func(tx *gorm.DB) error {
		d := NewArticleDAO(tx)
		var er error
		if art.Id > 0 {
			er = d.UpdateById(ctx, art)
		} else {
			_, er = d.Insert(ctx, art)
		}
		if er != nil {
			return er
		}
		now := time.Now()
		pubArt := PubArticle{
			Title:     art.Title,
			Content:   art.Content,
			CreatedAt: now,
			UpdatedAt: now,
		}
		return tx.WithContext(ctx).Model(&PubArticle{}).Clauses(clause.OnConflict{
			Columns: []clause.Column{{Name: "id"}},
			DoUpdates: clause.Assignments(map[string]interface{}{
				"updated_at": now,
				"content":    pubArt.Content,
				"title":      pubArt.Title,
			}),
		}).Create(art).Error
	})
	return err
}

func (dao *GormArticleDAO) Insert(ctx context.Context, art *Article) (int64, error) {
	now := time.Now()
	art.CreatedAt = now
	art.UpdatedAt = now
	err := dao.db.WithContext(ctx).Create(art).Error
	if err != nil {
		return 0, err
	}
	return art.Id, nil
}

func (dao *GormArticleDAO) UpdateById(ctx context.Context, art *Article) error {
	now := time.Now()
	art.UpdatedAt = now
	res := dao.db.WithContext(ctx).Model(art).Where("id = ? AND author_id = ?", art.Id, art.AuthorId).Updates(map[string]interface{}{
		"updated_at": now,
		"title":      art.Title,
		"content":    art.Content,
	})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return ErrArticleNotFound
	}
	return nil
}

type Article struct {
	Id        int64 `gorm:"primary_key;auto_increment"`
	Title     string
	Content   string `gorm:"type:BLOB"`
	AuthorId  int64  `gorm:"column:author_id"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type PubArticle Article
