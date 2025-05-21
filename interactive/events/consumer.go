package events

import (
	"context"
	"fmt"

	"github.com/IBM/sarama"
	"github.com/huangyul/go-webook/interactive/repository"
	"github.com/huangyul/go-webook/pkg/saramax"
)

const readEventTopic = "article_read"

type ArticleReadConsumer struct {
	client sarama.Client
	repo   repository.InteractiveRepository
}

func NewArticleReadConsumer(client sarama.Client, repo repository.InteractiveRepository) *ArticleReadConsumer {
	return &ArticleReadConsumer{
		client: client,
		repo:   repo,
	}
}

func (c *ArticleReadConsumer) Start() error {
	consumer, err := sarama.NewConsumerGroupFromClient("interactive", c.client)
	if err != nil {
		panic(err)
	}

	go func() {
		consumer.Consume(context.Background(), []string{readEventTopic}, saramax.NewHandler(func(evt ReadEvent) {
			err := c.repo.IncrReadCnt(context.Background(), evt.Biz, evt.ArtId)
			if err != nil {
				fmt.Printf("interactive:consumer error: %s", err.Error())
			}
		}))
	}()

	return err
}

type ReadEvent struct {
	ArtId int64
	Biz   string
}
