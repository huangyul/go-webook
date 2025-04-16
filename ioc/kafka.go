package ioc

import (
	"github.com/IBM/sarama"
	"github.com/huangyul/go-webook/internal/events"
	"github.com/huangyul/go-webook/internal/events/article"
	"github.com/huangyul/go-webook/internal/events/history"
	"github.com/spf13/viper"
)

func InitSaramaClient() sarama.Client {
	addr := viper.GetString("kafka.addr")
	if addr == "" {
		panic("kafka addr is empty")
	}
	cfg := sarama.NewConfig()
	cfg.Producer.Return.Successes = true
	cfg.Producer.Return.Errors = true
	client, err := sarama.NewClient([]string{addr}, cfg)
	if err != nil {
		panic(err)
	}
	return client
}

func InitSaramaProducer(client sarama.Client) sarama.SyncProducer {
	producer, err := sarama.NewSyncProducerFromClient(client)
	if err != nil {
		panic(err)
	}
	return producer
}

func InitConsumers(c1 *article.ArticleReadConsumer, c2 *history.Consumer) []events.Consumer {
	return []events.Consumer{c1, c2}
}
