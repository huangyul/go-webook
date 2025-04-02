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
	SyncStatus(ctx context.Context, userId, id int64, status uint8) error
	GetByAuthorId(ctx context.Context, userId, page, pageSize int64) ([]*Article, error)
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

func (dao *GormArticleDAO) GetByAuthorId(ctx context.Context, userId, page, pageSize int64) ([]*Article, error) {
	var arts []Article
	err := dao.db.WithContext(ctx).Model(&Article{}).Where("author_id = ?", userId).Limit(int(pageSize)).Offset(int((page - 1) * pageSize)).Find(&arts).Error
	if err != nil {
		return nil, err
	}
	result := make([]*Article, 0, len(arts))
	for _, art := range arts {
		result = append(result, &art)
	}
	return result, nil
}

func (dao *GormArticleDAO) SyncStatus(ctx context.Context, userId, id int64, status uint8) error {
	now := time.Now()
	return dao.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		res := tx.Model(&Article{}).Where("id = ? AND author_id = ?", id, userId).Updates(map[string]interface{}{
			"status":     status,
			"created_at": now,
		})
		if res.Error != nil {
			return res.Error
		}
		if res.RowsAffected == 0 {
			return ErrArticleNotFound
		}
		return tx.Model(&PubArticle{}).Where("id = ? AND author_id = ?", id, userId).Updates(map[string]interface{}{
			"status":     status,
			"created_at": now,
		}).Error
	})
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
		pubArt := PubArticle(*art)
		return tx.WithContext(ctx).Model(&PubArticle{}).Clauses(clause.OnConflict{
			Columns: []clause.Column{{Name: "id"}},
			DoUpdates: clause.Assignments(map[string]interface{}{
				"updated_at": now,
				"content":    pubArt.Content,
				"title":      pubArt.Title,
				"status":     pubArt.Status,
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
		"status":     art.Status,
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
	Status    uint8
	CreatedAt time.Time
	UpdatedAt time.Time
}

type PubArticle Article
