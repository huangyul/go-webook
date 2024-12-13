package sarama

import (
	"fmt"
	"github.com/IBM/sarama"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

var (
	Addr      = []string{"localhost:9094"}
	TopicName = "test_topic"
)

func TestSarama_SyncProducer(t *testing.T) {
	cfg := sarama.NewConfig()
	cfg.Producer.Return.Successes = true
	producer, err := sarama.NewSyncProducer(Addr, cfg)
	cfg.Producer.Partitioner = sarama.NewRoundRobinPartitioner
	assert.NoError(t, err)
	defer producer.Close()
	for i := 0; i < 10; i++ {
		msg := &sarama.ProducerMessage{
			Topic: TopicName,
			Value: sarama.StringEncoder(fmt.Sprintf("这是第%d条消息", i+1)),
		}
		_, _, err = producer.SendMessage(msg)
		assert.NoError(t, err)
	}
}

func TestSarama_AsyncProducer(t *testing.T) {
	cfg := sarama.NewConfig()
	cfg.Producer.Return.Successes = true
	cfg.Producer.Return.Errors = true
	cfg.Producer.Partitioner = sarama.NewRoundRobinPartitioner
	producer, err := sarama.NewAsyncProducer(Addr, cfg)
	assert.NoError(t, err)
	defer producer.Close()
	for i := 0; i < 10; i++ {
		msg := &sarama.ProducerMessage{
			Topic: TopicName,
			Value: sarama.StringEncoder(fmt.Sprintf("这是异步消息，第%d条", i+1)),
		}
		producer.Input() <- msg
	}

	for {
		select {
		case err := <-producer.Errors():
			fmt.Println(err)
		case msg := <-producer.Successes():
			fmt.Println(msg.Topic, msg.Value)
		case <-time.After(time.Second * 5):
			return
		}
	}

}
