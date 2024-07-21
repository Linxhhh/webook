package domain

import "time"

type Article struct {
	Id       int64         `json:"id"`
	Title    string        `json:"title"`
	Content  string        `json:"content"`
	AuthorId int64         `json:"authorId"`
	Status   ArticleStatus `json:"status"`
	Ctime    time.Time     `json:"ctime"`
	Utime    time.Time     `json:"utime"`

	// 需要通过 AuthorId 查询，只有查询线上库时，才显示 AuthorName
	AuthorName string `json:"authorName"`
}

// 帖子列表
type ArticleListElem struct {
	Id       int64         `json:"id"`
	Title    string        `json:"title"`
	Abstract string        `json:"abstract"`
	Status   ArticleStatus `json:"status"`
	Ctime    time.Time     `json:"ctime"`
	Utime    time.Time     `json:"utime"`
}

// 帖子状态
type ArticleStatus uint8

const (
	ArticleStatusUnpublished = iota // 未发表
	ArticleStatusPublished          // 已发表
	ArticleStatusPrivate            // 私有
)

// 获取文章内容摘要
func Abstract(content string) string {
	str := []rune(content)
	if len(str) > 128 {
		str = str[:128]
	}
	return string(str)
}
