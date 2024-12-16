package history

import (
	"encoding/json"
	"github.com/IBM/sarama"
	"github.com/huangyul/go-blog/internal/pkg/log"
)

const topic = "history_add"

type Event struct {
	ArticleID int64
	UserID    int64
}

type Producer interface {
	ProduceHistoryEvent(evt Event) error
}

type SaramaProducer struct {
	client sarama.SyncProducer
	l      log.Logger
}

func NewSaramaProducer(client sarama.SyncProducer, l log.Logger) Producer {
	return &SaramaProducer{client: client, l: l}
}

func (s *SaramaProducer) ProduceHistoryEvent(evt Event) error {
	data, err := json.Marshal(evt)
	if err != nil {
		s.l.Errorw("json marshal fail", "err", err)
		return err
	}
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.ByteEncoder(data),
	}
	_, _, err = s.client.SendMessage(msg)
	if err != nil {
		s.l.Errorw("send fail", "err", err)
		return err
	}
	return nil
}
