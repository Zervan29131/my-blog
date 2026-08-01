package handler

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"

	"personal-blog/backend/internal/service"
)

type DashboardHandler struct {
	dashboardService *service.DashboardService
}

func NewDashboardHandler(dashboardService *service.DashboardService) *DashboardHandler {
	return &DashboardHandler{dashboardService: dashboardService}
}

func (handler *DashboardHandler) Stats(context *gin.Context) {
	stats, err := handler.dashboardService.Stats(context.Request.Context())
	if err != nil {
		slog.Error("load dashboard statistics failed", "error", err)
		writeError(context, http.StatusInternalServerError, "INTERNAL_ERROR", "服务器内部错误")
		return
	}

	writeSuccess(context, http.StatusOK, gin.H{
		"article_total":     stats.ArticleTotal,
		"article_published": stats.ArticlePublished,
		"article_draft":     stats.ArticleDraft,
		"comment_pending":   stats.CommentPending,
		"comment_approved":  stats.CommentApproved,
	})
}
