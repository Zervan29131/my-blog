package handler

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"personal-blog/backend/internal/service"
)

const (
	defaultCommentPageSize = 20
	maximumCommentPageSize = 50
)

type CommentHandler struct {
	commentService *service.CommentService
}

type commentRequest struct {
	Nickname string          `json:"nickname"`
	Email    string          `json:"email"`
	Content  string          `json:"content"`
	Status   json.RawMessage `json:"status"`
}

type commentStatusRequest struct {
	Status string `json:"status"`
}

func NewCommentHandler(commentService *service.CommentService) *CommentHandler {
	return &CommentHandler{commentService: commentService}
}

func (handler *CommentHandler) Submit(context *gin.Context) {
	var request commentRequest
	if err := context.ShouldBindJSON(&request); err != nil || len(request.Status) > 0 {
		writeError(context, http.StatusBadRequest, "VALIDATION_ERROR", "请求参数不正确")
		return
	}

	comment, err := handler.commentService.Submit(
		context.Request.Context(),
		context.Param("slug"),
		service.CommentInput{
			Nickname: request.Nickname,
			Email:    request.Email,
			Content:  request.Content,
		},
	)
	if err != nil {
		handler.writeCommentError(context, err)
		return
	}

	context.JSON(http.StatusCreated, gin.H{
		"data": gin.H{
			"id":         comment.ID,
			"nickname":   comment.Nickname,
			"content":    comment.Content,
			"status":     comment.Status,
			"created_at": comment.CreatedAt,
		},
		"message": "评论已提交，审核通过后将会显示。",
	})
}

func (handler *CommentHandler) ListApproved(context *gin.Context) {
	page, pageSize, ok := parseCommentPagination(context)
	if !ok {
		writeError(context, http.StatusBadRequest, "VALIDATION_ERROR", "请求参数不正确")
		return
	}

	result, err := handler.commentService.ListApproved(
		context.Request.Context(),
		context.Param("slug"),
		page,
		pageSize,
	)
	if err != nil {
		handler.writeCommentError(context, err)
		return
	}

	items := make([]gin.H, 0, len(result.Items))
	for _, comment := range result.Items {
		items = append(items, gin.H{
			"id":         comment.ID,
			"nickname":   comment.Nickname,
			"content":    comment.Content,
			"created_at": comment.CreatedAt,
		})
	}
	writeSuccess(context, http.StatusOK, gin.H{
		"items":       items,
		"page":        result.Page,
		"page_size":   result.PageSize,
		"total":       result.Total,
		"total_pages": result.TotalPages,
	})
}

func (handler *CommentHandler) ListAdmin(context *gin.Context) {
	page, pageSize, ok := parseCommentPagination(context)
	if !ok {
		writeError(context, http.StatusBadRequest, "VALIDATION_ERROR", "请求参数不正确")
		return
	}

	result, err := handler.commentService.ListAdmin(
		context.Request.Context(),
		context.Query("status"),
		page,
		pageSize,
	)
	if err != nil {
		handler.writeCommentError(context, err)
		return
	}

	items := make([]gin.H, 0, len(result.Items))
	for _, comment := range result.Items {
		items = append(items, gin.H{
			"id": comment.ID,
			"article": gin.H{
				"id":    comment.ArticleID,
				"title": comment.ArticleTitle,
				"slug":  comment.ArticleSlug,
			},
			"nickname":   comment.Nickname,
			"email":      comment.Email,
			"content":    comment.Content,
			"status":     comment.Status,
			"created_at": comment.CreatedAt,
			"updated_at": comment.UpdatedAt,
		})
	}
	writeSuccess(context, http.StatusOK, gin.H{
		"items":       items,
		"page":        result.Page,
		"page_size":   result.PageSize,
		"total":       result.Total,
		"total_pages": result.TotalPages,
	})
}

func (handler *CommentHandler) UpdateStatus(context *gin.Context) {
	id, ok := parseCommentID(context)
	if !ok {
		return
	}

	var request commentStatusRequest
	if err := context.ShouldBindJSON(&request); err != nil {
		writeError(context, http.StatusBadRequest, "VALIDATION_ERROR", "请求参数不正确")
		return
	}
	if err := handler.commentService.UpdateStatus(context.Request.Context(), id, request.Status); err != nil {
		handler.writeCommentError(context, err)
		return
	}

	writeSuccess(context, http.StatusOK, gin.H{"id": id, "status": request.Status})
}

func (handler *CommentHandler) Delete(context *gin.Context) {
	id, ok := parseCommentID(context)
	if !ok {
		return
	}
	if err := handler.commentService.Delete(context.Request.Context(), id); err != nil {
		handler.writeCommentError(context, err)
		return
	}
	context.Status(http.StatusNoContent)
}

func (handler *CommentHandler) writeCommentError(context *gin.Context, err error) {
	switch {
	case errors.Is(err, service.ErrInvalidComment):
		writeError(context, http.StatusBadRequest, "VALIDATION_ERROR", "请求参数不正确")
	case errors.Is(err, service.ErrArticleNotFound):
		writeError(context, http.StatusNotFound, "ARTICLE_NOT_FOUND", "文章不存在")
	case errors.Is(err, service.ErrCommentNotFound):
		writeError(context, http.StatusNotFound, "COMMENT_NOT_FOUND", "评论不存在")
	default:
		slog.Error("comment request failed", "error", err)
		writeError(context, http.StatusInternalServerError, "INTERNAL_ERROR", "服务器内部错误")
	}
}

func parseCommentPagination(context *gin.Context) (int, int, bool) {
	page := 1
	pageSize := defaultCommentPageSize
	var err error

	if rawPage := context.Query("page"); rawPage != "" {
		page, err = strconv.Atoi(rawPage)
		if err != nil || page <= 0 {
			return 0, 0, false
		}
	}
	if rawPageSize := context.Query("page_size"); rawPageSize != "" {
		pageSize, err = strconv.Atoi(rawPageSize)
		if err != nil || pageSize <= 0 {
			return 0, 0, false
		}
		if pageSize > maximumCommentPageSize {
			pageSize = maximumCommentPageSize
		}
	}
	return page, pageSize, true
}

func parseCommentID(context *gin.Context) (uint64, bool) {
	id, err := strconv.ParseUint(context.Param("id"), 10, 64)
	if err != nil || id == 0 {
		writeError(context, http.StatusBadRequest, "VALIDATION_ERROR", "请求参数不正确")
		return 0, false
	}
	return id, true
}
