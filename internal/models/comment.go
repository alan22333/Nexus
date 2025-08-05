package models

import "gorm.io/gorm"

type Comment struct {
	gorm.Model
	Content string `gorm:"type:text;not null"`

	// --- 关联外键 ---
	UserID uint `gorm:"not null"`
	User   User `gorm:"foreignKey:UserID"`

	PostID uint `gorm:"not null"`

	// --- 回复机制 (Reply Mechanism) ---
	// ParentID 指向它所回复的另一条评论的 ID。
	// 如果是顶级评论，ParentID 为 0。
	ParentID uint `gorm:"default:0"`
}
