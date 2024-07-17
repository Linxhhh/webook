package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Linxhhh/webook/internal/domain"
	"github.com/go-redis/redis"
)

const (
	ErrKeyNotExist = redis.Nil
)

type UserCache interface {
	Set(ctx context.Context, u domain.User) error
	Get(ctx context.Context, id int64) (domain.User, error)
	Del(ctx context.Context, id int64) error
}

type RedisUserCache struct {
	cmd       redis.Cmdable
	expiresAt time.Duration
}

func NewUserCache(cmd redis.Cmdable) UserCache {
	return &RedisUserCache{
		cmd:       cmd,
		expiresAt: time.Minute * 15,
	}
}

func (uc *RedisUserCache) Key(id int64) string {
	// 格式化 Key
	return fmt.Sprintf("user:info:%d", id)
}

func (uc *RedisUserCache) Set(ctx context.Context, u domain.User) error {
	
	// 序列化 -> []byte
	val, err := json.Marshal(u)
	if err != nil {
		return err
	}
	key := uc.Key(u.Id)

	// 存储 kv
	return uc.cmd.Set(key, val, uc.expiresAt).Err()
}

func (uc *RedisUserCache) Get(ctx context.Context, id int64) (domain.User, error) {
	
	// 获取 kv
	key := uc.Key(id)
	val, err := uc.cmd.Get(key).Result()
	if err != nil {
		return domain.User{}, err
	}
	
	// 反序列化 -> domain.User
	var u domain.User
	err = json.Unmarshal([]byte(val), &u)
	return u, err
}

func (uc *RedisUserCache) Del(ctx context.Context, id int64) error {
	
	// 获取 kv
	key := uc.Key(id)
	return uc.cmd.Del(key).Err()
}