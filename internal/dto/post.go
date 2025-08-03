package dto

import (
	"time"
)

type ListPostsReqDTO struct {
	Tag  string `form:"tag"`
	Page int    `form:"page,default=1"`
	Size int    `form:"size,default=10"`
}

type ListPostsResDTO struct {
	Total int64            `json:"total"`
	Post  []PostInfoResDTO `json:posts`
}

type PostInfoResDTO struct {
	ID            uint         `json:id`
	Title         string       `json:"title"`
	Author        UserInfoDTO  `json:"author"` // 关联作者信息
	Tags          []TagInfoDTO `json:"tags"`
	ViewCount     int          `json:"view_count"`
	LikeCount     int          `json:"like_count"`
	CommentCount  int          `json:"comment_count"`
	FavoriteCount int          `json:"favorite_count"`
	CreatedAt     time.Time    `json:"created_at"`
}

type PostDetailResDTO struct {
	ID            uint         `json:id`
	Title         string       `json:"title"`
	Content       string       `json:"content"`
	Author        UserInfoDTO  `json:"author"` // 关联作者信息
	Tags          []TagInfoDTO `json:"tags"`
	ViewCount     int          `json:"view_count"`
	LikeCount     int          `json:"like_count"`
	CommentCount  int          `json:"comment_count"`
	FavoriteCount int          `json:"favorite_count"`
	CreatedAt     time.Time    `json:"created_at"`
	UpdatedAt     time.Time    `json:"updated_at"`
}

type CreatePostReqDTO struct {
	Title   string   `json:"title" binding:"required,min=3"`
	Content string   `json:"content" binding:"required,min=5"`
	Tags    []string `json:"tags"`
}

type UpdatePostReqDTO struct {
	Title   string   `json:"title" binding:"required,min=3"`
	Content string   `json:"content" binding:"required,min=5"`
	Tags    []string `json:"tags"`
}
