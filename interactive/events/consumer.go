package events

import (
	"context"
	"github.com/IBM/sarama"
	"github.com/huangyul/go-blog/interactive/repository"
	"github.com/huangyul/go-blog/pkg/saramax"
	"log"
	"time"
)

const topicName = "read_cnt"

type ReadEvent struct {
	ArticleID int64
	UserID    int64
	Biz       string
}

type InteractiveReadEventConsumer struct {
	repo   repository.InteractiveRepository
	client sarama.Client
}

func (c *InteractiveReadEventConsumer) Start() error {
	group, err := sarama.NewConsumerGroupFromClient("read_cnt", c.client)
	if err != nil {
		return err
	}

	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		er := group.Consume(ctx, []string{topicName}, saramax.NewHandler[ReadEvent](c.consume))
		if er != nil {
			log.Println(er)
		}
	}()

	return nil
}

func (c *InteractiveReadEventConsumer) consume(msg *sarama.ConsumerMessage, event ReadEvent) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	return c.repo.IncrReadCnt(ctx, event.ArticleID, event.Biz)

}

func NewInteractiveReadEventConsumer(repo repository.InteractiveRepository, client sarama.Client) *InteractiveReadEventConsumer {
	return &InteractiveReadEventConsumer{repo: repo, client: client}
}
