package saramax

import (
	"encoding/json"

	"github.com/IBM/sarama"
)

var _ sarama.ConsumerGroupHandler = (*Handler[any])(nil)

type Handler[T any] struct {
	fn func(evt T)
}

func NewHandler[T any](fn func(evt T)) *Handler[T] {
	return &Handler[T]{
		fn: fn,
	}
}

func (h *Handler[T]) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

func (h *Handler[T]) Setup(sarama.ConsumerGroupSession) error {
	return nil
}

func (h *Handler[T]) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	msgs := claim.Messages()

	for msg := range msgs {
		var evt T
		err := json.Unmarshal(msg.Value, &evt)
		if err != nil {
			return err
		}
		h.fn(evt)
		session.MarkMessage(msg, "")
	}
	return nil
}
