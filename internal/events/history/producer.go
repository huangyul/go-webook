package history

import (
	"encoding/json"
	"github.com/IBM/sarama"
)

const AddHistoryTopic = "history_add"

type AddHistoryEvent struct {
	UserId    int64
	ArticleId int64
}

type Producer interface {
	AddHistory(*AddHistoryEvent) error
}

type ProducerImpl struct {
	producer sarama.SyncProducer
}

func NewHistoryProducer(producer sarama.SyncProducer) Producer {
	return &ProducerImpl{
		producer: producer,
	}
}

func (producer *ProducerImpl) AddHistory(event *AddHistoryEvent) error {
	data, err := json.Marshal(event)
	if err != nil {
		return err
	}
	_, _, err = producer.producer.SendMessage(&sarama.ProducerMessage{
		Topic: AddHistoryTopic,
		Value: sarama.ByteEncoder(data),
	})
	return err
}
