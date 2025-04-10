package events

import (
	"encoding/json"
	"github.com/IBM/sarama"
)

type Handler[T any] struct {
	fn func(evt T)
}

func NewHandler[T any](fn func(evt T)) *Handler[T] {
	return &Handler[T]{
		fn: fn,
	}
}

func (h Handler[T]) Setup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (h Handler[T]) Cleanup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (h Handler[T]) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	msgs := claim.Messages()
	for msg := range msgs {
		var evt T
		if err := json.Unmarshal(msg.Value, &evt); err != nil {
			return err
		}
		h.fn(evt)
		session.MarkMessage(msg, "")
	}
	return nil
}
