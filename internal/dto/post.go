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

// -------------------评论--------------------------------
type ListCommentResDto struct {
	Total    int64         `json:"total"`
	Comments []CommentInfo `json:"comments"`
}

type CommentInfo struct {
	Id        uint        `json:"id"`
	Content   string      `json:"content"`
	Author    UserInfoDTO `json:"author"`
	ParentId  uint        `json:"parent_id"`
	CreatedAt time.Time   `json:"created_at"`
}

type CreateCommentReqDTO struct {
	Content  string `json:"content" binding:"required,min=1,max=500"`
	ParentId uint   `json:"parent_id"`
}

// ---------------------点赞、收藏---------------------------------
// ToggleActionResDTO 用于点赞/收藏操作的统一响应
type ToggleActionResDTO struct {
	ActionState  bool  `json:"action_state"`
	CurrentCount int64 `json:"current_count"`
}

type GetUserStatusResDTO struct {
	Liked     bool `json:"liked"`
	Favorited bool `json:"favorited"`
}
