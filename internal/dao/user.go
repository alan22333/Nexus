package dao

import (
	"Nuxus/internal/models"
	"strings"

	"gorm.io/gorm"
)

type UserDAO struct {
	db *gorm.DB
}

func NewUserDAO(db *gorm.DB) *UserDAO {
	return &UserDAO{db: db}
}

func (u *UserDAO) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	err := u.db.Where("email=?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (u *UserDAO) GetUserById(userId uint) (*models.User, error) {
	var user models.User
	err := u.db.Where("id=?", userId).First(&user).Error
	return &user, err
}

func (u *UserDAO) CreateUser(user *models.User) (*models.User, error) {
	err := u.db.Create(user).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (u *UserDAO) GetUserByIdentifier(identifier string) (*models.User, error) {
	var user models.User
	var err error

	// 施展探查法术：检查是否包含“@”符文
	if strings.Contains(identifier, "@") {
		err = u.db.Where("email = ?", identifier).First(&user).Error
	} else {
		err = u.db.Where("username = ?", identifier).First(&user).Error
	}

	return &user, err
}

func (u *UserDAO) UpdateUserPassword(id uint, password string) error {
	return u.db.Model(&models.User{}).Where("id=?", id).Update("password", password).Error
}

func (u *UserDAO) UpdateProfile(user *models.User) (*models.User, error) {
	err := u.db.Model(user).Where("id=?", user.ID).Updates(user).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}

// ------------------头像--------------------------------
// UpdateUserAvatar 更新指定用户的头像 URL
func (u *UserDAO) UpdateUserAvatar(userID uint, avatarURL string) error {
	// 使用 Model 和 Where 来定位用户，并用 Update 更新单个字段
	// 这是最高效的方式
	return u.db.Model(&models.User{}).Where("id = ?", userID).Update("avatar", avatarURL).Error
}
