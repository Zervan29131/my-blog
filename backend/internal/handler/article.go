package handler

import (
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"personal-blog/backend/internal/model"
	"personal-blog/backend/internal/service"
)

const (
	defaultArticlePageSize = 10
	maximumArticlePageSize = 50
)

type ArticleHandler struct {
	articleService *service.ArticleService
}

type articleRequest struct {
	Title   string `json:"title"`
	Slug    string `json:"slug"`
	Summary string `json:"summary"`
	Content string `json:"content"`
	Status  string `json:"status"`
}

func NewArticleHandler(articleService *service.ArticleService) *ArticleHandler {
	return &ArticleHandler{articleService: articleService}
}

func (handler *ArticleHandler) ListPublished(context *gin.Context) {
	page, pageSize, ok := parseArticlePagination(context)
	if !ok {
		writeError(context, http.StatusBadRequest, "VALIDATION_ERROR", "请求参数不正确")
		return
	}

	result, err := handler.articleService.ListPublished(context.Request.Context(), page, pageSize)
	if err != nil {
		slog.Error("list published articles failed", "error", err)
		writeError(context, http.StatusInternalServerError, "INTERNAL_ERROR", "服务器内部错误")
		return
	}

	items := make([]gin.H, 0, len(result.Items))
	for _, article := range result.Items {
		items = append(items, gin.H{
			"id":            article.ID,
			"title":         article.Title,
			"slug":          article.Slug,
			"summary":       article.Summary,
			"published_at":  article.PublishedAt,
			"comment_count": article.CommentCount,
		})
	}
	writeArticlePage(context, items, result)
}

func (handler *ArticleHandler) GetPublished(context *gin.Context) {
	article, err := handler.articleService.GetPublishedBySlug(
		context.Request.Context(),
		context.Param("slug"),
	)
	if err != nil {
		handler.writeArticleError(context, err)
		return
	}

	writeSuccess(context, http.StatusOK, gin.H{
		"id":           article.ID,
		"title":        article.Title,
		"slug":         article.Slug,
		"summary":      article.Summary,
		"content":      article.Content,
		"published_at": article.PublishedAt,
		"created_at":   article.CreatedAt,
		"updated_at":   article.UpdatedAt,
	})
}

func (handler *ArticleHandler) ListAdmin(context *gin.Context) {
	page, pageSize, ok := parseArticlePagination(context)
	if !ok {
		writeError(context, http.StatusBadRequest, "VALIDATION_ERROR", "请求参数不正确")
		return
	}

	status := context.Query("status")
	if status != "" && status != model.ArticleStatusPublished && status != model.ArticleStatusDraft {
		writeError(context, http.StatusBadRequest, "VALIDATION_ERROR", "请求参数不正确")
		return
	}

	result, err := handler.articleService.ListAllByStatus(context.Request.Context(), page, pageSize, status)
	if err != nil {
		slog.Error("list administrator articles failed", "error", err)
		writeError(context, http.StatusInternalServerError, "INTERNAL_ERROR", "服务器内部错误")
		return
	}

	items := make([]gin.H, 0, len(result.Items))
	for _, article := range result.Items {
		items = append(items, gin.H{
			"id":           article.ID,
			"title":        article.Title,
			"slug":         article.Slug,
			"summary":      article.Summary,
			"status":       article.Status,
			"published_at": article.PublishedAt,
			"created_at":   article.CreatedAt,
			"updated_at":   article.UpdatedAt,
		})
	}
	writeArticlePage(context, items, result)
}

func (handler *ArticleHandler) GetAdmin(context *gin.Context) {
	id, ok := parseArticleID(context)
	if !ok {
		return
	}

	article, err := handler.articleService.GetByID(context.Request.Context(), id)
	if err != nil {
		handler.writeArticleError(context, err)
		return
	}
	writeSuccess(context, http.StatusOK, fullArticleResponse(article))
}

func (handler *ArticleHandler) Create(context *gin.Context) {
	var request articleRequest
	if err := context.ShouldBindJSON(&request); err != nil {
		writeError(context, http.StatusBadRequest, "VALIDATION_ERROR", "请求参数不正确")
		return
	}

	article, err := handler.articleService.Create(context.Request.Context(), request.toServiceInput())
	if err != nil {
		handler.writeArticleError(context, err)
		return
	}
	writeSuccess(context, http.StatusCreated, fullArticleResponse(article))
}

func (handler *ArticleHandler) Update(context *gin.Context) {
	id, ok := parseArticleID(context)
	if !ok {
		return
	}

	var request articleRequest
	if err := context.ShouldBindJSON(&request); err != nil {
		writeError(context, http.StatusBadRequest, "VALIDATION_ERROR", "请求参数不正确")
		return
	}

	article, err := handler.articleService.Update(context.Request.Context(), id, request.toServiceInput())
	if err != nil {
		handler.writeArticleError(context, err)
		return
	}
	writeSuccess(context, http.StatusOK, fullArticleResponse(article))
}

func (handler *ArticleHandler) Delete(context *gin.Context) {
	id, ok := parseArticleID(context)
	if !ok {
		return
	}

	if err := handler.articleService.Delete(context.Request.Context(), id); err != nil {
		handler.writeArticleError(context, err)
		return
	}
	context.Status(http.StatusNoContent)
}

func (request articleRequest) toServiceInput() service.ArticleInput {
	return service.ArticleInput{
		Title:   request.Title,
		Slug:    request.Slug,
		Summary: request.Summary,
		Content: request.Content,
		Status:  request.Status,
	}
}

func (handler *ArticleHandler) writeArticleError(context *gin.Context, err error) {
	switch {
	case errors.Is(err, service.ErrInvalidArticle):
		writeError(context, http.StatusBadRequest, "VALIDATION_ERROR", "请求参数不正确")
	case errors.Is(err, service.ErrArticleNotFound):
		writeError(context, http.StatusNotFound, "ARTICLE_NOT_FOUND", "文章不存在")
	case errors.Is(err, service.ErrSlugConflict):
		writeError(context, http.StatusConflict, "SLUG_CONFLICT", "文章 Slug 已存在")
	default:
		slog.Error("article request failed", "error", err)
		writeError(context, http.StatusInternalServerError, "INTERNAL_ERROR", "服务器内部错误")
	}
}

func parseArticlePagination(context *gin.Context) (int, int, bool) {
	page := 1
	pageSize := defaultArticlePageSize
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
		if pageSize > maximumArticlePageSize {
			pageSize = maximumArticlePageSize
		}
	}

	return page, pageSize, true
}

func parseArticleID(context *gin.Context) (uint64, bool) {
	id, err := strconv.ParseUint(context.Param("id"), 10, 64)
	if err != nil || id == 0 {
		writeError(context, http.StatusBadRequest, "VALIDATION_ERROR", "请求参数不正确")
		return 0, false
	}
	return id, true
}

func fullArticleResponse(article model.Article) gin.H {
	return gin.H{
		"id":           article.ID,
		"title":        article.Title,
		"slug":         article.Slug,
		"summary":      article.Summary,
		"content":      article.Content,
		"status":       article.Status,
		"published_at": article.PublishedAt,
		"created_at":   article.CreatedAt,
		"updated_at":   article.UpdatedAt,
	}
}

func writeArticlePage(context *gin.Context, items []gin.H, result service.ArticlePage) {
	writeSuccess(context, http.StatusOK, gin.H{
		"items":       items,
		"page":        result.Page,
		"page_size":   result.PageSize,
		"total":       result.Total,
		"total_pages": result.TotalPages,
	})
}
