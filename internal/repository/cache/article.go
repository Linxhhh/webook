package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Linxhhh/webook/internal/domain"
	"github.com/go-redis/redis"
)

type ArticleCache interface {
	GetFirstPage(ctx context.Context, uid int64) ([]domain.ArticleListElem, error)
	SetFirstPage(ctx context.Context, uid int64, arts []domain.ArticleListElem) error
	DelFirstPage(ctx context.Context, uid int64) error
	Get(ctx context.Context, id int64) (domain.Article, error)
	Set(ctx context.Context, art domain.Article) error
	GetPub(ctx context.Context, id int64) (domain.Article, error)
	SetPub(ctx context.Context, art domain.Article) error
}

type RedisArticleCache struct {
	cmd redis.Cmdable
}

func NewArticleCache(cmd redis.Cmdable) ArticleCache {
	return &RedisArticleCache{
		cmd: cmd,
	}
}

/*
缓存制作库的帖子列表首页
- listKey
- GetFirstPage
- SetFirstPage
- DelFirstPage
*/

func (ac *RedisArticleCache) listKey(uid int64) string {
	return fmt.Sprintf("article:list:%d", uid)
}

func (ac *RedisArticleCache) GetFirstPage(ctx context.Context, uid int64) ([]domain.ArticleListElem, error) {
	
	// 获取 kv
	key := ac.listKey(uid)
	val, err := ac.cmd.Get(key).Bytes()
	if err != nil {
		return nil, err
	}
	
	// 反序列化 -> []domain.ArticleListElem
	var arts []domain.ArticleListElem
	err = json.Unmarshal([]byte(val), &arts)
	return arts, err
}

func (ac *RedisArticleCache) SetFirstPage(ctx context.Context, uid int64, arts []domain.ArticleListElem) error {
	
	// 设置 kv
	key := ac.listKey(uid)
	val, err := json.Marshal(arts)
	if err != nil {
		return err
	}
	return ac.cmd.Set(key, val, 10 * time.Minute).Err()
}

func (ac *RedisArticleCache) DelFirstPage(ctx context.Context, uid int64) error {
	
	// 删除 kv
	key := ac.listKey(uid)
	return ac.cmd.Del(key).Err()
}

/*
缓存制作库的帖子详情
- key
- Get
- Set
*/

func (ac *RedisArticleCache) key(id int64) string {
	return fmt.Sprintf("article:detail:%d", id)
}

func (ac *RedisArticleCache) Get(ctx context.Context, id int64) (domain.Article, error) {
	
	// 获取 kv
	key := ac.key(id)
	val, err := ac.cmd.Get(key).Bytes()
	if err != nil {
		return domain.Article{}, err
	}
	
	// 反序列化 -> domain.Article
	var art domain.Article
	err = json.Unmarshal([]byte(val), &art)
	return art, err
}

func (ac *RedisArticleCache) Set(ctx context.Context, art domain.Article) error {
	
	// 设置 kv
	key := ac.key(art.Id)
	val, err := json.Marshal(art)
	if err != nil {
		return err
	}
	return ac.cmd.Set(key, val, 10 * time.Minute).Err()
}

/*
缓存线上库的帖子详情
- pubKey
- GetPub
- SetPub
*/

func (ac *RedisArticleCache) pubKey(id int64) string {
	return fmt.Sprintf("article:pub:%d", id)
}

func (ac *RedisArticleCache) GetPub(ctx context.Context, id int64) (domain.Article, error) {
	
	// 获取 kv
	key := ac.pubKey(id)
	val, err := ac.cmd.Get(key).Bytes()
	if err != nil {
		return domain.Article{}, err
	}
	
	// 反序列化 -> domain.Article
	var art domain.Article
	err = json.Unmarshal([]byte(val), &art)
	return art, err
}

func (ac *RedisArticleCache) SetPub(ctx context.Context, art domain.Article) error {
	
	// 设置 kv
	key := ac.pubKey(art.Id)
	val, err := json.Marshal(art)
	if err != nil {
		return err
	}
	return ac.cmd.Set(key, val, 10 * time.Minute).Err()
}