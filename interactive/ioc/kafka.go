package ioc

import (
	"github.com/IBM/sarama"
	intrConsumer "github.com/huangyul/go-webook/interactive/events"
	"github.com/huangyul/go-webook/internal/events"
	"github.com/spf13/viper"
)

func InitSaramaClient() sarama.Client {
	addr := viper.GetString("kafka.addr")
	cfg := sarama.NewConfig()
	client, err := sarama.NewClient([]string{addr}, cfg)
	if err != nil {
		panic(err)
	}
	return client
}

func InitConsumers(c1 *intrConsumer.ArticleReadConsumer) []events.Consumer {
	return []events.Consumer{
		c1,
	}
}
