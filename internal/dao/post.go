package dao

import (
	"Nuxus/internal/dto"
	"Nuxus/internal/models"
	"errors"

	"gorm.io/gorm"
)

func ListPosts(reqDto *dto.ListPostsReqDTO) ([]*models.Post, int64, error) {
	var posts []*models.Post
	var total int64

	// 1. 构建基础查询
	// Preload("Tags") 是一个 GORM 的强大功能，它会高效地执行另一条查询，
	query := DB.Model(&models.Post{}).Preload("Tags").Preload("User")

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
	err := DB.Where("id=?", id).Preload("User").First(&post).Error
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

// -----------------评论----------------------------
func ListComment(postID uint, page int, size int) ([]*models.Comment, int64, error) {
	var comments []*models.Comment
	var total int64

	offset := (page - 1) * size

	DB.Model(&models.Comment{}).Where("post_id = ?", postID).Count(&total)

	// 查询分页数据，并预加载 User 信息以避免 N+1 查询
	err := DB.Where("post_id = ?", postID).
		Order("created_at ASC"). // 按创建时间升序
		Limit(size).
		Offset(offset).
		Preload("User"). // 关键！预加载作者信息
		Find(&comments).Error

	return comments, total, err
}

func CreateComment(tx *gorm.DB, comment *models.Comment) error {
	return tx.Create(comment).Error
}

// tool for update
func UpdatePostCounter(tx *gorm.DB, postID uint, column string, amount int) error {
	// gorm.Expr能保证原子性
	// TODO:搞不懂
	return tx.Model(&models.Post{}).Where("id = ?", postID).
		Update(column, gorm.Expr(column+" + ?", amount)).Error
}

func GetCommentById(commentId uint) (*models.Comment, error) {
	var comment models.Comment
	err := DB.Where("id=?", commentId).Preload("User").First(&comment).Error
	if err != nil {
		return nil, err
	}
	return &comment, nil
}

func DeleteComment(postId uint) error {
	return DB.Delete(&models.Comment{}, postId).Error
}

func UpdateComment(comment *models.Comment) error {
	res := DB.Model(comment).Where("id=?", comment.ID).Updates(comment)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return errors.New("更新失败")
	}
	res.First(comment)
	return nil
}

// --------------------点赞、收藏------------------------------
// IsLiked 检查用户是否已点赞某帖子
// 思路：检查中间表数量，不差出模型
func IsLiked(userID, postID uint) (bool, error) {
	var count int64
	// 直接在中间表上执行 COUNT 查询
	// 我们甚至不需要 .Model()，因为 Count() 不需要模型来确定表名
	err := DB.Table("user_post_likes"). // 使用我们之前定义的常量
						Where("user_id = ? AND post_id = ?", userID, postID).
						Count(&count).Error

	if err != nil {
		// 如果在 COUNT 过程中发生任何数据库错误，直接返回
		return false, err
	}

	// 如果 count > 0，说明存在，返回 true。否则返回 false。
	return count > 0, nil
}

func IsFavorite(userId, postId uint) (bool, error) {
	var count int64
	err := DB.Table("user_post_favorites").
		Where("user_id = ? AND post_id = ?", userId, postId).
		Count(&count).Error
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// RemoveLike 在事务中移除点赞关系
func RemoveLike(tx *gorm.DB, userID, postID uint) error {
	// .Delete() 方法需要一个模型实例来确定表名
	// 如果不提供模型，我们就需要用原生 SQL 或者指定 Table

	// 方法一：使用 Table() 和 Where() (推荐)
	// 注意，Delete 需要一个“模板”来知道删除什么，但内容不重要
	return tx.Table("user_post_likes").
		Where("user_id = ? AND post_id = ?", userID, postID).
		Delete(map[string]any{}).Error // 传入一个空的 map 即可

	// 方法二：原生 SQL
	// sql := "DELETE FROM user_post_likes WHERE user_id = ? AND post_id = ?"
	// return tx.Exec(sql, userID, postID).Error
}

// AddLike 在事务中添加点赞关系
func AddLike(tx *gorm.DB, userID, postID uint) error {
	// 当没有模型时，我们无法使用 .Create()
	// 但我们可以用 .Exec() 来执行原生的 SQL 语句
	// 或者，更 GORM-style 的方式是使用 .Create() 配合 map

	// 方法一：使用 map (推荐)
	// GORM 允许你用 map[string]interface{} 来创建记录
	likeRecord := map[string]any{
		"user_id": userID,
		"post_id": postID,
		// 如果你的表有 created_at，GORM 配合 map 时可能不会自动填充
		// "created_at": time.Now(), // 可能需要手动指定
	}
	return tx.Table("user_post_likes").Create(&likeRecord).Error

	// 方法二：原生 SQL (也非常清晰)
	// sql := "INSERT INTO user_post_likes (user_id, post_id) VALUES (?, ?)"
	// return tx.Exec(sql, userID, postID).Error
}

func RemoveFavorite(tx *gorm.DB, userID, postID uint) error {
	return tx.Table("user_post_favorites").
		Where("user_id = ? AND post_id = ?", userID, postID).
		Delete(map[string]any{}).Error // 传入一个空的 map 即可
}

func AddFavorite(tx *gorm.DB, userID, postID uint) error {
	favRecord := map[string]any{
		"user_id": userID,
		"post_id": postID,
	}
	return tx.Table("user_post_favorites").Create(&favRecord).Error
}

// ------------------头像--------------------------------
// UpdateUserAvatar 更新指定用户的头像 URL
func UpdateUserAvatar(userID uint, avatarURL string) error {
	// 使用 Model 和 Where 来定位用户，并用 Update 更新单个字段
	// 这是最高效的方式
	return DB.Model(&models.User{}).Where("id = ?", userID).Update("avatar", avatarURL).Error
}
