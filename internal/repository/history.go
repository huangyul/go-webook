package repository

import (
	"context"
	"github.com/huangyul/go-webook/internal/domain"
	"github.com/huangyul/go-webook/internal/repository/dao"
)

type HistoryRepository interface {
	Insert(ctx context.Context, h *domain.History) error
	ListByUserId(ctx context.Context, userId int64) ([]*domain.History, error)
}

type HistoryRepositoryImpl struct {
	dao dao.HistoryDAO
}

func NewHistoryRepository(dao dao.HistoryDAO) HistoryRepository {
	return &HistoryRepositoryImpl{dao: dao}
}

func (repo *HistoryRepositoryImpl) Insert(ctx context.Context, h *domain.History) error {
	return repo.dao.Insert(ctx, repo.toEntity(h))
}

func (repo *HistoryRepositoryImpl) ListByUserId(ctx context.Context, userId int64) ([]*domain.History, error) {
	res, err := repo.dao.ListByUserId(ctx, userId)
	if err != nil {
		return nil, err
	}
	datas := make([]*domain.History, 0)
	for _, v := range res {
		datas = append(datas, repo.toDomain(v))
	}
	return datas, nil
}

func (repo *HistoryRepositoryImpl) toDomain(data *dao.History) *domain.History {
	return &domain.History{
		Id:        data.Id,
		AuthorId:  data.UserId,
		ArticleId: data.ArticleId,
		CreatedAt: data.CreatedAt,
		UpdatedAt: data.UpdatedAt,
	}
}

func (repo *HistoryRepositoryImpl) toEntity(h *domain.History) *dao.History {

	return &dao.History{
		Id:        h.Id,
		UserId:    h.AuthorId,
		ArticleId: h.ArticleId,
	}
}
