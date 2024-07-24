package events

import (
	"encoding/json"

	"github.com/IBM/sarama"
)

type ArticleEvent struct {
	Uid   int64
	Aid   int64
	Title string
}

type ArticleEventProducer struct {
	producer sarama.SyncProducer
}

func NewArticleEventProducer(producer sarama.SyncProducer) *ArticleEventProducer {
	return &ArticleEventProducer{producer: producer}
}

func (s *ArticleEventProducer) ProduceEvent(evt ArticleEvent) error {
	val, err := json.Marshal(evt)
	if err != nil {
		return err
	}
	_, _, err = s.producer.SendMessage(&sarama.ProducerMessage{
		Topic: TopicArticleEvent,
		Value: sarama.StringEncoder(val),
	})
	return err
}
