package controller

import (
	"Nuxus/internal/dto"
	"Nuxus/internal/models"
	"Nuxus/internal/res"
	"Nuxus/internal/service"
	"Nuxus/pkg/erru"
	"log"
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
	// postId := c.Param("id")
	postId, _ := strconv.ParseUint(c.Param("id"), 10, 32)

	if postId == 0 {
		c.Error(erru.ErrInvalidParams)
		return
	}
	// log.Println("postId:", postId)
	post, err := service.GetPostById(uint(postId))
	if err != nil {
		c.Error(err)
		return
	}

	log.Println(post.User)

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
		ID:            post.ID,
		Title:         post.Title,
		Author:        *userModel2InfoDto(&post.User),
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
	// postId := c.Param("id")
	postId, _ := strconv.ParseUint(c.Param("id"), 10, 32)

	err := c.ShouldBindJSON(&reqDto)
	if postId == 0 || err != nil {
		c.Error(erru.ErrInvalidParams)
		return
	}
	userId := c.MustGet("userID").(uint)

	post, err := service.UpdatePost(userId, uint(postId), reqDto)
	if err != nil {
		c.Error(err)
		return
	}

	resDto := postModel2DetailDTO(post)
	res.OkWithData(c, resDto)
}

func DeletePost(c *gin.Context) {
	postId, _ := strconv.ParseUint(c.Param("id"), 10, 32)

	if postId == 0 {
		c.Error(erru.ErrInvalidParams)
		return
	}
	userId := c.MustGet("userID").(uint)

	err := service.DeletePost(uint(postId), userId)
	if err != nil {
		c.Error(err)
		return
	}

	res.OkWithMsg(c, "删除成功")
}

// --------------交互相关：点赞、收藏、评论------------------------------

func ListComment(c *gin.Context) {
	postID, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))

	comments, total, err := service.ListComment(uint(postID), page, size)
	if err != nil {
		c.Error(err)
		return
	}

	commentsDto := make([]dto.CommentInfo, 0, len(comments))
	for _, comment := range comments {
		commentsDto = append(commentsDto, *commentModel2ResDTO(comment))
	}

	resDto := dto.ListCommentResDto{
		Total:    total,
		Comments: commentsDto,
	}

	res.OkWithData(c, resDto)
}

func GetUserStatus(c *gin.Context) {
	postId, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	userId := c.MustGet("userID").(uint)

	liked, favorited, err := service.GetUserStatus(userId, uint(postId))
	if err != nil {
		c.Error(err)
	}

	resDto := dto.GetUserStatusResDTO{
		Liked:     liked,
		Favorited: favorited,
	}

	res.OkWithData(c, resDto)
}

func commentModel2ResDTO(comment *models.Comment) *dto.CommentInfo {
	return &dto.CommentInfo{
		Id:        comment.ID,
		Content:   comment.Content,
		Author:    *userModel2InfoDto(&comment.User),
		ParentId:  comment.ParentID,
		CreatedAt: comment.CreatedAt,
	}
}

func CreateComment(c *gin.Context) {
	var reqDto dto.CreateCommentReqDTO
	err := c.ShouldBindJSON(&reqDto)
	if err != nil {
		c.Error(erru.ErrInvalidParams.Wrap(err))
		return
	}

	postId, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	userId := c.MustGet("userID").(uint)

	comment, err := service.CreateComment(&reqDto, userId, uint(postId))
	if err != nil {
		c.Error(err)
		return
	}
	resDto := commentModel2ResDTO(comment)
	res.OkWithData(c, resDto)
}

func DeleteComment(c *gin.Context) {
	commentId, _ := strconv.ParseUint(c.Param("commentId"), 10, 32)
	userId := c.MustGet("userID").(uint)

	err := service.DeleteComment(uint(commentId), userId)
	if err != nil {
		c.Error(err)
		return
	}

	res.OkWithMsg(c, "删除评论成功")
}

// ---------------------点赞、收藏--------------------------------
func LikePost(c *gin.Context) {
	postId, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	userId := c.MustGet("userID").(uint)

	actionState, newLikeCount, err := service.LikePost(uint(postId), userId)
	if err != nil {
		c.Error(err)
		return
	}

	resDto := &dto.ToggleActionResDTO{
		ActionState:  actionState,
		CurrentCount: newLikeCount,
	}

	res.OkWithData(c, resDto)
}

func FavoritePost(c *gin.Context) {
	postId, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	userId := c.MustGet("userID").(uint)

	actionState, newFavoriteCount, err := service.FavoritePost(uint(postId), userId)
	if err != nil {
		c.Error(err)
		return
	}

	resDto := &dto.ToggleActionResDTO{
		ActionState:  actionState,
		CurrentCount: newFavoriteCount,
	}

	res.OkWithData(c, resDto)
}
