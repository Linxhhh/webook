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

//go:embed lua/incrCnt.lua
var luaIncrCnt string

const fieldReadCnt = "read_cnt"
const fieldLikeCnt = "like_cnt"
const fieldCollectCnt = "collect_cnt"

type InteractionCache interface {
	IncrReadCnt(ctx context.Context, biz string, bizId int64) error
	IncrLikeCnt(ctx context.Context, biz string, bizId int64) error
	DecrLikeCnt(ctx context.Context, biz string, bizId int64) error
	Get(ctx context.Context, biz string, bizId int64) (domain.Interaction, error)
	Set(ctx context.Context, biz string, bizId int64, interaction domain.Interaction) error
}

type RedisInteractionCache struct {
	cmd       redis.Cmdable
	expiresAt time.Duration
}

func NewInteractionCache(cmd redis.Cmdable) InteractionCache {
	return &RedisInteractionCache{
		cmd:       cmd,
		expiresAt: time.Minute * 15,
	}
}

func (i *RedisInteractionCache) key(biz string, bizId int64) string {
	return fmt.Sprintf("interaction:%s:%d", biz, bizId)
}

func (i *RedisInteractionCache) IncrReadCnt(ctx context.Context, biz string, bizId int64) error {
	key := i.key(biz, bizId)
	return i.cmd.Eval(luaIncrCnt, []string{key}, fieldReadCnt, 1).Err()
}

func (i *RedisInteractionCache) IncrLikeCnt(ctx context.Context, biz string, bizId int64) error {
	key := i.key(biz, bizId)
	return i.cmd.Eval(luaIncrCnt, []string{key}, fieldLikeCnt, 1).Err()
}

func (i *RedisInteractionCache) DecrLikeCnt(ctx context.Context, biz string, bizId int64) error {
	key := i.key(biz, bizId)
	return i.cmd.Eval(luaIncrCnt, []string{key}, fieldLikeCnt, -1).Err()
}

func (i *RedisInteractionCache) Get(ctx context.Context, biz string, bizId int64) (domain.Interaction, error) {
	key := i.key(biz, bizId)
	res, err := i.cmd.HGetAll(key).Result()
	if err != nil {
		return domain.Interaction{}, err
	}

	if len(res) == 0 {
		return domain.Interaction{}, ErrKeyNotExist
	}
	var ia domain.Interaction
	ia.ReadCnt, _ = strconv.ParseInt(res[fieldReadCnt], 10, 64)
	ia.LikeCnt, _ = strconv.ParseInt(res[fieldLikeCnt], 10, 64)
	ia.CollectCnt, _ = strconv.ParseInt(res[fieldCollectCnt], 10, 64)
	return ia, nil
}

func (i *RedisInteractionCache) Set(ctx context.Context, biz string, bizId int64, ia domain.Interaction) error {

	// 设置 kv
	key := i.key(biz, bizId)
	if err := i.cmd.HSet(key, fieldReadCnt, ia.ReadCnt).Err(); err != nil {
		return err
	}
	if err := i.cmd.HSet(key, fieldLikeCnt, ia.LikeCnt).Err(); err != nil {
		return err
	}
	if err := i.cmd.HSet(key, fieldCollectCnt, ia.CollectCnt).Err(); err != nil {
		return err
	}
	return i.cmd.Expire(key, i.expiresAt).Err()
}
