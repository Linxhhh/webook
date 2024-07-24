package events

import (
	"encoding/json"
	"github.com/IBM/sarama"
)

type ReadEvent struct {
	Aid int64
	Uid int64
}

type BatchReadEvent struct {
	Aids []int64
	Uids []int64
}

type SaramaReadProducer struct {
	producer sarama.SyncProducer
}

func NewSaramaSyncProducer(producer sarama.SyncProducer) *SaramaReadProducer {
	return &SaramaReadProducer{producer: producer}
}

func (s *SaramaReadProducer) ProduceEvent(evt ReadEvent) error {
	val, err := json.Marshal(evt)
	if err != nil {
		return err
	}
	_, _, err = s.producer.SendMessage(&sarama.ProducerMessage{
		Topic: TopicReadEvent,
		Value: sarama.StringEncoder(val),
	})
	return err
}
