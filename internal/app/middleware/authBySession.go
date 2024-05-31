package middleware

import (
	"net/http"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

/*
AuthBySession 鉴权中间件：
基于 Session 对用户登录进行校验，并对其刷新 Session 有效期
*/
func AuthBySession() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		path := ctx.Request.URL.Path
		for _, ignorePath := range ignorePaths {
			// 不需要登录鉴权
			if path == ignorePath {
				return
			}
		}

		// 获取 Session
		session := sessions.Default(ctx)
		id := session.Get("userId")
		if id == nil {
			// 如果用户未登录，拦截请求
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// 刷新有效期
		now := time.Now().UnixMilli()
		val := session.Get("update_time")
		if val == nil {
			// 第一次刷新有效期
			session.Set("update_time", now)
			session.Set("userId", id)
			session.Options(sessions.Options{MaxAge: 3600})
			if err := session.Save(); err != nil {
				println(err)
			}
		} else {
			// 十分钟刷新一次有效期
			lastUpdateTime := val.(int64)
			if now - lastUpdateTime >= 600*1000 {
				session.Set("update_time", now)
				session.Set("userId", id)
				session.Options(sessions.Options{MaxAge: 3600})
				if err := session.Save(); err != nil {
					println(err)
				}
			}
		}
	}
}