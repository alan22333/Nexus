package dao

import (
	"Nuxus/configs"
	"Nuxus/internal/models"
	"log"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// ctxTxKey 是在Context中存储事务对象的键名
// 用于在事务执行过程中传递数据库事务实例
// const ctxTxKey = "TxKey"

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func NewDB(config *configs.Config) *gorm.DB {
	// 读取配置
	dsn := config.MySQL.DSN

	// 连接数据库
	var err error
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Connect to mysql err: %v", err)
	}

	// 自动迁移
	err = db.AutoMigrate(&models.User{}, &models.Post{}, &models.Tag{}, &models.Comment{})
	if err != nil {
		log.Fatalf("Failed to auto migrate err: %v", err)
	}

	// 配置数据库连接池
	// 连接池可以提高数据库访问性能并控制资源使用
	sqlDB, err := db.DB()
	if err != nil {
		panic(err)
	}
	sqlDB.SetMaxIdleConns(10)           // 设置空闲连接池中连接的最大数量
	sqlDB.SetMaxOpenConns(100)          // 设置打开数据库连接的最大数量
	sqlDB.SetConnMaxLifetime(time.Hour) // 设置连接可复用的最大时间

	log.Println("Database connection and migration successful!")
	return db
}

// DB 获取数据库连接实例
// 如果当前上下文中存在事务，则返回事务连接；否则返回普通连接
// 这是Repository模式的核心方法，确保在事务和非事务场景下都能正确获取DB实例

func (r *Repository) DB() *gorm.DB {
	return r.db
}

// func (r *Repository) Transaction(ctx context.Context, fn func(ctx context.Context) error) error {
// 	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
// 		// 将事务实例存储到上下文中，供后续DB()方法使用
// 		ctx = context.WithValue(ctx, ctxTxKey, tx)
// 		return fn(ctx)
// 	})
// }

// func InitMysql() {
// 	// 读取配置
// 	dsn := Conf.MySQL.DSN

// 	// 连接数据库
// 	var err error
// 	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
// 	if err != nil {
// 		log.Fatalf("Connect to mysql err: %v", err)
// 	}

// 	// 自动迁移
// 	err = DB.AutoMigrate(&models.User{}, &models.Post{}, &models.Tag{}, &models.Comment{})
// 	if err != nil {
// 		log.Fatalf("Failed to auto migrate err: %v", err)
// 	}

// 	log.Println("Database connection and migration successful!")
// }
