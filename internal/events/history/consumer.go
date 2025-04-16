package history

import (
	"context"
	"github.com/IBM/sarama"
	"github.com/huangyul/go-webook/internal/domain"
	"github.com/huangyul/go-webook/internal/events"
	"github.com/huangyul/go-webook/internal/repository"
	"log"
)

type Consumer struct {
	client sarama.Client
	repo   repository.HistoryRepository
}

func NewConsumer(client sarama.Client, repo repository.HistoryRepository) *Consumer {
	return &Consumer{
		client: client,
		repo:   repo,
	}
}

func (c *Consumer) Start() error {
	cg, err := sarama.NewConsumerGroupFromClient("history", c.client)
	if err != nil {
		return err
	}

	go func() {
		cg.Consume(context.Background(), []string{AddHistoryTopic}, events.NewHandler[AddHistoryEvent](func(evt AddHistoryEvent) {
			c.consume(evt)
		}))
	}()
	return nil
}

func (c *Consumer) consume(evt AddHistoryEvent) {
	err := c.repo.Insert(context.Background(), &domain.History{
		AuthorId:  evt.UserId,
		ArticleId: evt.ArticleId,
	})
	if err != nil {
		log.Println(err)
	}
}
