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
	SyncStatus(ctx context.Context, uid int64, id int64, status uint8) error
	ListByAuthor(ctx context.Context, uid int64, page int64, pageSize int64) ([]Article, error)
	GetById(ctx context.Context, uid int64, id int64) (Article, error)
	GetPubById(ctx context.Context, uid int64, id int64) (PublishedArticle, error)
}

var _ ArticleDao = (*GormArticleDao)(nil)

type GormArticleDao struct {
	db *gorm.DB
}

func NewArticleDao(db *gorm.DB) ArticleDao {
	return &GormArticleDao{db: db}
}

func (dao *GormArticleDao) GetPubById(ctx context.Context, uid int64, id int64) (PublishedArticle, error) {
	var art PublishedArticle
	err := dao.db.WithContext(ctx).Where("id = ? AND author_id = ?", id, uid).First(&art).Error
	if err != nil {
		return PublishedArticle{}, err
	}
	return art, nil
}

func (dao *GormArticleDao) ListByAuthor(ctx context.Context, uid int64, page int64, pageSize int64) ([]Article, error) {
	var res []Article
	err := dao.db.WithContext(ctx).Where("author_id = ?", uid).Offset(int((page - 1) * pageSize)).Limit(int(pageSize)).Find(&res).Error
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (dao *GormArticleDao) GetById(ctx context.Context, uid int64, id int64) (Article, error) {
	var res Article
	err := dao.db.WithContext(ctx).Where("id = ? AND author_id = ?", id, uid).Find(&res).Error
	if err != nil {
		return Article{}, err
	}
	return res, nil
}

func (dao *GormArticleDao) SyncStatus(ctx context.Context, uid int64, id int64, status uint8) error {
	now := time.Now().UnixMilli()
	return dao.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		res := tx.Model(&Article{}).Where("id = ? AND author_id = ?", id, uid).Updates(map[string]interface{}{
			"status":     status,
			"updated_at": now,
		})
		if res.Error != nil {
			return res.Error
		}
		if res.RowsAffected == 0 {
			return gorm.ErrRecordNotFound
		}
		return tx.Model(&PublishedArticle{}).Where("id = ?", id).Updates(map[string]interface{}{
			"status":     status,
			"updated_at": now,
		}).Error
	})
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
				"status":     pubArt.Status,
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
		"status":     art.Status,
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
	Status    uint8
	CreatedAt int64
	UpdatedAt int64
}

type PublishedArticle Article
