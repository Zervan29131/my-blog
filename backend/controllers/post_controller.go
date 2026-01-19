package controllers

import (
	"my-blog-backend/models"
	"my-blog-backend/services"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// GetPostList 获取文章列表
func GetPostList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	keyword := c.Query("q")

	// 🟢 新增：获取 category_id 参数
	categoryID, _ := strconv.Atoi(c.Query("category_id"))

	// 🟢 修改：传入 categoryID (转为 uint)
	posts, total, err := services.GetPostList(page, pageSize, keyword, uint(categoryID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取列表失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  posts,
		"total": total,
		"page":  page,
	})
}

// CreatePost, GetPostDetail, UpdatePost, DeletePost 保持不变...
// (为节省篇幅，这里只列出修改过的 GetPostList，但你需要保留其他函数)
// ⚠️ 请确保文件中包含 CreatePost, GetPostDetail, UpdatePost, DeletePost 的原有代码
// 这里为了文件完整性，我会把它们补全，防止覆盖导致丢失

func CreatePost(c *gin.Context) {
	var post models.Post
	if err := c.ShouldBindJSON(&post); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数格式错误: " + err.Error()})
		return
	}
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未获取到用户信息"})
		return
	}
	post.AuthorID = userID.(uint)
	if post.Summary == "" && len(post.Content) > 0 {
		runes := []rune(post.Content)
		if len(runes) > 100 {
			post.Summary = string(runes[:100]) + "..."
		} else {
			post.Summary = string(runes)
		}
	}
	if err := services.CreatePost(&post); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建文章失败"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "文章发布成功", "data": post})
}

func GetPostDetail(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.Atoi(idStr)
	post, err := services.GetPostByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "文章不存在"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": post})
}

func UpdatePost(c *gin.Context) {
	idStr := strings.TrimSpace(c.Param("id"))
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的 ID 格式"})
		return
	}
	var updateData models.Post
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := services.UpdatePost(uint(id), &updateData); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "更新成功"})
}

func DeletePost(c *gin.Context) {
	idStr := strings.TrimSpace(c.Param("id"))
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的 ID 格式"})
		return
	}
	if err := services.DeletePost(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}
