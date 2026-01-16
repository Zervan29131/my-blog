package controllers

import (
	"my-blog-backend/database"
	"my-blog-backend/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// CreateComment 发表评论
func CreateComment(c *gin.Context) {
	var comment models.Comment
	if err := c.ShouldBindJSON(&comment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := database.DB.Create(&comment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "评论失败"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": comment, "message": "评论成功"})
}

// GetComments GetCommentsByPost 获取某文章的评论
func GetComments(c *gin.Context) {
	postIDStr := c.Query("post_id")
	postID, _ := strconv.Atoi(postIDStr)

	var comments []models.Comment
	// 查询该文章下的所有评论，按时间倒序
	// 这里简化处理：前端拿到扁平列表后自己组装树形结构，或者后端组装。
	// 为了灵活性，我们返回扁平列表。
	if err := database.DB.Where("post_id = ?", postID).Order("created_at desc").Find(&comments).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取评论失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": comments})
}

//#### 3. 注册路由 (`backend/routes/routes.go`)
//在 `v1` 组的**公开接口**区域添加（评论通常允许游客查看和发布）：
//
//// ...
//v1.GET("/comments", controllers.GetComments)
//v1.POST("/comments", controllers.CreateComment)
//// ...
