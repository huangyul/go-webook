package service

import (
	"context"
	"log"
	"time"

	"github.com/huangyul/go-webook/internal/domain"
	"github.com/huangyul/go-webook/internal/events/article"
	"github.com/huangyul/go-webook/internal/events/history"
	"github.com/huangyul/go-webook/internal/repository"
	"golang.org/x/sync/errgroup"
)

type ArticleService interface {
	Save(ctx context.Context, art *domain.Article) (int64, error)
	Publish(ctx context.Context, art *domain.Article) error
	WithDraw(ctx context.Context, userId, id int64) error
	GetByAuthor(ctx context.Context, userId, page, pageSize int64) ([]*domain.Article, error)
	GetById(ctx context.Context, id int64, userId int64) (*domain.Article, error)
	GetPudDetailById(ctx context.Context, id int64) (*domain.Article, error)
	GetPudById(ctx context.Context, id int64, userId int64, biz string) (*domain.Article, error)
	ListPub(ctx context.Context, start time.Time, offset int, size int) ([]*domain.Article, error)
}

type articleService struct {
	repo            repository.ArticleRepository
	userRepo        repository.UserRepository
	producer        article.ReadProducer
	historyProducer history.Producer
}

// ListPub implements ArticleService.
func (svc *articleService) ListPub(ctx context.Context, start time.Time, offset int, size int) ([]*domain.Article, error) {
	return svc.repo.ListPub(ctx, start, offset, size)
}

func (svc *articleService) GetPudDetailById(ctx context.Context, id int64) (*domain.Article, error) {
	return svc.repo.GetPudDetailById(ctx, id)
}

func (svc *articleService) GetById(ctx context.Context, id int64, userId int64) (*domain.Article, error) {
	return svc.repo.GetById(ctx, id, userId)
}

func (svc *articleService) GetPudById(ctx context.Context, id int64, userId int64, biz string) (*domain.Article, error) {
	art, err := svc.repo.GetPubById(ctx, id, userId)
	go func() {
		if err == nil {
			var eg errgroup.Group
			eg.Go(func() error {
				return svc.producer.Produce(&article.ReadEvent{
					UserId: userId,
					ArtId:  id,
					Biz:    biz,
				})
			})
			eg.Go(func() error {
				return svc.historyProducer.AddHistory(&history.AddHistoryEvent{
					UserId:    art.Author.Id,
					ArticleId: art.Id,
				})
			})
			er := eg.Wait()
			if er != nil {
				log.Println(er)
			}
		}
	}()
	return art, err
}

func (svc *articleService) GetByAuthor(ctx context.Context, userId, page, pageSize int64) ([]*domain.Article, error) {
	arts, err := svc.repo.GetByAuthorId(ctx, userId, page, pageSize)
	if err != nil {
		return nil, err
	}
	if len(arts) == 0 {
		return []*domain.Article{}, nil
	}
	user, err := svc.userRepo.FindByID(ctx, userId)
	if err == nil && user != nil {
		for _, article := range arts {
			article.Author.Name = user.Nickname
		}
	}

	return arts, nil
}

func (svc *articleService) WithDraw(ctx context.Context, userId, id int64) error {
	return svc.repo.SyncStatus(ctx, userId, id, domain.ArticleStatusPrivate)
}

func (svc *articleService) Save(ctx context.Context, art *domain.Article) (int64, error) {
	art.Status = domain.ArticleStatusUnPublished
	if art.Id > 0 {
		return art.Id, svc.repo.UpdateById(ctx, art)
	} else {
		return svc.repo.Insert(ctx, art)
	}
}

func (svc *articleService) Publish(ctx context.Context, art *domain.Article) error {
	art.Status = domain.ArticleStatusPublished
	return svc.repo.Sync(ctx, art)
}

func NewArticleService(repo repository.ArticleRepository, userRepo repository.UserRepository, producer article.ReadProducer, historyProducer history.Producer) ArticleService {
	return &articleService{
		repo:            repo,
		userRepo:        userRepo,
		producer:        producer,
		historyProducer: historyProducer,
	}
}
