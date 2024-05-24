package main

import (
	"context"
	"log"
	"strings"
	"time"

	"github.com/Linxhhh/webook/internal/app"
	"github.com/Linxhhh/webook/internal/app/middleware"
	"github.com/Linxhhh/webook/internal/repository"
	"github.com/Linxhhh/webook/internal/repository/cache"
	"github.com/Linxhhh/webook/internal/repository/dao"
	"github.com/Linxhhh/webook/internal/service"
	"github.com/gin-contrib/cors"
	"github.com/go-redis/redis"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	// 初始化
	db := initDB()
	cache := initCache()
	router := initRouter()
	initUserRouter(db, cache, router)
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

func initUserRouter(db *gorm.DB, cmd *redis.Client, router *gin.Engine) {
	ud := dao.NewUserDAO(db)
	uc := cache.NewUserCache(cmd)
	ur := repository.NewUserRepository(ud, uc)
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

func initCache() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     "locahost:6379",
		Password: "",
		DB:       0,
	})

	// 设置超时时间
	_, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	// 测试 Redis 连接是否正常
	_, err := client.Ping().Result()
	if err != nil {
		log.Fatalf("Redis 连接失败，%s", err.Error())
	}
	return client
}
