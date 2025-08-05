package dao

import (
	"Nuxus/internal/models"
	"strings"
)

func GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	err := DB.Where("email=?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func GetUserById(userId uint) (*models.User, error) {
	var user models.User
	err := DB.Where("id=?", userId).First(&user).Error
	return &user, err
}

func CreateUser(user *models.User) (*models.User, error) {
	err := DB.Create(user).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}

func GetUserByIdentifier(identifier string) (*models.User, error) {
	var user models.User
	var err error

	// 施展探查法术：检查是否包含“@”符文
	if strings.Contains(identifier, "@") {
		err = DB.Where("email = ?", identifier).First(&user).Error
	} else {
		err = DB.Where("username = ?", identifier).First(&user).Error
	}

	return &user, err
}

func UpdateUserPassword(id uint, password string) error {
	return DB.Model(&models.User{}).Where("id=?", id).Update("password", password).Error
}

func UpdateProfile(user *models.User)(*models.User,error){
	err := DB.Model(user).Where("id=?", user.ID).Updates(user).Error
	if err != nil {
		return nil,err
	}
	return user,nil
}
