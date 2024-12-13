package sarama

import (
	"context"
	"github.com/IBM/sarama"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
	"time"
)

type consumerHandler struct{}

func (c consumerHandler) Setup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (c consumerHandler) Cleanup(session sarama.ConsumerGroupSession) error {
	return nil
}

// batch consume
func (c consumerHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	batchSize := 10
	msgs := claim.Messages()
	for {
		log.Println("start a new round")
		batch := make([]*sarama.ConsumerMessage, 0, batchSize)
		done := false
		ctx, cancel := context.WithTimeout(context.Background(), time.Minute*10)
		for i := 0; i < batchSize && !done; i++ {
			select {
			case <-ctx.Done():
				done = true
			case msg, ok := <-msgs:
				if !ok {
					cancel()
					return nil
				}
				batch = append(batch, msg)
			}
		}
		cancel()
		for _, msg := range batch {
			log.Println(string(msg.Value))
			session.MarkMessage(msg, "")
		}
	}
	return nil
}

// ConsumeClaimV1 single consume
func (c consumerHandler) ConsumeClaimV1(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	msgs := claim.Messages()
	for msg := range msgs {
		log.Println(string(msg.Value))
		session.MarkMessage(msg, "")
	}
	return nil
}

func TestSarama_Consumer(t *testing.T) {
	config := sarama.NewConfig()
	cg, err := sarama.NewConsumerGroup(Addr, "demo", config)
	assert.NoError(t, err)
	defer cg.Close()

	err = cg.Consume(context.Background(), []string{TopicName}, consumerHandler{})
	assert.NoError(t, err)
	select {}
}
