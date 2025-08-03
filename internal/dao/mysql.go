package dao

import (
	. "Nuxus/configs"
	"Nuxus/internal/models"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitMysql() {
	// 读取配置
	dsn := Conf.MySQL.DSN

	// 连接数据库
	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Connect to mysql err: %v", err)
	}

	// 自动迁移
	err = DB.AutoMigrate(&models.User{}, &models.Post{}, &models.Tag{}, &models.Comment{})
	if err != nil {
		log.Fatalf("Failed to auto migrate err: %v", err)
	}

	log.Println("Database connection and migration successful!")
}
