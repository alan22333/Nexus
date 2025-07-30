package erru

// 定义业务错误码
const (
	// ================== 通用错误 ==================
	OK                  = 0
	InternalServerError = 10001
	InvalidParams       = 10002

	// ================== 用户相关错误 =================
	UserNotFound      = 20001
	PasswordIncorrect = 20002
	EmailAlreadyUsed  = 20003
	InvalidVerifyCode = 20004

	// ================== 认证授权相关 =================
	TokenNotFound = 30001
	TokenInvalid  = 30002
	TokenExpired  = 30003
	Unauthorized  = 30004 // 已认证，但无权访问资源

	// ================== 资源相关错误 =================
	ResourceNotFound     = 40001
	InvalidRequestHeader = 40002

	BusinessLogicError = 50001
)

// 预先定义好常用的错误，可以直接在代码中使用
var (
	ErrOK             = &AppError{Code: OK, Msg: "成功"}
	ErrInternalServer = &AppError{Code: InternalServerError, Msg: "服务器内部错误"}
	ErrInvalidParams  = &AppError{Code: InvalidParams, Msg: "参数无效"}

	ErrUserNotFound      = &AppError{Code: UserNotFound, Msg: "用户不存在"}
	ErrPasswordIncorrect = &AppError{Code: PasswordIncorrect, Msg: "密码错误"}
	ErrEmailAlreadyUsed  = &AppError{Code: EmailAlreadyUsed, Msg: "邮箱已被注册"}
	ErrInvaliVerifyCode  = &AppError{Code: InvalidVerifyCode, Msg: "验证码错误"}

	ErrTokenNotFound = &AppError{Code: TokenNotFound, Msg: "未找到认证Token"}
	ErrTokenInvalid  = &AppError{Code: TokenInvalid, Msg: "认证Token无效"}
	ErrTokenExpired  = &AppError{Code: TokenExpired, Msg: "认证Token已过期"}
	ErrUnauthorized  = &AppError{Code: Unauthorized, Msg: "无权执行此操作"}

	ErrResourceNotFound     = &AppError{Code: ResourceNotFound, Msg: "资源未找到"}
	ErrInvalidRequestHeader = &AppError{Code: InvalidRequestHeader, Msg: "请求头错误"}
)
