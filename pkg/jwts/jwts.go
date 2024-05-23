package jwts

import (
	"time"

	"github.com/dgrijalva/jwt-go/v4"
)

type JwtPayload struct {
	UserId    int64  `json:"userId"`
	UserAgent string `json:"userAgent"`
}

type CustomClaims struct {
	JwtPayload
	jwt.StandardClaims // 标准声明结构体
}

// 生成用户 token
func GenToken(user JwtPayload) (string, error) {

	claims := CustomClaims{
		JwtPayload: user,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: jwt.At(time.Now().Add(time.Duration(8) * time.Hour)), // 到期时间
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims) // 加密算法

	return token.SignedString([]byte("uis&*jbb55dHRhf5")) // 使用密钥，生成带签名的JWT
}

// 解析用户 token
func ParseToken(token string) (*CustomClaims, error) {

	Token, err := jwt.ParseWithClaims(token, &CustomClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte("uis&*jbb55dHRhf5"), nil
	})
	if err != nil {
		// 解析 token 异常
		return nil, err
	}
	if !Token.Valid {
		// 令牌无效
		return nil, &jwt.TokenNotValidYetError{}
	}

	claims, ok := Token.Claims.(*CustomClaims)
	if !ok {
		// 数据不一致
		return nil, &jwt.InvalidClaimsError{}
	}

	return claims, nil
}
