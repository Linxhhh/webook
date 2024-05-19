package middleware

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		path := ctx.Request.URL.Path
		if path == "/user/signup" || path == "/user/login" {
			// 不需要登录鉴权
			return
		}

		// 获取 Session
		session := sessions.Default(ctx)
		if session.Get("userId") == nil {
			// 如果用户未登录，拦截请求
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
	}
}