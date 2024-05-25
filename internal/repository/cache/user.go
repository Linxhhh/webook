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

type UserCache struct {
	cmd       redis.Cmdable
	expiresAt time.Duration
}

func NewUserCache(cmd redis.Cmdable) *UserCache {
	return &UserCache{
		cmd:       cmd,
		expiresAt: time.Minute * 15,
	}
}

func (uc *UserCache) Key(id int64) string {
	// 格式化 Key
	return fmt.Sprintf("user:info:%d", id)
}

func (uc *UserCache) Set(ctx context.Context, u domain.User) error {
	
	// 序列化 -> []byte
	val, err := json.Marshal(u)
	if err != nil {
		return err
	}
	key := uc.Key(u.Id)

	// 存储 kv
	return uc.cmd.Set(key, val, uc.expiresAt).Err()
}

func (uc *UserCache) Get(ctx context.Context, id int64) (domain.User, error) {
	
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