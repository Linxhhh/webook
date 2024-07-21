package ioc

import (
	"github.com/Linxhhh/webook/internal/repository/dao"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitDB() (master *gorm.DB, slaves []*gorm.DB) {
	master, err := gorm.Open(mysql.Open("root:123456@tcp(localhost:13306)/webook"))
	if err != nil {
		panic(err)
	}

	s1, err := gorm.Open(mysql.Open("root:123456@tcp(localhost:23306)/webook"))
	if err != nil {
		panic(err)
	}

	slaves = append(slaves, s1)

	err = master.AutoMigrate(
		&dao.User{},
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

	return
}