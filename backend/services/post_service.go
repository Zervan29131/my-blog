package services

import (
	"my-blog-backend/database"
	"my-blog-backend/models"
)

// CreatePost 创建文章
func CreatePost(post *models.Post) error {
	result := database.DB.Create(post)
	return result.Error
}

// GetPostByID 根据ID获取文章
func GetPostByID(id uint) (*models.Post, error) {
	var post models.Post
	// Preload 加载关联的 Category 信息
	result := database.DB.Preload("Category").First(&post, id)
	return &post, result.Error
}

// GetPostList 获取文章列表 (分页)
// page: 当前页码, pageSize: 每页数量
func GetPostList(page int, pageSize int) ([]models.Post, int64, error) {
	var posts []models.Post
	var total int64

	// 计算偏移量
	offset := (page - 1) * pageSize

	// 1. 获取总数
	if err := database.DB.Model(&models.Post{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 2. 获取分页数据，按创建时间倒序
	// Preload("Category") 确保返回分类信息
	// Omit("Content") 列表页通常不需要返回长文章内容，优化性能
	err := database.DB.Omit("Content").Preload("Category").Offset(offset).Limit(pageSize).Order("created_at desc").Find(&posts).Error

	return posts, total, err
}

// UpdatePost 更新文章
func UpdatePost(id uint, updateData *models.Post) error {
	var post models.Post
	if err := database.DB.First(&post, id).Error; err != nil {
		return err
	}

	// 更新指定字段
	return database.DB.Model(&post).Updates(updateData).Error
}

// DeletePost 删除文章
func DeletePost(id uint) error {
	return database.DB.Delete(&models.Post{}, id).Error
}
