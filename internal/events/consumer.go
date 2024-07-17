package events

import (

	"github.com/IBM/sarama"
	"github.com/Linxhhh/webook/internal/repository"
)

type Consumer interface {
	Start() error
}

type ReadEventConsumer struct {
	repo   repository.InteractionRepository
	client sarama.Client
}

func NewReadEventConsumer(repo repository.InteractionRepository, client sarama.Client) *ReadEventConsumer {
	return &ReadEventConsumer{
		repo: repo, 
		client: client,
	}
}

/*
func (i *ReadEventConsumer) Start() error {
	cg, err := sarama.NewConsumerGroupFromClient("interaction  ", i.client)
	if err != nil {
		return err
	}
	go func() {
		er := cg.Consume(context.Background(),
			[]string{TopicReadEvent},
			samarax.NewBatchHandler[ReadEvent](i.l, i.BatchConsume))
		if er != nil {
			log.Print("退出消费", er)
		}
	}()
	return err
}

func (i *ReadEventConsumer) StartV1() error {
	cg, err := sarama.NewConsumerGroupFromClient("interactive", i.client)
	if err != nil {
		return err
	}
	go func() {
		er := cg.Consume(context.Background(),
			[]string{TopicReadEvent},
			samarax.NewHandler[ReadEvent](i.l, i.Consume))
		if er != nil {
			i.l.Error("退出消费", logger.Error(er))
		}
	}()
	return err
}
func (i *ReadEventConsumer) BatchConsume(msgs []*sarama.ConsumerMessage,
	events []ReadEvent) error {
	bizs := make([]string, 0, len(events))
	bizIds := make([]int64, 0, len(events))
	for _, evt := range events {
		bizs = append(bizs, "article")
		bizIds = append(bizIds, evt.Aid)
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	return i.repo.BatchIncrReadCnt(ctx, bizs, bizIds)
}

func (i *ReadEventConsumer) Consume(msg *sarama.ConsumerMessage,
	event ReadEvent) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	return i.repo.IncrReadCnt(ctx, "article", event.Aid)
}
*/