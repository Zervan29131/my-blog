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

// GetPostList 获取文章列表 (支持分页和关键词搜索)
// page: 当前页码, pageSize: 每页数量, keyword: 搜索关键词(可选)
func GetPostList(page int, pageSize int, keyword string) ([]models.Post, int64, error) {
	var posts []models.Post
	var total int64

	offset := (page - 1) * pageSize

	// 1. 创建基础查询
	tx := database.DB.Model(&models.Post{})

	// 2. 如果有关键词，添加模糊查询条件 (匹配标题、摘要或内容)
	if keyword != "" {
		likeStr := "%" + keyword + "%"
		tx = tx.Where("title LIKE ? OR summary LIKE ? OR content LIKE ?", likeStr, likeStr, likeStr)
	}

	// 3. 计算总数 (基于筛选条件)
	if err := tx.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 4. 获取数据 (预加载关联表)
	// Omit("Content") 列表页不查长文本，优化性能
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
	// 先检查是否存在
	if err := database.DB.First(&post, id).Error; err != nil {
		return err
	}

	// 更新关联关系 (Tags) 需要特殊处理
	// 如果 updateData.Tags 不为空，GORM 的 Association Mode 会处理中间表
	if len(updateData.Tags) > 0 {
		if err := database.DB.Model(&post).Association("Tags").Replace(updateData.Tags); err != nil {
			return err
		}
	}

	// 更新普通字段
	return database.DB.Model(&post).Updates(updateData).Error
}

// DeletePost 删除文章
func DeletePost(id uint) error {
	// 也是软删除，因为 Post 模型包含 gorm.Model
	return database.DB.Delete(&models.Post{}, id).Error
}
