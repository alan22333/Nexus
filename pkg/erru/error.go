package erru

import "fmt"

// AppError 是我们项目中所有错误的标准格式
type AppError struct {
	Code int    // 业务错误码
	Msg  string // 面向用户的错误信息
	Err  error  // 包装的原始错误，用于日志记录
}

// Error 实现 Go 内置的 error 接口
func (e *AppError) Error() string {
	if e.Err != nil {
		// 包含原始错误信息，方便日志记录
		return fmt.Sprintf("AppError: Code=%d, Msg=%s, Err=%v", e.Code, e.Msg, e.Err)
	}
	return fmt.Sprintf("AppError: Code=%d, Msg=%s", e.Code, e.Msg)
}

// Unwrap 提供对原始错误的支持，以便使用 errors.Is/As
func (e *AppError) Unwrap() error {
	return e.Err
}

// Wrap 用于包装一个已有的 error，并附加错误码和消息
func (e *AppError) Wrap(err error) *AppError {
	// 返回一个新的 AppError 实例，以避免修改原始的错误定义
	return &AppError{
		Code: e.Code,
		Msg:  e.Msg,
		Err:  err,
	}
}

// New 创建一个新的 AppError 实例，它使用通用的业务错误码，但允许自定义错误消息。
// 这对于那些不需要预定义、消息内容不固定的业务错误非常有用。
func New(msg string) *AppError {
	return &AppError{
		Code: BusinessLogicError, // 使用我们新定义的通用错误码
		Msg:  msg,
		Err:  nil, // 初始时没有包装的底层错误
	}
}
