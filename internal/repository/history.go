package repository

import (
	"context"
	"github.com/huangyul/go-blog/internal/domain"
	"github.com/huangyul/go-blog/internal/repository/dao"
	"time"
)

type HistoryRepository interface {
	GetListByUser(ctx context.Context, userID int64) ([]domain.History, error)
	Create(ctx context.Context, history domain.History) error
}

type historyRepository struct {
	dao dao.HistoryDao
}

func NewHistoryRepository(dao dao.HistoryDao) HistoryRepository {
	return &historyRepository{dao: dao}
}

func (h *historyRepository) Create(ctx context.Context, history domain.History) error {
	return h.dao.Insert(ctx, h.toEntity(history))
}

func (h *historyRepository) GetListByUser(ctx context.Context, userID int64) ([]domain.History, error) {
	hs, err := h.dao.GetByUser(ctx, userID)
	if err != nil {
		return nil, err
	}
	result := make([]domain.History, 0, len(hs))
	for _, hh := range hs {
		result = append(result, h.toDomain(hh))
	}
	return result, nil
}

func (h *historyRepository) toDomain(his dao.History) domain.History {
	return domain.History{
		BizID:     his.BizID,
		UserID:    his.UserID,
		UpdatedAt: time.UnixMilli(his.UpdatedAt),
	}
}

func (h *historyRepository) toEntity(his domain.History) dao.History {
	return dao.History{
		BizID:  his.BizID,
		UserID: his.UserID,
	}
}
