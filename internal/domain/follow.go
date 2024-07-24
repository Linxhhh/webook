package domain

type FollowData struct {
	Id        int64
	Uid       int64
	Followers int64 // 粉丝数量
	Followees int64 // 关注数量
	Ctime     int64
	Utime     int64

	// 上面数据是每位用户的公共数据
	// 下面数据是针对其它粉丝的数据
	IsFollowed bool
}

type FollowRelation struct {
	Id       int64
	Follower int64 // 粉丝
	Followee int64 // 博主
	Ctime    int64
	Utime    int64
}