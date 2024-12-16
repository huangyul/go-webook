package ioc

import (
	"github.com/IBM/sarama"
	"github.com/huangyul/go-blog/internal/event"
	"github.com/huangyul/go-blog/internal/event/article"
	"github.com/huangyul/go-blog/internal/event/history"
	"github.com/spf13/viper"
)

func InitSaramaClient() sarama.Client {
	addr := viper.GetString("kafka.addr")
	if addr == "" {
		panic("kafka addr is empty")
	}
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	client, err := sarama.NewClient([]string{addr}, config)
	if err != nil {
		panic(err)
	}
	return client
}

func InitProducer(client sarama.Client) sarama.SyncProducer {
	producer, err := sarama.NewSyncProducerFromClient(client)
	if err != nil {
		panic(err)
	}
	return producer
}

func InitConsumers(c1 *article.InteractiveReadConsumer, c2 *history.Consumer) []event.Consumer {
	return []event.Consumer{c1, c2}
}
