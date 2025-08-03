package service

import (
	"Nuxus/internal/dao"
	"Nuxus/internal/dto"
	"Nuxus/internal/models"
	"Nuxus/pkg/erru"
	"log"
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

	return post, nil
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
