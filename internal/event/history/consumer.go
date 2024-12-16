package history

import (
	"context"
	"github.com/IBM/sarama"
	"github.com/huangyul/go-blog/internal/domain"
	"github.com/huangyul/go-blog/internal/pkg/log"
	"github.com/huangyul/go-blog/internal/repository"
	"github.com/huangyul/go-blog/pkg/saramax"
	"time"
)

type Consumer struct {
	client sarama.Client
	repo   repository.HistoryRepository
	l      log.Logger
}

func NewConsumer(client sarama.Client, repo repository.HistoryRepository, l log.Logger) *Consumer {
	return &Consumer{client: client, repo: repo, l: l}
}

func (c *Consumer) Start() error {
	cg, err := sarama.NewConsumerGroupFromClient("history", c.client)
	if err != nil {
		c.l.Errorw("consumer", "err", err)
		return err
	}
	go func() {
		err = cg.Consume(context.Background(), []string{topic}, saramax.NewHandler[Event](c.consume))
		if err != nil {
			c.l.Errorw("consumer", "err", err)
		}
	}()
	return nil
}

func (c *Consumer) consume(msg *sarama.ConsumerMessage, event Event) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()
	return c.repo.Create(ctx, domain.History{
		UserID: event.UserID,
		BizID:  event.ArticleID,
	})
}
