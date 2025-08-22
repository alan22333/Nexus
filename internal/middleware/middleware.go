package middleware

import (
	"Nuxus/configs"

	"github.com/gin-gonic/gin"
)

// MiddlewareManager 管理所有中间件
type MiddlewareManager struct {
	jwtMiddleware *JWTMiddleware
	config        *configs.Config
}

func NewMiddlewareManager(config *configs.Config) *MiddlewareManager {
	return &MiddlewareManager{
		jwtMiddleware: NewJWTMiddleware(config),
		config:        config,
	}
}

// GetJWTMiddleware 获取JWT中间件
func (mm *MiddlewareManager) GetJWTMiddleware() *JWTMiddleware {
	return mm.jwtMiddleware
}

// JWTAuth 返回JWT认证中间件
func (mm *MiddlewareManager) JWTAuth() gin.HandlerFunc {
	return mm.jwtMiddleware.JWTAuth()
}

// GenerateToken 生成JWT token
func (mm *MiddlewareManager) GenerateToken(userID uint) (string, error) {

	return mm.jwtMiddleware.GenerateToken(userID)
}

// ErrorHandler 错误处理中间件（保持不变）
func (mm *MiddlewareManager) ErrorHandler() gin.HandlerFunc {
	return ErrorHandler()
}

// CORSMiddleware CORS中间件（保持不变）
func (mm *MiddlewareManager) CORSMiddleware() gin.HandlerFunc {
	return CORSMiddleware()
}

// Recovery 恢复中间件（保持不变）
func (mm *MiddlewareManager) Recovery() gin.HandlerFunc {
	return Recovery()
}

// Logger 日志中间件（保持不变）
func (mm *MiddlewareManager) Logger() gin.HandlerFunc {
	return Logger()
}
