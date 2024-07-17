package cache

import (
	"context"
	_ "embed"
	"fmt"
	"strconv"
	"time"

	"github.com/Linxhhh/webook/internal/domain"
	"github.com/go-redis/redis"
)

const (
	fieldFollower = "follower_cnt"
	fieldFollowee = "followee_cnt"
)

type FollowCache interface {
	Get(ctx context.Context, uid int64) (domain.FollowData, error)
	Set(ctx context.Context, data domain.FollowData) error
}

type RedisFollowCache struct {
	cmd       redis.Cmdable
	expiresAt time.Duration
}

func NewFollowCache(cmd redis.Cmdable) FollowCache {
	return &RedisFollowCache{
		cmd:       cmd,
		expiresAt: 15 * time.Minute,
	}
}

func (c *RedisFollowCache) key(uid int64) string {
	return fmt.Sprintf("user:followdata:%d", uid)
}

func (c *RedisFollowCache) Get(ctx context.Context, uid int64) (domain.FollowData, error) {
	key := c.key(uid)
	res, err := c.cmd.HGetAll(key).Result()
	if err != nil {
		return domain.FollowData{}, err
	}

	if len(res) == 0 {
		return domain.FollowData{}, ErrKeyNotExist
	}

	var data domain.FollowData
	data.Followers, _ = strconv.ParseInt(res[fieldFollower], 10, 64)
	data.Followees, _ = strconv.ParseInt(res[fieldFollowee], 10, 64)
	return data, nil
}

func (c *RedisFollowCache) Set(ctx context.Context, data domain.FollowData) error {

	// 设置 kv
	key := c.key(data.Uid)
	if err := c.cmd.HSet(key, fieldFollower, data.Followers).Err(); err != nil {
		return err
	}
	if err := c.cmd.HSet(key, fieldFollowee, data.Followees).Err(); err != nil {
		return err
	}
	return c.cmd.Expire(key, c.expiresAt).Err()
}
