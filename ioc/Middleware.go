package ioc

import (
	"net/http"
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
			ExposeHeaders:    []string{"jwt-token", "Content-Length", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers", "Content-Type"},
			// AllowAllOrigins:  true,
			AllowOriginFunc: func(origin string) bool {
				// 允许开发环境的 localhost 和 127.0.0.1
				if strings.HasPrefix(origin, "http://localhost") || strings.HasPrefix(origin, "https://localhost") ||
					strings.HasPrefix(origin, "http://127.0.0.1") || strings.HasPrefix(origin, "https://127.0.0.1") {
					return true
				}
				return strings.Contains(origin, "webook.com")
			},
			MaxAge: 12 * time.Hour,
		}),

		// 处理 Options
		handleOptions(),
	}
}

// 处理 OPTIONS 请求
func handleOptions() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}

/* 注册 Session 会话中间件
store, err := redis.NewStore(16, "tcp", "localhost:6379", "", []byte("sgpLG7yh8mUYnh619gO0P5HdYftPKpAQ"), []byte("FlIESLxvbN5wiYZS6v7HgLkqsTmED0yh"))
if err != nil {
	panic(err)
}
router.Use(sessions.Sessions("ssid", store))
*/
