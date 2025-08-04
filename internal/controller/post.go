package controller

import (
	"Nuxus/internal/dto"
	"Nuxus/internal/models"
	"Nuxus/internal/res"
	"Nuxus/internal/service"
	"Nuxus/pkg/erru"
	"strconv"

	"github.com/gin-gonic/gin"
)

func ListPosts(c *gin.Context) {
	var reqDto dto.ListPostsReqDTO
	err := c.ShouldBindQuery(&reqDto)
	if err != nil {
		c.Error(erru.ErrInvalidParams.Wrap(err))
		return
	}

	if reqDto.Page == 0 {
		reqDto.Page = 1
	}
	if reqDto.Size == 0 {
		reqDto.Size = 10
	}

	posts, total, err := service.ListPosts(&reqDto)
	if err != nil {
		c.Error(err)
		return
	}

	postInfos := make([]dto.PostInfoResDTO, 0, len(posts))
	for _, post := range posts {
		postInfos = append(postInfos, *postModel2InfoDTO(post))
	}

	listPostsResDTO := dto.ListPostsResDTO{
		Total: total,
		Post:  postInfos,
	}

	res.OkWithData(c, listPostsResDTO)
}

func ListPopularPosts(c *gin.Context) {
	limit := c.Param("limit")
	limitNum, _ := strconv.Atoi(limit)
	if limitNum == 0 {
		limitNum = 10
	}

	posts, err := service.ListPopularPosts(limitNum)
	if err != nil {
		c.Error(err)
		return
	}

	postInfos := make([]dto.PostInfoResDTO, 0, len(posts))
	for _, post := range posts {
		postInfos = append(postInfos, *postModel2InfoDTO(post))
	}

	res.OkWithData(c, postInfos)
}

func GetPost(c *gin.Context) {
	// var postId string
	// err := c.ShouldBindUri(postId)
	postId := c.Param("id")
	if postId == "" {
		c.Error(erru.ErrInvalidParams)
		return
	}
	// log.Println("postId:", postId)
	post, err := service.GetPostById(postId)
	if err != nil {
		c.Error(err)
		return
	}

	// encapsulate
	postDetail := postModel2DetailDTO(post)

	res.OkWithData(c, postDetail)
}

func postModel2InfoDTO(post *models.Post) *dto.PostInfoResDTO {
	postInfo := &dto.PostInfoResDTO{
		ID:    post.ID,
		Title: post.Title,
		Author: dto.UserInfoDTO{
			ID:       post.User.ID,
			UserName: post.User.Username,
			Email:    post.User.Email,
			Avatar:   post.User.Avatar,
		},

		ViewCount:     post.ViewCount,
		LikeCount:     post.LikeCount,
		CommentCount:  post.CommentCount,
		FavoriteCount: post.FavoriteCount,
		CreatedAt:     post.CreatedAt,
	}
	tags := make([]dto.TagInfoDTO, 0, len(post.Tags))
	for _, tag := range post.Tags {
		tags = append(tags, *tagModel2InfoDTO(tag))
	}
	postInfo.Tags = tags
	return postInfo
}

func postModel2DetailDTO(post *models.Post) *dto.PostDetailResDTO {
	postInfo := &dto.PostDetailResDTO{
		ID:    post.ID,
		Title: post.Title,
		Author: dto.UserInfoDTO{
			ID:       post.User.ID,
			UserName: post.User.Username,
			Email:    post.User.Email,
			Avatar:   post.User.Avatar,
		},
		Content:       post.Content,
		ViewCount:     post.ViewCount,
		LikeCount:     post.LikeCount,
		CommentCount:  post.CommentCount,
		FavoriteCount: post.FavoriteCount,
		CreatedAt:     post.CreatedAt,
		UpdatedAt:     post.UpdatedAt,
	}
	tags := make([]dto.TagInfoDTO, 0, len(post.Tags))
	for _, tag := range post.Tags {
		tags = append(tags, *tagModel2InfoDTO(tag))
	}
	postInfo.Tags = tags
	return postInfo
}

func tagModel2InfoDTO(tag *models.Tag) *dto.TagInfoDTO {
	return &dto.TagInfoDTO{
		ID:   tag.ID,
		Name: tag.Name,
	}
}

func CreatePost(c *gin.Context) {
	var reqDto dto.CreatePostReqDTO
	err := c.ShouldBindJSON(&reqDto)
	if err != nil {
		c.Error(erru.ErrInvalidParams.Wrap(err))
		return
	}
	// log.Println(reqDto)
	userId := c.MustGet("userID").(uint)

	post, err := service.CreatePost(userId, &reqDto)
	if err != nil {
		c.Error(err)
		return
	}

	postDetailResDto := postModel2DetailDTO(post)

	res.Ok(c, postDetailResDto, "创建成功")
}

func UpdatePost(c *gin.Context) {
	var reqDto dto.UpdatePostReqDTO
	postId := c.Param("id")
	err := c.ShouldBindJSON(&reqDto)
	if postId == "" || err != nil {
		c.Error(erru.ErrInvalidParams)
		return
	}
	userId := c.MustGet("userID").(uint)

	post, err := service.UpdatePost(userId, postId, reqDto)
	if err != nil {
		c.Error(err)
		return
	}

	resDto := postModel2DetailDTO(post)
	res.OkWithData(c, resDto)
}

func DeletePost(c *gin.Context) {
	postId := c.Param("id")
	if postId == "" {
		c.Error(erru.ErrInvalidParams)
		return
	}
	userId := c.MustGet("userID").(uint)

	err := service.DeletePost(postId, userId)
	if err != nil {
		c.Error(err)
		return
	}

	res.OkWithMsg(c, "删除成功")
}
