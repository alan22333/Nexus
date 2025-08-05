package dto

import "time"

type PrivacyInfo struct {
	IsPhonePublic  bool `json:"is_phone_public" binding:"boolean"` // 确保是布尔值
	IsEmailPublic  bool `json:"is_email_public" binding:"boolean"`
	IsQQPublic     bool `json:"is_qq_public" binding:"boolean"`
	IsWechatPublic bool `json:"is_wechat_public" binding:"boolean"`
	IsGenderPublic bool `json:"is_gender_public" binding:"boolean"`
}
type ProfileResDTO struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Role     string `json:"role"`
	Avatar   string `json:"avatar"`

	// --- 个人资料 (Profile Details) ---
	Gender int    `json:"gender"` // 0: 未设置, 1: 男, 2: 女
	Phone  string `json:"phone"`
	QQ     string `json:"qq"`
	Wechat string `json:"wechat"`
	Bio    string `json:"bio"` // 个人简介

	Privacy PrivacyInfo `json:"privacy"`

	CreatedAt time.Time `json:"created_at"`
}

type ProfileReqDTO struct {
	// oneof=0 1 2: 确保 gender 的值只能是 0, 1, 或 2 之一
	Gender int    `json:"gender" binding:"required,oneof=0 1 2"`
	Phone  string `json:"phone" binding:"omitempty,len=11"`
	QQ     string `json:"qq" binding:"omitempty,numeric,min=5,max=11"`
	Wechat string `json:"wechat" binding:"omitempty,min=6,max=20"`
	Bio    string `json:"bio" binding:"omitempty,max=200"`

	Privacy PrivacyInfo `json:"privacy" binding:"required"`
}
