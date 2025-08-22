package middleware

import (
	"Nuxus/configs"
	"Nuxus/internal/res"
	"Nuxus/pkg/erru"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

type JWTMiddleware struct {
	config *configs.Config
}

func NewJWTMiddleware(config *configs.Config) *JWTMiddleware {
	return &JWTMiddleware{
		config: config,
	}
}

// var duration = configs.Conf.JWT.ExpireHours

type MyClaims struct {
	UserID uint `json:"user_id"`
	jwt.RegisteredClaims
}

// 中间件方法
// JWTAuth 中间件方法 - 完全使用注入的配置
func (jm *JWTMiddleware) JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader == "" {
			res.FailWithAppErr(c, erru.ErrInvalidRequestHeader)
			c.Abort()
			return
		}

		// 按空格分割，格式应为 "Bearer <token>"
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			res.FailWithAppErr(c, erru.ErrInvalidRequestHeader)
			c.Abort()
			return
		}

		// 提取 token 字符串部分
		tokenString := parts[1]

		// 使用注入的配置中的密钥
		jwtSecret := []byte(jm.config.JWT.Secret)

		// 使用 ParseWithClaims 解析 JWT
		token, err := jwt.ParseWithClaims(tokenString, &MyClaims{}, func(token *jwt.Token) (any, error) {
			return jwtSecret, nil
		})
		if err != nil {
			res.FailWithAppErr(c, erru.ErrTokenInvalid.Wrap(err))
			c.Abort()
			return
		}

		// 校验通过后，取出 claims 中的数据
		if claims, ok := token.Claims.(*MyClaims); ok && token.Valid {
			c.Set("userID", claims.UserID)
			c.Next()
		} else {
			res.FailWithAppErr(c, erru.ErrTokenInvalid)
			c.Abort()
			return
		}
	}
}

// GenerateToken 生成JWT token
func (jm *JWTMiddleware) GenerateToken(userID uint) (string, error) {
	expireHours := time.Duration(jm.config.JWT.ExpireHours)
	jwtSecret := []byte(jm.config.JWT.Secret)

	claims := MyClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expireHours * time.Hour)),
			Issuer:    "Nexus",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}
