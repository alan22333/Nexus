package dto

type TagInfoDTO struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type ListTagsResDTO struct {
	ID        uint   `json:"id"`
	Name      string `json:"name"`
	PostCount int64  `json:"post_count"`
}
