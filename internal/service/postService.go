package service

import (
	"Nuxus/internal/dao"
	"Nuxus/internal/dto"
	"Nuxus/internal/models"
	"Nuxus/pkg/erru"
	"log"

	"gorm.io/gorm"
)

type PostService struct {
	postDAO     *dao.PostDAO
	tagDAO      *dao.TagDAO
	repository  *dao.Repository
	redisClient *dao.RedisClient
}

func NewPostService(postDAO *dao.PostDAO, tagDAO *dao.TagDAO, repository *dao.Repository, redisClient *dao.RedisClient) *PostService {
	return &PostService{
		postDAO:     postDAO,
		tagDAO:      tagDAO,
		repository:  repository,
		redisClient: redisClient,
	}
}

func (p *PostService) ListPosts(reqDto *dto.ListPostsReqDTO) ([]*models.Post, int64, error) {
	posts, total, err := p.postDAO.ListPosts(reqDto)
	if err != nil {
		return nil, 0, erru.ErrInternalServer.Wrap(err)
	}
	return posts, total, nil
}

func (p *PostService) GetPostById(id uint) (*models.Post, error) {
	post, err := p.postDAO.GetPostById(id)
	if err != nil {
		return nil, erru.ErrInternalServer.Wrap(err)
	}

	// 异步更新
	go func() {
		// 增加浏览量
		err := p.redisClient.IncrementPostViewCount(id)
		if err != nil {
			// TODO:异步错误处理
			log.Fatal(erru.New("增加浏览量失败"))
		}
		// 增加热门积分
		// 阅读 +10
		// 点赞 +30
		// 评论 +20
		// 收藏 +30
		err = p.redisClient.IncrementPostRank(id, 10)
		if err != nil {
			log.Fatal(erru.New("增加热门积分失败"))
		}
	}()
	return post, nil
}

func (p *PostService) ListPopularPosts(limit int) ([]*models.Post, error) {
	// get popular postIds from redis
	ids, err := p.redisClient.GetPopularPostIDs(int64(limit))
	if err != nil {
		return nil, erru.ErrInternalServer.Wrap(err)
	}

	if len(ids) == 0 {
		return nil, erru.New("没有热门文章")
	}

	// get post from mysql
	posts, err := p.postDAO.GetPostsByIds(ids)
	if err != nil {
		return nil, erru.ErrInternalServer.Wrap(err)
	}

	return posts, err
}

func (p *PostService) CreatePost(userID uint, reqDto *dto.CreatePostReqDTO) (*models.Post, error) {
	// find or create Tags
	tags, err := p.tagDAO.FindOrCreateTagsByNames(reqDto.Tags)
	if err != nil {
		return nil, erru.ErrInternalServer.Wrap(err)
	}

	post := &models.Post{}

	post.Title = reqDto.Title
	post.Content = reqDto.Content
	post.UserID = userID
	post.Tags = tags

	err = p.postDAO.CreatePost(post)
	if err != nil {
		return nil, erru.ErrInternalServer.Wrap(err)
	}

	fullPost, err := p.postDAO.GetPostById(post.ID)
	if err != nil {
		return nil, erru.ErrInternalServer.Wrap(err)
	}

	return fullPost, nil
}

func (p *PostService) UpdatePost(userId uint, postId uint, reqDto dto.UpdatePostReqDTO) (*models.Post, error) {
	// 更新逻辑
	// 1.检查post是否存在
	// 2.检查是不是当前user的post
	// 3.更新post

	post, err := p.postDAO.GetPostById(postId)
	if err != nil {
		return nil, erru.ErrInternalServer.Wrap(err)
	}
	if post.UserID != userId {
		return nil, erru.ErrUnauthorized
	}

	tags, err := p.tagDAO.FindOrCreateTagsByNames(reqDto.Tags)
	if err != nil {
		return nil, erru.ErrInternalServer.Wrap(err)
	}

	post.Title = reqDto.Title
	post.Content = reqDto.Content
	post.Tags = tags

	post, err = p.postDAO.UpdatePost(post)
	if err != nil {
		return nil, erru.ErrInternalServer.Wrap(err)
	}
	return post, nil
}

func (p *PostService) DeletePost(postId uint, userId uint) error {
	post, err := p.postDAO.GetPostById(postId)
	if err != nil {
		return erru.ErrInternalServer.Wrap(err)
	}
	if post.UserID != userId {
		return erru.ErrUnauthorized
	}

	return p.postDAO.DeletePost(postId)
}

// -------------------评论相关------------------------------
func (p *PostService) ListComment(postId uint, page int, size int) ([]*models.Comment, int64, error) {
	// 评论列表逻辑
	// 1.检查帖子是否存在
	// 2.dao进行分页查询

	_, err := p.postDAO.GetPostById(postId)
	if err != nil {
		return nil, 0, erru.ErrInternalServer.Wrap(err)
	}

	comments, total, err := p.postDAO.ListComment(postId, page, size)
	if err != nil {
		return nil, 0, erru.ErrInternalServer.Wrap(err)
	}

	return comments, total, nil
}

func (p *PostService) CreateComment(req *dto.CreateCommentReqDTO, userId uint, postId uint) (*models.Comment, error) {

	_, err := p.postDAO.GetPostById(postId)
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
	p.repository.DB().Transaction(func(tx *gorm.DB) error {
		// 1. 在事务中创建评论
		if err := p.postDAO.CreateComment(tx, comment); err != nil {
			return err
		}
		// 2. 在事务中更新帖子的评论数
		if err := p.postDAO.UpdatePostCounter(tx, postId, "comment_count", 1); err != nil {
			return err
		}
		return nil // 返回 nil，事务就会被提交
	})

	fullComment, err := p.postDAO.GetCommentById(comment.ID)
	if err != nil {
		return nil, erru.ErrInternalServer.Wrap(err)
	}

	return fullComment, nil
}

func (p *PostService) DeleteComment(commentId uint, userId uint) error {
	comment, err := p.postDAO.GetCommentById(commentId)
	if err != nil {
		return erru.ErrInternalServer.Wrap(err)
	}

	if userId != comment.UserID {
		return erru.ErrUnauthorized
	}

	// TODO:递归删除 or 逻辑删除
	// 我选 后者
	comment.Content = "该评论已删除"
	err = p.postDAO.UpdateComment(comment)
	if err != nil {
		return erru.ErrInternalServer.Wrap(err)
	}
	return err
}

// --------------------点赞、收藏------------------------------
func (p *PostService) LikePost(postId uint, userId uint) (bool, int64, error) {
	// 验证
	post, err := p.postDAO.GetPostById(postId)
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
	isLiked, err := p.postDAO.IsLiked(userId, postId)
	if err != nil {
		return false, 0, erru.ErrInternalServer.Wrap(err)
	}

	var actionState bool
	var newLikeCount int

	err = p.repository.DB().Transaction(func(tx *gorm.DB) error {

		if isLiked == true {
			err := p.postDAO.RemoveLike(tx, userId, postId)
			if err != nil {
				return err
			}
			p.postDAO.UpdatePostCounter(tx, postId, "like_count", -1)
			actionState = false
		} else {
			err := p.postDAO.AddLike(tx, userId, postId)
			if err != nil {
				return err
			}
			p.postDAO.UpdatePostCounter(tx, postId, "like_count", 1)
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

func (p *PostService) FavoritePost(postId uint, userId uint) (bool, int64, error) {
	// 验证
	post, err := p.postDAO.GetPostById(postId)
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
	isFavorite, err := p.postDAO.IsFavorite(userId, postId)
	if err != nil {
		return false, 0, erru.ErrInternalServer.Wrap(err)
	}

	var actionState bool
	var newFavoriteCount int

	err = p.repository.DB().Transaction(func(tx *gorm.DB) error {
		if isFavorite == true {
			err := p.postDAO.RemoveFavorite(tx, userId, postId)
			if err != nil {
				return err
			}
			p.postDAO.UpdatePostCounter(tx, postId, "favorite_count", -1)
			actionState = false
		} else {
			err := p.postDAO.AddFavorite(tx, userId, postId)
			if err != nil {
				return err
			}
			p.postDAO.UpdatePostCounter(tx, postId, "favorite_count", 1)
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

func (p *PostService) GetUserStatus(userId, postId uint) (bool, bool, error) {
	liked, err := p.postDAO.IsFavorite(userId, postId)
	if err != nil {
		return false, false, erru.ErrInternalServer.Wrap(err)
	}

	favorited, err := p.postDAO.IsLiked(userId, postId)
	if err != nil {
		return false, false, erru.ErrInternalServer.Wrap(err)
	}

	return liked, favorited, nil
}
