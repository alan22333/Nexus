package middleware

import (
	"Nuxus/internal/res"
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

// 颜色代码
const (
	green  = "\033[32m"
	yellow = "\033[33m"
	red    = "\033[31m"
	blue   = "\033[34m"
	reset  = "\033[0m"
)

// Logger 是一个美化后的 Gin 日志中间件
// 感觉不如自带的哈哈
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		c.Next()

		latency := time.Since(start)
		clientIP := c.ClientIP()
		method := c.Request.Method
		statusCode := c.Writer.Status()

		// 拼接状态码颜色
		var statusColor string
		switch {
		case statusCode >= 200 && statusCode < 300:
			statusColor = green
		case statusCode >= 300 && statusCode < 400:
			statusColor = blue
		case statusCode >= 400 && statusCode < 500:
			statusColor = yellow
		default:
			statusColor = red
		}

		// 拼接请求路径（带 query）
		fullPath := path
		if raw != "" {
			fullPath = path + "?" + raw
		}

		log.Printf("[GIN] |%s %3d %s| %13v | %-15s | %-7s %s\n",
			statusColor, statusCode, reset, // 彩色状态码
			latency,  // 请求耗时
			clientIP, // 客户端 IP
			method,   // 请求方法
			fullPath, // 完整路径
		)
	}
}

// CORS Middleware
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*") // 允许所有来源，生产环境应配置为前端域名
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
		c.Header("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}

// Recovery Middleware
func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("panic recovered: %v", err)
				// esponse 包返回一个标准的 500 错误
				res.Fail(c, 500, nil, "服务器内部错误")
				c.Abort()
			}
		}()
		c.Next()
	}
}
