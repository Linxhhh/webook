package middleware

import (
	"log"
	"time"

	"github.com/Linxhhh/webook/pkg/jwts"
	"github.com/dgrijalva/jwt-go/v4"
	"github.com/gin-gonic/gin"
)

/*
AuthByJWT 鉴权中间件：
*/
func AuthByJWT() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		path := ctx.Request.URL.Path
		for _, ignorePath := range ignorePaths {
			// 不需要登录鉴权
			if path == ignorePath {
				return
			}
		}

		// 获取 Token
		token := ctx.Request.Header.Get("jwt-token")
		if token == "" {
			ctx.String(200, "未携带令牌!")
			ctx.Abort()
			return
		}

		claims, err := jwts.ParseToken(token)
		if err != nil {
			ctx.String(200, "令牌错误!")
			ctx.Abort()
			return
		}
		
		// 对用户代理进行校验
		if claims.UserAgent != ctx.GetHeader("User-Agent") {
			ctx.String(200, "用户代理更改!")
			ctx.Abort()
			return
		}

		// 如果剩余有效期小于30分钟，则刷新有效期
		if time.Until(claims.ExpiresAt.Time) < time.Minute*30 {

			claims.ExpiresAt = &jwt.Time{Time: time.Now().Add(time.Hour * 8)} // 刷新有效期：未来两小时

			newToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
			tokenStr, err := newToken.SignedString([]byte("secret"))
			if err != nil {
				log.Printf("jwt 续约失败，err : %s", err)
			} else {
				ctx.Header("jwt-token", tokenStr)
			}
		}

		// 如果通过验证，则设置 claims 上下文
		ctx.Set("claims", claims)
	}
}
