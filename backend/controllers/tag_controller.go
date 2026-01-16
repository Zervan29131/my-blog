package controllers

import (
	"my-blog-backend/database"
	"my-blog-backend/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// CreateTag 创建标签
func CreateTag(c *gin.Context) {
	var tag models.Tag
	// 绑定 JSON
	if err := c.ShouldBindJSON(&tag); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 写入数据库
	if err := database.DB.Create(&tag).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建标签失败，可能标签名已存在"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "创建成功",
		"data":    tag,
	})
}

// GetTagList 获取所有标签
func GetTagList(c *gin.Context) {
	var tags []models.Tag
	// 标签数量通常不多，直接查询全部返回即可，方便前端做下拉选择
	if err := database.DB.Find(&tags).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取标签失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": tags})
}

// DeleteTag 删除标签
func DeleteTag(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.Atoi(idStr)

	// 硬删除
	if err := database.DB.Unscoped().Delete(&models.Tag{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}
