package dto

type RegisterReqDTO struct {
	Email string `json:"email" binding:"required,email"`
	// Password string `json:"password" binding:"required,min=6,max=15"`
}

type VerifyRegisterReqDTO struct {
	Email    string `json:"email" binding:"required,email"`
	UserName string `json:"username" binding:"required,min=1,max=20"`
	Password string `json:"password" binding:"required,min=6,max=15"`
	Code     string `json:code binding:"required,len=6"`
}

type LoginReqDTO struct {
	Identifier string `json:"identifier" binding:"required"`
	// Email    string `json:"email" binding:"required,email"`
	// UserName string `json:"username" binding:"required,min=1,max=20"`
	Password string `json:"password" binding:"required,min=6,max=15"`
}

// RegisterResponseDTO 定义了注册成功后返回的数据结构（不含密码）
type UserInfoDTO struct {
	ID       uint   `json:"id"`
	UserName string `json:"username"`
	Email    string `json:"email"`
}

type LoginResponseDTO struct {
	User  UserInfoDTO `json:"user"`
	Token string      `json:"token"`
}

type RequestResetReqDTO struct {
	Email string `json:"email" binding:"required,email"`
}

type VerifyResetReqDTO struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6,max=15"`
	Code     string `json:code binding:"required,len=6"`

}
