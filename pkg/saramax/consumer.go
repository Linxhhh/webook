package samarax

import (
	"encoding/json"
	"log"

	"github.com/IBM/sarama"
)

type Consumer[T any] struct {
	fn func(msg *sarama.ConsumerMessage, event T) error
}

func NewConsumer[T any](fn func(msg *sarama.ConsumerMessage, event T) error) *Consumer[T] {
	return &Consumer[T]{fn: fn}
}

func (h *Consumer[T]) Setup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (h *Consumer[T]) Cleanup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (h *Consumer[T]) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	
	msgs := claim.Messages()
	for msg := range msgs {
		var t T
		err := json.Unmarshal(msg.Value, &t)
		if err != nil {
			log.Println("反序列消息体失败:")
			log.Println("topic", msg.Topic)
			log.Println("partition", msg.Partition)
			log.Println("offset", msg.Offset)
			log.Println(err)
		}
		err = h.fn(msg, t)
		if err != nil {
			log.Println("处理消息失败")
			log.Println("topic", msg.Topic)
			log.Println("partition", msg.Partition)
			log.Println("offset", msg.Offset)
			log.Println(err)
		}
		session.MarkMessage(msg, "")
	}
	return nil
}