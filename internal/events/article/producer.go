package article

import (
	"encoding/json"
	"github.com/IBM/sarama"
)

const ReadEventTopic = "article_read"

type ReadProducer interface {
	Produce(evt *ReadEvent) error
}

type ReadEvent struct {
	ArtId  int64
	UserId int64
	Biz    string
}

type ArticleReadProducer struct {
	producer sarama.SyncProducer
}

func NewArticleReadProducer(producer sarama.SyncProducer) ReadProducer {
	return &ArticleReadProducer{producer: producer}
}

func (p *ArticleReadProducer) Produce(evt *ReadEvent) error {
	data, err := json.Marshal(evt)
	if err != nil {
		return err
	}
	_, _, err = p.producer.SendMessage(&sarama.ProducerMessage{
		Topic: ReadEventTopic,
		Value: sarama.ByteEncoder(data),
	})
	return err
}
