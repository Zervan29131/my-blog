package controllers

import (
	"my-blog-backend/models"
	"my-blog-backend/services"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// CreatePost 处理创建文章请求
func CreatePost(c *gin.Context) {
	var post models.Post

	// 1. 绑定前端传来的 JSON 数据 (Title, Content, CategoryID 等)
	if err := c.ShouldBindJSON(&post); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数格式错误: " + err.Error()})
		return
	}

	// 2. 🟢 关键修复：获取当前登录用户的 ID
	// 这个值是在 middleware/jwt_auth.go 中设置的 c.Set("userID", ...)
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未获取到用户信息，请重新登录"})
		return
	}

	// 将 interface{} 类型的 userID 断言为 uint 并赋值给 post.AuthorID
	// 这样数据库就知道这篇文章是谁写的了
	post.AuthorID = userID.(uint)

	// 3. 自动生成摘要 (如果前端没传)
	if post.Summary == "" && len(post.Content) > 0 {
		// 将内容转为 rune 切片以正确处理中文字符计数
		runes := []rune(post.Content)
		if len(runes) > 100 {
			post.Summary = string(runes[:100]) + "..."
		} else {
			post.Summary = string(runes)
		}
	}

	// 4. 调用 Service 层写入数据库
	if err := services.CreatePost(&post); err != nil {
		// 建议在控制台打印错误，方便调试数据库报错 (如 Error 1452 外键错误)
		println("CreatePost Error:", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建文章失败"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "文章发布成功",
		"data":    post,
	})
}

// GetPostDetail 获取文章详情
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

// GetPostList 获取文章列表
func GetPostList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	// 🟢 获取搜索关键词 q
	keyword := c.Query("q")

	// 🟢 将 keyword 传给 Service
	posts, total, err := services.GetPostList(page, pageSize, keyword)
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

// UpdatePost 更新文章
func UpdatePost(c *gin.Context) {
	// 1. 获取 URL 中的 ID 并去除可能的空格
	idStr := strings.TrimSpace(c.Param("id"))
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的 ID 格式"})
		return
	}

	// 2. 绑定要更新的数据
	var updateData models.Post
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 3. 调用 Service 更新
	if err := services.UpdatePost(uint(id), &updateData); err != nil {
		if strings.Contains(err.Error(), "record not found") {
			c.JSON(http.StatusNotFound, gin.H{"error": "文章不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "更新成功"})
}

// DeletePost 删除文章
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
