package dao

import (
	"Nuxus/internal/dto"
	"Nuxus/internal/models"
	"errors"
)

func ListPosts(reqDto *dto.ListPostsReqDTO) ([]*models.Post, int64, error) {
	var posts []*models.Post
	var total int64

	// 1. 构建基础查询
	// Preload("Tags") 是一个 GORM 的强大功能，它会高效地执行另一条查询，
	query := DB.Model(&models.Post{}).Preload("Tags")

	// 2. 如果提供了 tag，则添加过滤条件
	if reqDto.Tag != "" {
		// 这是 GORM 中进行多对多查询的关键！
		// 我们需要 JOIN 中间表和目标表，然后在外键上进行筛选。
		query = query.Joins("JOIN post_tags ON posts.id = post_tags.post_id").
			Joins("JOIN tags ON tags.id = post_tags.tag_id").
			Where("tags.name = ?", reqDto.Tag)
	}

	// 3. 首先，在不应用分页的情况下，计算总数
	// 我们在应用 `LIMIT` 和 `OFFSET` 之前执行 `Count`，这样就能得到满足条件的总记录数。
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 4. 然后，添加分页和排序条件
	// Offset 计算：(页码 - 1) * 每页数量
	offset := (reqDto.Page - 1) * reqDto.Size
	// 按创建时间降序排序，最新的在前面
	query = query.Offset(offset).Limit(reqDto.Size).Order("created_at DESC")

	// 5. 执行最终查询，获取当前页的数据
	if err := query.Find(&posts).Error; err != nil {
		return nil, 0, err
	}

	return posts, total, nil
}

func GetPostById(id string) (*models.Post, error) {
	var post models.Post
	err := DB.Where("id=?", id).First(&post).Error
	if err != nil {
		return nil, err
	}

	return &post, nil
}

func AddPostViewCount(postId uint, incr int) error {
	var post models.Post
	err := DB.Where("id=?", postId).First(&post).Error
	if err != nil {
		return err
	}
	post.ViewCount += incr
	// gorm的更新设计简直是逆天
	return DB.Model(&post).Where("id=?", post.ID).Updates(post).Error
}

func GetPostsByIds(ids []string) ([]*models.Post, error) {
	var posts []*models.Post
	err := DB.Where("id IN (?)", ids).Find(&posts).Error
	if err != nil {
		return nil, err
	}
	return posts, nil
}

func CreatePost(post *models.Post) error {
	return DB.Create(post).Error
}

func UpdatePost(post *models.Post) (*models.Post, error) {
	res := DB.Model(post).Where("id=?", post.ID).Updates(post)
	if res.Error != nil {
		return nil, res.Error
	}
	if res.RowsAffected == 0 {
		return nil, errors.New("更新失败")
	}
	res.First(post)
	return post, nil
}

func DeletePost(postId string) error {
	return DB.Delete(&models.Post{}, postId).Error
}
