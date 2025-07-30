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

// 假设的配置，你应该从你的 config 包导入
var jwtSecret = []byte(configs.Conf.JWT.Secret)
var duration = configs.Conf.JWT.ExpireHours

type MyClaims struct {
	UserID uint `json:"user_id"`
	jwt.RegisteredClaims
}

// 中间件方法
func JWTAuth() gin.HandlerFunc {
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
			// 格式不合法，如没有 Bearer 前缀，返回失败
			res.FailWithAppErr(c, erru.ErrInvalidRequestHeader)
			c.Abort()
			return
		}

		// 提取 token 字符串部分
		tokenString := parts[1]

		// 使用 ParseWithClaims 解析 JWT，并指定我们自定义的 claims 结构体
		token, err := jwt.ParseWithClaims(tokenString, &MyClaims{}, func(token *jwt.Token) (any, error) {
			// 返回密钥，用于验证签名
			return jwtSecret, nil
		})
		if err != nil {
			// token 解析失败（可能是签名错误、过期等）
			res.FailWithAppErr(c, erru.ErrTokenInvalid.Wrap(err))
			c.Abort()
			return
		}

		// 校验通过后，取出 claims 中的数据
		if claims, ok := token.Claims.(*MyClaims); ok && token.Valid {
			// 将用户 ID 写入 Gin 的 context 中，方便后续处理器使用
			c.Set("userID", claims.UserID)
			c.Next() // 继续执行后续处理器
		} else {
			// Token 无效
			res.FailWithAppErr(c, erru.ErrTokenInvalid)
			c.Abort()
			return
		}
	}
}

// GenerateToken 生成一个 JWT token（传入用户ID）
func GenerateToken(userID uint) (string, error) {
	claims := MyClaims{
		UserID: userID, // 写入自定义字段
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Duration(duration) * time.Hour)), // 设置7天过期时间
			Issuer:    "Nexus",                                                               // 谁签发了这个token（可选）
		},
	}

	// 生成 token，使用 HS256 对称加密算法
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 用密钥签名并转成字符串
	return token.SignedString(jwtSecret)
}
