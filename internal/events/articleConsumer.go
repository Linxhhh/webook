package events

import (
	"context"
	"log"
	"strconv"
	"time"

	"github.com/IBM/sarama"
	"github.com/Linxhhh/webook/internal/domain"
	"github.com/Linxhhh/webook/internal/service"
	samarax "github.com/Linxhhh/webook/pkg/saramax"
)

type ArticleEventConsumer struct {
	client sarama.Client
	svc    *service.FeedEventService
}

func NewArticleEventConsumer(client sarama.Client, svc *service.FeedEventService) *ArticleEventConsumer {
	return &ArticleEventConsumer{
		svc:    svc,
		client: client,
	}
}

// Start 启动 goroutine 消费事件
func (r *ArticleEventConsumer) Start() error {

	cg, err := sarama.NewConsumerGroupFromClient("articleFeed", r.client)
	if err != nil {
		return err
	}

	go func() {
		err := cg.Consume(context.Background(), []string{TopicArticleEvent}, samarax.NewConsumer[ArticleEvent](r.Consume))
		if err != nil {
			log.Println("退出了消费循环异常", err)
		}
	}()
	return err
}

// Consume 消费 ArticleEvent
func (r *ArticleEventConsumer) Consume(msg *sarama.ConsumerMessage, evt ArticleEvent) error {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	return r.svc.CreateFeedEvent(ctx, domain.FeedEvent{
		Type: TopicArticleEvent,
		Ext: map[string]string{
			"uid":   strconv.FormatInt(evt.Uid, 10),
			"aid":   strconv.FormatInt(evt.Aid, 10),
			"title": evt.Title,
		},
	})
}
