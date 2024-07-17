package domain

import "time"

type Article struct {
	Id       int64
	Title    string
	Content  string
	AuthorId int64
	Status   ArticleStatus
	Ctime    time.Time
	Utime    time.Time

	// 需要通过 AuthorId 查询，只有查询线上库时，才显示 AuthorName
	AuthorName string
}

// 帖子列表
type ArticleListElem struct {
	Id       int64
	Title    string
	Abstract string
	Status   ArticleStatus
	Ctime    time.Time
	Utime    time.Time
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


