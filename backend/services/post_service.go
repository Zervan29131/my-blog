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
	// Preload 加载关联的 Category, Author, Tags 信息
	result := database.DB.Preload("Category").Preload("Author").Preload("Tags").First(&post, id)
	return &post, result.Error
}

// GetPostList 获取文章列表 (支持分页、关键词搜索、分类筛选)
// 🟢 修改点：增加 categoryID 参数
func GetPostList(page int, pageSize int, keyword string, categoryID uint) ([]models.Post, int64, error) {
	var posts []models.Post
	var total int64

	offset := (page - 1) * pageSize

	tx := database.DB.Model(&models.Post{})

	// 🟢 新增：如果有分类ID，添加筛选条件
	if categoryID > 0 {
		tx = tx.Where("category_id = ?", categoryID)
	}

	// 搜索逻辑
	if keyword != "" {
		likeStr := "%" + keyword + "%"
		tx = tx.Where("title LIKE ? OR summary LIKE ? OR content LIKE ?", likeStr, likeStr, likeStr)
	}

	// 计算总数
	if err := tx.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取数据
	err := tx.Omit("Content").
		Preload("Category").
		Preload("Author").
		Preload("Tags").
		Offset(offset).
		Limit(pageSize).
		Order("created_at desc").
		Find(&posts).Error

	return posts, total, err
}

// UpdatePost 更新文章
func UpdatePost(id uint, updateData *models.Post) error {
	var post models.Post
	if err := database.DB.First(&post, id).Error; err != nil {
		return err
	}
	if len(updateData.Tags) > 0 {
		if err := database.DB.Model(&post).Association("Tags").Replace(updateData.Tags); err != nil {
			return err
		}
	}
	return database.DB.Model(&post).Updates(updateData).Error
}

// DeletePost 删除文章
func DeletePost(id uint) error {
	return database.DB.Delete(&models.Post{}, id).Error
}
