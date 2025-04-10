package article

import (
	"context"
	"fmt"
	"github.com/IBM/sarama"
	"github.com/huangyul/go-webook/internal/events"
	"github.com/huangyul/go-webook/internal/repository"
	"time"
)

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

func (con *ArticleReadConsumer) Start() error {
	cg, err := sarama.NewConsumerGroupFromClient("article", con.client)
	if err != nil {
		return err
	}
	go func() {
		err = cg.Consume(context.Background(), []string{ReadEventTopic}, events.NewHandler[ReadEvent](con.consume))
		if err != nil {
			panic(err)
		}
	}()

	return nil
}

func (con *ArticleReadConsumer) consume(evt ReadEvent) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err := con.repo.IncrReadCnt(ctx, evt.Biz, evt.ArtId)
	if err != nil {
		fmt.Printf("Error incrementing read cnt: %v\n", err)
	}
}
