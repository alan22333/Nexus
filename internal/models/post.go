// alan-nexus/internal/models/post.go
package models

import "gorm.io/gorm"

type Post struct {
	gorm.Model

	// --- 核心内容 (Core Content) ---
	Title string `gorm:"not null;size:100"`
	// 存储原始 Markdown 文本，前端直接用此内容进行渲染
	Content string `gorm:"type:text;not null"`

	// --- 关联外键 (Foreign Keys) ---
	UserID uint `gorm:"not null"`
	User   User `gorm:"foreignKey:UserID"` // 关联作者信息

	// --- 互动计数 (Interaction Counters) ---
	// 这些计数字段可以通过触发器或在业务代码中维护，用于列表页快速展示
	// 也可以在查询时动态计算，取决于性能需求
	ViewCount     int `gorm:"default:0"`
	LikeCount     int `gorm:"default:0"`
	FavoriteCount int `gorm:"default:0"`
	CommentCount  int `gorm:"default:0"`

	// --- 关联关系 (Associations) ---
	Comments         []*Comment `gorm:"foreignKey:PostID"` // 帖子的所有评论
	Tags             []*Tag     `gorm:"many2many:post_tags;"`
	LikedByUser      []*User    `gorm:"many2many:user_post_likes;"`     // 用户点赞的帖子
	FavoritedByUsers []*User    `gorm:"many2many:user_post_favorites;"` // 用户收藏的帖子
}
