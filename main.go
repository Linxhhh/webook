package main

import (
	"strings"
	"time"

	"github.com/Linxhhh/webook/internal/app"
	"github.com/Linxhhh/webook/internal/app/middleware"
	"github.com/Linxhhh/webook/internal/repository"
	"github.com/Linxhhh/webook/internal/repository/dao"
	"github.com/Linxhhh/webook/internal/service"
	"github.com/gin-contrib/cors"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	// 初始化数据库
	db := initDB()

	// 初始化路由
	router := initRouter()
	initUserRouter(db, router)
	router.Run()
}

func initRouter() *gin.Engine {
	router := gin.Default()

	/* 注册会话中间件
	store, err := redis.NewStore(16, "tcp", "localhost:6379", "", []byte("sgpLG7yh8mUYnh619gO0P5HdYftPKpAQ"), []byte("FlIESLxvbN5wiYZS6v7HgLkqsTmED0yh"))
	if err != nil {
		panic(err)
	}
	router.Use(sessions.Sessions("ssid", store))
	*/

	// 注册鉴权中间件
	router.Use(middleware.AuthByJWT())

	// 配置 CORS
	router.Use(cors.New(cors.Config{
		AllowCredentials: true,
		AllowHeaders:     []string{"Content-Type", "jwt-token"},
		ExposeHeaders:    []string{"jwt-token"},
		AllowOriginFunc: func(origin string) bool {
			// 开发环境下允许
			if strings.HasPrefix(origin, "http://localhost") {
				return true
			}
			return strings.Contains(origin, "webook.com")
		},
		MaxAge: 12 * time.Hour,
	}))

	return router
}

func initUserRouter(db *gorm.DB, router *gin.Engine) {
	ud := dao.NewUserDAO(db)
	ur := repository.NewUserRepository(ud)
	us := service.NewUserService(ur)
	hdl := app.NewUserHandler(us)
	hdl.RegistryRouter(router)
}

func initDB() *gorm.DB {
	db, err := gorm.Open(mysql.Open("root:root@tcp(localhost:13316)/webook"))
	if err != nil {
		panic(err)
	}

	err = db.AutoMigrate(&dao.User{})
	if err != nil {
		panic(err)
	}
	return db
}
