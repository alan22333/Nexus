package middleware

import (
	"Nuxus/internal/res"
	"Nuxus/pkg/erru"
	"errors"
	"log"

	"github.com/gin-gonic/gin"
)

// ErrorHandler 是一个中央错误处理中间件
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next() // 先执行后续的 handler

		// 检查 context 中是否有错误
		// c.Errors 是一个 []*gin.Error 类型的切片
		if len(c.Errors) > 0 {
			// 我们只关心第一个错误
			err := c.Errors.Last().Err

			var appErr *erru.AppError
			// 使用 errors.As 来判断错误的具体类型是否是我们的 *AppError
			if errors.As(err, &appErr) {
				// 如果是 AppError，我们知道如何处理它
				// 记录包含完整上下文的错误日志
				log.Printf("Application error: %v\n", appErr)
				// 使用 response 包返回格式化的 JSON
				res.Fail(c, appErr.Code, nil, appErr.Msg)
				return
			}

			// 如果不是我们定义的 AppError，说明是未知的内部错误
			// 记录详细错误
			log.Printf("Internal server error: %v\n", err)
			// 向用户返回一个通用的服务器内部错误，隐藏实现细节
			res.Fail(c, erru.InternalServerError, nil, erru.ErrInternalServer.Msg)
		}
	}
}
