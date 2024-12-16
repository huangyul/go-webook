package article

import (
	"context"
	"github.com/IBM/sarama"
	"github.com/huangyul/go-blog/internal/pkg/log"
	"github.com/huangyul/go-blog/internal/repository"
	"github.com/huangyul/go-blog/pkg/saramax"
	"time"
)

type InteractiveReadConsumer struct {
	client sarama.Client
	repo   repository.InteractiveRepository
	l      log.Logger
}

func NewInteractiveReadConsumer(client sarama.Client, repo repository.InteractiveRepository, l log.Logger) *InteractiveReadConsumer {
	return &InteractiveReadConsumer{client: client, repo: repo, l: l}
}

// Start get message to consume
func (i *InteractiveReadConsumer) Start() error {
	cg, err := sarama.NewConsumerGroupFromClient("article", i.client)
	if err != nil {
		i.l.Errorw("create consumer group error", "error", err)
		return err
	}
	go func() {
		err = cg.Consume(context.Background(), []string{topicName}, saramax.NewHandler[ReadEvent](i.consumeMessage))
		if err != nil {
			i.l.Errorw("consume message error", "error", err)
		}
	}()

	return nil
}

// consumeMessage use repo
func (i *InteractiveReadConsumer) consumeMessage(msg *sarama.ConsumerMessage, evt ReadEvent) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	return i.repo.IncrReadCnt(ctx, evt.ArticleID, evt.Biz)
}
