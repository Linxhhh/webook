package ioc

import (
	"github.com/Linxhhh/webook/internal/repository/dao"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitDB() *gorm.DB {
	db, err := gorm.Open(mysql.Open("root:root@tcp(localhost:13316)/webook"))
	if err != nil {
		panic(err)
	}

	err = db.AutoMigrate(
		// 用户模块
		&dao.User{},

		// 帖子模块
		&dao.Article{},
		&dao.PublishedArticle{},
		&dao.Interaction{},
		&dao.UserLike{},
		&dao.UserCollection{},
		&dao.FollowData{},
		&dao.FollowRelation{},
	)
	if err != nil {
		panic(err)
	}
	return db
}