package ioc

import (
	"github.com/IBM/sarama"
	"github.com/huangyul/go-blog/interactive/events"
	events2 "github.com/huangyul/go-blog/internal/event"
	"github.com/spf13/viper"
)

func InitSaramaClient() sarama.Client {
	saramaCfg := sarama.NewConfig()
	saramaCfg.Producer.Return.Successes = true
	saramaCfg.Producer.Return.Errors = true
	client, err := sarama.NewClient([]string{viper.GetString("kafka.addr")}, saramaCfg)
	if err != nil {
		panic(err)
	}
	return client
}

func InitConsumers(c1 *events.InteractiveReadEventConsumer) []events2.Consumer {
	return []events2.Consumer{
		c1,
	}
}
