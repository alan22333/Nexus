package service

import (
	"Nuxus/internal/dao"
	"Nuxus/internal/dto"
	"Nuxus/internal/models"
	"Nuxus/pkg/erru"
	"fmt"
	"log"

	"gorm.io/gorm"
)

func ListPosts(reqDto *dto.ListPostsReqDTO) ([]*models.Post, int64, error) {
	posts, total, err := dao.ListPosts(reqDto)
	if err != nil {
		return nil, 0, erru.ErrInternalServer.Wrap(err)
	}
	return posts, total, nil
}

func GetPostById(id string) (*models.Post, error) {
	post, err := dao.GetPostById(id)
	if err != nil {
		return nil, erru.ErrInternalServer.Wrap(err)
	}

	// 异步更新
	go func() {
		// 增加浏览量
		err := dao.IncrementPostViewCount(id)
		if err != nil {
			// TODO:异步错误处理
			log.Fatal(erru.New("增加浏览量失败"))
		}
		// 增加热门积分
		// 阅读 +10
		// 点赞 +30
		// 评论 +20
		// 收藏 +30
		err = dao.IncrementPostRank(id, 10)
		if err != nil {
			log.Fatal(erru.New("增加热门积分失败"))
		}
	}()
	return post, nil
}

func ListPopularPosts(limit int) ([]*models.Post, error) {
	// get popular postIds from redis
	ids, err := dao.GetPopularPostIDs(int64(limit))
	if err != nil {
		return nil, erru.ErrInternalServer.Wrap(err)
	}

	if len(ids) == 0 {
		return nil, erru.New("没有热门文章")
	}

	// get post from mysql
	posts, err := dao.GetPostsByIds(ids)
	if err != nil {
		return nil, erru.ErrInternalServer.Wrap(err)
	}

	return posts, err
}

func CreatePost(userID uint, reqDto *dto.CreatePostReqDTO) (*models.Post, error) {
	// find or create Tags
	tags, err := dao.FindOrCreateTagsByNames(reqDto.Tags)
	if err != nil {
		return nil, erru.ErrInternalServer.Wrap(err)
	}

	post := &models.Post{}

	post.Title = reqDto.Title
	post.Content = reqDto.Content
	post.UserID = userID
	post.Tags = tags

	err = dao.CreatePost(post)
	if err != nil {
		return nil, erru.ErrInternalServer.Wrap(err)
	}

	fullPost, err := dao.GetPostById(fmt.Sprint(post.ID))
	if err != nil {
		return nil, erru.ErrInternalServer.Wrap(err)
	}

	return fullPost, nil
}

func UpdatePost(userId uint, postId string, reqDto dto.UpdatePostReqDTO) (*models.Post, error) {
	// 更新逻辑
	// 1.检查post是否存在
	// 2.检查是不是当前user的post
	// 3.更新post

	post, err := dao.GetPostById(postId)
	if err != nil {
		return nil, erru.ErrInternalServer.Wrap(err)
	}
	if post.UserID != userId {
		return nil, erru.ErrUnauthorized
	}

	tags, err := dao.FindOrCreateTagsByNames(reqDto.Tags)
	if err != nil {
		return nil, erru.ErrInternalServer.Wrap(err)
	}

	post.Title = reqDto.Title
	post.Content = reqDto.Content
	post.Tags = tags

	post, err = dao.UpdatePost(post)
	if err != nil {
		return nil, erru.ErrInternalServer.Wrap(err)
	}
	return post, nil
}

func DeletePost(postId string, userId uint) error {
	post, err := dao.GetPostById(postId)
	if err != nil {
		return erru.ErrInternalServer.Wrap(err)
	}
	if post.UserID != userId {
		return erru.ErrUnauthorized
	}

	return dao.DeletePost(postId)
}

// -------------------评论相关------------------------------
func ListComment(postId uint, page int, size int) ([]*models.Comment, int64, error) {
	// 评论列表逻辑
	// 1.检查帖子是否存在
	// 2.dao进行分页查询

	_, err := dao.GetPostById(fmt.Sprint(postId))
	if err != nil {
		return nil, 0, erru.ErrInternalServer.Wrap(err)
	}

	comments, total, err := dao.ListComment(postId, page, size)
	if err != nil {
		return nil, 0, erru.ErrInternalServer.Wrap(err)
	}

	return comments, total, nil
}

func CreateComment(req *dto.CreateCommentReqDTO, userId uint, postId uint) (*models.Comment, error) {

	_, err := dao.GetPostById(fmt.Sprint(postId))
	if err != nil {
		return nil, erru.ErrInternalServer.Wrap(err)
	}

	comment := &models.Comment{
		Content:  req.Content,
		UserID:   userId,
		ParentID: req.ParentId,
		PostID:   postId,
	}

	// TODO:（数据库事务）
	// 确保“创建评论”和“帖子评论数+1”这两个操作，要么都成功，要么都失败
	dao.DB.Transaction(func(tx *gorm.DB) error {
		// 1. 在事务中创建评论
		if err := dao.CreateComment(tx, comment); err != nil {
			return err
		}
		// 2. 在事务中更新帖子的评论数
		if err := dao.UpdatePostCounter(tx, postId, "comment_count", 1); err != nil {
			return err
		}
		return nil // 返回 nil，事务就会被提交
	})

	fullComment, err := dao.GetCommentById(comment.ID)
	if err != nil {
		return nil, erru.ErrInternalServer.Wrap(err)
	}

	return fullComment, nil
}

func DeleteComment(commentId uint, userId uint) error {
	comment, err := dao.GetCommentById(commentId)
	if err != nil {
		return erru.ErrInternalServer.Wrap(err)
	}

	if userId != comment.UserID {
		return erru.ErrUnauthorized
	}

	// TODO:递归删除 or 逻辑删除
	// 我选 后者
	comment.Content = "该评论已删除"
	err = dao.UpdateComment(comment)
	if err != nil {
		return erru.ErrInternalServer.Wrap(err)
	}
	return err
}

// --------------------点赞、收藏------------------------------
func LikePost(postId uint, userId uint) (bool, int64, error) {
	// 验证
	post, err := dao.GetPostById(fmt.Sprint(postId))
	if err != nil {
		return false, 0, erru.ErrInternalServer.Wrap(err)
	}
	// if post.UserID != userId {
	// 	return false, 0, erru.ErrUnauthorized
	// }
	if userId == 0 {
		return false, 0, erru.ErrUnauthorized
	}
	// 查询是否点赞
	isLiked, err := dao.IsLiked(userId, postId)
	if err != nil {
		return false, 0, erru.ErrInternalServer.Wrap(err)
	}

	var actionState bool
	var newLikeCount int

	err = dao.DB.Transaction(func(tx *gorm.DB) error {

		if isLiked == true {
			err := dao.RemoveLike(tx, userId, postId)
			if err != nil {
				return err
			}
			dao.UpdatePostCounter(tx, postId, "like_count", -1)
			actionState = false
		} else {
			err := dao.AddLike(tx, userId, postId)
			if err != nil {
				return err
			}
			dao.UpdatePostCounter(tx, postId, "like_count", 1)
			actionState = true
		}
		return nil
	})

	if err != nil {
		return false, 0, erru.ErrInternalServer.Wrap(err)
	}

	// 点赞更新
	newLikeCount = post.LikeCount
	if actionState == false {
		newLikeCount--
	} else {
		newLikeCount++
	}

	return actionState, int64(newLikeCount), nil
}

func FavoritePost(postId uint, userId uint) (bool, int64, error) {
	// 验证
	post, err := dao.GetPostById(fmt.Sprint(postId))
	if err != nil {
		return false, 0, erru.ErrInternalServer.Wrap(err)
	}
	// if post.UserID != userId {
	// 	return false, 0, erru.ErrUnauthorized
	// }
	if userId == 0 {
		return false, 0, erru.ErrUnauthorized
	}
	// 查询是否收藏
	isFavorite, err := dao.IsFavorite(userId, postId)
	if err != nil {
		return false, 0, erru.ErrInternalServer.Wrap(err)
	}

	var actionState bool
	var newFavoriteCount int

	err = dao.DB.Transaction(func(tx *gorm.DB) error {
		if isFavorite == true {
			err := dao.RemoveFavorite(tx, userId, postId)
			if err != nil {
				return err
			}
			dao.UpdatePostCounter(tx, postId, "favorite_count", -1)
			actionState = false
		} else {
			err := dao.AddFavorite(tx, userId, postId)
			if err != nil {
				return err
			}
			dao.UpdatePostCounter(tx, postId, "favorite_count", 1)
			actionState = true
		}
		return nil
	})

	if err != nil {
		return false, 0, erru.ErrInternalServer.Wrap(err)
	}

	// 点赞更新
	newFavoriteCount = post.FavoriteCount
	if actionState == false {
		newFavoriteCount--
	} else {
		newFavoriteCount++
	}

	return actionState, int64(newFavoriteCount), nil
}

func GetUserStatus(userId, postId uint) (bool, bool, error) {
	liked, err := dao.IsFavorite(userId, postId)
	if err != nil {
		return false, false, erru.ErrInternalServer.Wrap(err)
	}

	favorited, err := dao.IsLiked(userId, postId)
	if err != nil {
		return false, false, erru.ErrInternalServer.Wrap(err)
	}

	return liked, favorited, nil
}
