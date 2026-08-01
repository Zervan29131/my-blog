package service

import (
	"context"
	"fmt"

	"gorm.io/gorm"

	"personal-blog/backend/internal/model"
)

type DashboardService struct {
	database *gorm.DB
}

type DashboardStats struct {
	ArticleTotal     int64 `gorm:"column:article_total" json:"article_total"`
	ArticlePublished int64 `gorm:"column:article_published" json:"article_published"`
	ArticleDraft     int64 `gorm:"column:article_draft" json:"article_draft"`
	CommentPending   int64 `gorm:"column:comment_pending" json:"comment_pending"`
	CommentApproved  int64 `gorm:"column:comment_approved" json:"comment_approved"`
}

func NewDashboardService(database *gorm.DB) *DashboardService {
	return &DashboardService{database: database}
}

func (service *DashboardService) Stats(ctx context.Context) (DashboardStats, error) {
	var stats DashboardStats
	result := service.database.WithContext(ctx).Raw(`
		SELECT
			(SELECT COUNT(*) FROM articles) AS article_total,
			(SELECT COUNT(*) FROM articles WHERE status = ?) AS article_published,
			(SELECT COUNT(*) FROM articles WHERE status = ?) AS article_draft,
			(SELECT COUNT(*) FROM comments WHERE status = ?) AS comment_pending,
			(SELECT COUNT(*) FROM comments WHERE status = ?) AS comment_approved
	`,
		model.ArticleStatusPublished,
		model.ArticleStatusDraft,
		model.CommentStatusPending,
		model.CommentStatusApproved,
	).Scan(&stats)
	if result.Error != nil {
		return DashboardStats{}, fmt.Errorf("load dashboard statistics: %w", result.Error)
	}
	return stats, nil
}
