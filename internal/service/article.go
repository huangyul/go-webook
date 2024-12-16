package service

import (
	"context"
	"github.com/huangyul/go-blog/internal/domain"
	"github.com/huangyul/go-blog/internal/event/article"
	"github.com/huangyul/go-blog/internal/pkg/log"
	"github.com/huangyul/go-blog/internal/repository"
)

type ArticleService interface {
	Save(ctx context.Context, art domain.Article) (int64, error)
	Publish(ctx context.Context, art domain.Article) (int64, error)
	Withdraw(ctx context.Context, uid, id int64) error
	List(ctx context.Context, uid int64, page int64, pageSize int64) ([]domain.Article, error)
	Detail(ctx context.Context, uid int64, id int64) (domain.Article, error)
	PubDetail(ctx context.Context, uid int64, id int64, biz string) (domain.Article, error)
}

var _ ArticleService = (*articleService)(nil)

type articleService struct {
	repo     repository.ArticleRepository
	userRepo repository.UserRepository
	producer article.Producer
	l        log.Logger
}

func NewArticleService(repo repository.ArticleRepository, userRepo repository.UserRepository, producer article.Producer, l log.Logger) ArticleService {
	return &articleService{
		repo:     repo,
		userRepo: userRepo,
		producer: producer,
		l:        l,
	}
}

func (svc *articleService) PubDetail(ctx context.Context, uid int64, id int64, biz string) (domain.Article, error) {
	art, err := svc.repo.GetPubById(ctx, uid, id)
	if err != nil {
		return domain.Article{}, err
	}
	user, err := svc.userRepo.FindById(ctx, uid)
	if err != nil {
		return domain.Article{}, err
	}
	art.Author.Name = user.Nickname
	go func() {
		er := svc.producer.ProduceReadEvent(article.ReadEvent{
			ArticleID: art.ID,
			UserID:    user.ID,
			Biz:       biz,
		})
		if er != nil {
			log.Errorw("article produce read event error", "id", art.ID, "err", er)
		}
	}()
	return art, nil
}

func (svc *articleService) Detail(ctx context.Context, uid int64, id int64) (domain.Article, error) {
	art, err := svc.repo.GetById(ctx, uid, id)
	if err != nil {
		return domain.Article{}, err
	}
	user, err := svc.userRepo.FindById(ctx, uid)
	if err != nil {
		return domain.Article{}, err
	}
	art.Author.Name = user.Nickname
	return art, nil
}

func (svc *articleService) Withdraw(ctx context.Context, uid int64, id int64) error {
	return svc.repo.SyncStatus(ctx, uid, id, domain.ArticleStatusWithdraw)
}

func (svc *articleService) Save(ctx context.Context, art domain.Article) (int64, error) {
	art.Status = domain.ArticleStatusUnPublished
	if art.ID == 0 {
		return svc.repo.Create(ctx, art)
	} else {
		err := svc.repo.Update(ctx, art)
		return art.ID, err
	}
}

func (svc *articleService) Publish(ctx context.Context, art domain.Article) (int64, error) {
	art.Status = domain.ArticleStatusPublished
	return svc.repo.Sync(ctx, art)
}

func (svc *articleService) List(ctx context.Context, uid int64, page int64, pageSize int64) ([]domain.Article, error) {
	return svc.repo.ListByAuthor(ctx, uid, page, pageSize)
}
