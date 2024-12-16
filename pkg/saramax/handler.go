package saramax

import (
	"encoding/json"
	"fmt"
	"github.com/IBM/sarama"
)

type Handler[T any] struct {
	fn func(msg *sarama.ConsumerMessage, event T) error
}

func NewHandler[T any](fn func(msg *sarama.ConsumerMessage, event T) error) *Handler[T] {
	return &Handler[T]{fn: fn}
}

func (h *Handler[T]) Setup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (h *Handler[T]) Cleanup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (h *Handler[T]) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	msgs := claim.Messages()
	for msg := range msgs {
		var t T
		err := json.Unmarshal(msg.Value, &t)
		if err != nil {
			fmt.Println("sarama consume json unmarshal error:", err)
			session.MarkMessage(msg, "")
			continue
		}
		err = h.fn(msg, t)
		if err != nil {
			fmt.Println("sarama consume error:", err)
		}
		session.MarkMessage(msg, "")
	}
	return nil
}
