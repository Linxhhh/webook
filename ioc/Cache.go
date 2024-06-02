package ioc

import (
	"context"
	"log"
	"time"

	"github.com/go-redis/redis"
)

func InitCache() redis.Cmdable {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	// 设置超时时间
	_, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	// 测试 Redis 连接是否正常
	_, err := client.Ping().Result()
	if err != nil {
		log.Fatalf("Redis 连接失败，%s", err.Error())
	}
	return client
}
