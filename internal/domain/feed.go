package domain

import (
	"errors"
	"fmt"
	"time"
)

type FeedEvent struct {
	Id    int64
	Uid   int64
	Type  string
	Ctime time.Time
	Ext   ExtendFields
}

const (
	ArticleFeedEvent = "article_feed_event"
	ReadFeedEvent    = "read_feed_event"
	LikeFeedEvent    = "like_feed_event"
	CollectFeedEvent = "coll_feed_event"
)

/*
拓展字段，Feed 应该可以推送帖子、点赞消息、收藏消息、关注消息等。
*/
type ExtendFields map[string]string

var errKeyNotFound = errors.New("没有找到对应的 key")

func (f ExtendFields) Get(key string) (val string, err error) {
	val, ok := f[key]
	if !ok {
		return "", fmt.Errorf("%w, key %s", errKeyNotFound, key)
	}
	return val, nil
}
