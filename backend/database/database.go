package database

import (
	"log"
	"my-blog-backend/config"
	"my-blog-backend/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

// InitDB 初始化数据库连接
func InitDB(cfg *config.Config) {
	var err error
	dsn := cfg.Database.Source

	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("无法连接数据库: %v", err)
	}

	log.Println("数据库连接成功")

	// 自动迁移模式 (Auto Migrate)
	// 注意：这里引用了我们在上一步设计的 models 包
	// 如果你还没有 models 包，请先注释掉这一行
	err = DB.AutoMigrate(&models.User{}, &models.Post{}, &models.Category{})
	if err != nil {
		log.Fatalf("数据库迁移失败: %v", err)
	}
}
