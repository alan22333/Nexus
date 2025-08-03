package models

import "gorm.io/gorm"

type User struct {
	gorm.Model

	// --- 核心账户信息 (Core Account Info) ---
	Username string `gorm:"unique;not null;size:50"`
	Email    string `gorm:"unique;not null;size:100"`
	Password string `gorm:"not null"`
	Role     string `gorm:"size:20;default:'user'" `
	Avatar   string `gorm:"size:255"`

	// --- 个人资料 (Profile Details) ---
	Gender int    `gorm:"default:0"` // 0: 未设置, 1: 男, 2: 女
	Phone  string `gorm:"size:20"`
	QQ     string `gorm:"size:20"`
	Wechat string `gorm:"size:50"`
	Bio    string `gorm:"type:text"` // 个人简介

	// --- 隐私设置 (Privacy Settings) ---
	IsPhonePublic  bool `gorm:"default:false"`
	IsEmailPublic  bool `gorm:"default:false"`
	IsQQPublic     bool `gorm:"default:false"`
	IsWechatPublic bool `gorm:"default:false"`
	IsGenderPublic bool `gorm:"default:true"`

	// --- 关联关系 (Associations) ---
	Posts     []*Post    `gorm:"foreignKey:UserID"`              // 用户发表的帖子
	Comments  []*Comment `gorm:"foreignKey:UserID"`              // 用户发表的评论
	Likes     []*Post    `gorm:"many2many:user_post_likes;"`     // 用户点赞的帖子
	Favorites []*Post    `gorm:"many2many:user_post_favorites;"` // 用户收藏的帖子
}
