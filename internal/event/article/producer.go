package article

import (
	"encoding/json"
	"github.com/IBM/sarama"
)

const topicName = "article_read"

type ReadEvent struct {
	ArticleID int64
	UserID    int64
	Biz       string
}

type Producer interface {
	ProduceReadEvent(evt ReadEvent) error
}

type SaramaSyncProducer struct {
	client sarama.SyncProducer
}

func NewSaramaSyncProducer(client sarama.SyncProducer) Producer {
	return &SaramaSyncProducer{client: client}
}

func (s *SaramaSyncProducer) ProduceReadEvent(ev ReadEvent) error {
	data, err := json.Marshal(ev)
	if err != nil {
		return err
	}
	msg := &sarama.ProducerMessage{
		Value: sarama.ByteEncoder(data),
		Topic: topicName,
	}
	_, _, err = s.client.SendMessage(msg)
	return err
}
