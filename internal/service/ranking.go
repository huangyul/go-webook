package service

import (
	"context"
	"math"
	"time"

	"github.com/ecodeclub/ekit/queue"
	"github.com/ecodeclub/ekit/slice"
	interactive "github.com/huangyul/go-webook/interactive/service"
	"github.com/huangyul/go-webook/internal/domain"
	"github.com/huangyul/go-webook/internal/repository"
)

type RankingService interface {
	TopN(ctx context.Context) error
	GetTopN(ctx context.Context) ([]domain.Article, error)
}

type RankingServiceImpl struct {
	intrSvc   interactive.InteractiveService
	artSvc    ArticleService
	batchSize int
	scoreFunc func(likeCnt int64, updatedAt time.Time) float64
	n         int
	repo      repository.RankingRepository
}

func (r *RankingServiceImpl) TopN(ctx context.Context) error {
	arts, err := r.topN(ctx)
	if err != nil {
		return err
	}
	return r.repo.ReplaceTopN(ctx, arts)
}

func (r *RankingServiceImpl) GetTopN(ctx context.Context) ([]domain.Article, error) {
	return r.repo.GetTopN(ctx)
}

func (r *RankingServiceImpl) topN(ctx context.Context) ([]domain.Article, error) {
	offset := 0
	start := time.Now()
	ddl := start.Add(-7 * 24 * time.Hour)

	type Score struct {
		score float64
		art   domain.Article
	}

	topN := queue.NewPriorityQueue[Score](r.n, func(src Score, dst Score) int {
		if src.score > dst.score {
			return 1
		} else if src.score == dst.score {
			return 0
		} else {
			return -1
		}
	})

	for {
		arts, err := r.artSvc.ListPub(ctx, start, offset, r.batchSize)
		if err != nil {
			return nil, err
		}

		ids := slice.Map(arts, func(idx int, src *domain.Article) int64 {
			return src.Id
		})
		intrMap, err := r.intrSvc.GetByIds(ctx, "article", ids)
		if err != nil {
			return nil, err
		}

		for _, art := range arts {
			intr := intrMap[art.Id]

			score := r.scoreFunc(intr.LikeCnt, art.UpdatedAt)
			ele := Score{score, *art}
			err = topN.Enqueue(ele)
			if err == queue.ErrOutOfCapacity {
				minEle, _ := topN.Dequeue()
				if minEle.score < score {
					_ = topN.Enqueue(ele)
				} else {
					_ = topN.Enqueue(ele)
				}
			}
		}
		offset = offset + len(arts)
		if len(arts) < r.batchSize ||
			arts[len(arts)-1].UpdatedAt.Before(ddl) {
			break
		}
	}

	res := make([]domain.Article, topN.Len())
	for i := topN.Len() - 1; i >= 0; i-- {
		ele, _ := topN.Dequeue()
		res[i] = ele.art
	}
	return res, nil
}

func NewRankingService(intrSvc interactive.InteractiveService, artSvc ArticleService, repo repository.RankingRepository) RankingService {
	return &RankingServiceImpl{intrSvc: intrSvc, artSvc: artSvc, batchSize: 100, scoreFunc: func(likeCnt int64, updatedAt time.Time) float64 {
		duration := time.Since(updatedAt).Seconds()
		return float64(likeCnt-1) / math.Pow(duration+2, 1.5)
	}, n: 100, repo: repo}
}
