package domain

// 用户互动数据表
type Interaction struct {
	Id          int64  `json:"id"`
	Biz         string `json:"biz"`
	BizId       int64  `json:"bizId"`
	ReadCnt     int64  `json:"readCnt"`
	LikeCnt     int64  `json:"likeCnt"`
	CollectCnt  int64  `json:"collectCnt"`

	// 上面数据是一篇帖子的公共数据
	// 下面数据是针对具体用户的数据
	IsLiked     bool   `json:"isLiked"`
	IsCollected bool   `json:"isCollected"`
}
