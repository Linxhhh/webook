package ioc

import (
	"strings"
	"time"

	"github.com/Linxhhh/webook/internal/app/middleware"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func InitMiddleware() []gin.HandlerFunc {
	return []gin.HandlerFunc{
		// 注册鉴权中间件
		middleware.AuthByJWT(),

		// 配置 CORS
		cors.New(cors.Config{
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
		}),
	}
}

/* 注册 Session 会话中间件
store, err := redis.NewStore(16, "tcp", "localhost:6379", "", []byte("sgpLG7yh8mUYnh619gO0P5HdYftPKpAQ"), []byte("FlIESLxvbN5wiYZS6v7HgLkqsTmED0yh"))
if err != nil {
	panic(err)
}
router.Use(sessions.Sessions("ssid", store))
*/
