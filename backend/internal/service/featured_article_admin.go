package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"

	"personal-blog/backend/internal/model"
)

const maximumFeaturedArticles = 10

type FeaturedArticleAdminService struct {
	database *gorm.DB
}

type FeaturedArticleInput struct {
	ArticleID uint64
	SortOrder int
	IsVisible bool
}

type AdminFeaturedArticle struct {
	ArticleID   uint64     `json:"article_id"`
	Title       string     `json:"title"`
	Slug        string     `json:"slug"`
	Summary     string     `json:"summary"`
	Status      string     `json:"status"`
	PublishedAt *time.Time `json:"published_at"`
	SortOrder   int        `json:"sort_order"`
	IsVisible   bool       `json:"is_visible"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

func NewFeaturedArticleAdminService(database *gorm.DB) *FeaturedArticleAdminService {
	return &FeaturedArticleAdminService{database: database}
}

func (service *FeaturedArticleAdminService) List(ctx context.Context) ([]AdminFeaturedArticle, error) {
	items := make([]AdminFeaturedArticle, 0)
	result := service.database.WithContext(ctx).
		Table("featured_articles").
		Select(
			"featured_articles.article_id, articles.title, articles.slug, articles.summary, articles.status, " +
				"articles.published_at, featured_articles.sort_order, featured_articles.is_visible, " +
				"featured_articles.created_at, featured_articles.updated_at",
		).
		Joins("JOIN articles ON articles.id = featured_articles.article_id").
		Order("featured_articles.sort_order ASC, featured_articles.id ASC").
		Scan(&items)
	if result.Error != nil {
		return nil, fmt.Errorf("list administrator featured articles: %w", result.Error)
	}
	return items, nil
}

func (service *FeaturedArticleAdminService) Create(
	ctx context.Context,
	input FeaturedArticleInput,
	administratorID uint64,
) (model.FeaturedArticle, error) {
	if input.ArticleID == 0 {
		return model.FeaturedArticle{}, ErrInvalidSiteContent
	}
	featured := model.FeaturedArticle{
		ArticleID: input.ArticleID, SortOrder: input.SortOrder, IsVisible: input.IsVisible,
	}
	err := service.database.WithContext(ctx).Transaction(func(transaction *gorm.DB) error {
		var article model.Article
		result := transaction.Where("id = ?", input.ArticleID).Limit(1).Find(&article)
		if result.Error != nil {
			return fmt.Errorf("find featured article candidate: %w", result.Error)
		}
		if result.RowsAffected == 0 || article.Status != model.ArticleStatusPublished {
			return ErrArticleNotPublished
		}

		var duplicateCount int64
		if err := transaction.Model(&model.FeaturedArticle{}).
			Where("article_id = ?", input.ArticleID).
			Count(&duplicateCount).Error; err != nil {
			return fmt.Errorf("check featured article duplicate: %w", err)
		}
		if duplicateCount > 0 {
			return ErrSiteContentConflict
		}

		var total int64
		if err := transaction.Model(&model.FeaturedArticle{}).Count(&total).Error; err != nil {
			return fmt.Errorf("count featured articles: %w", err)
		}
		if total >= maximumFeaturedArticles {
			return ErrSiteContentLimit
		}
		if err := transaction.Create(&featured).Error; err != nil {
			if errors.Is(err, gorm.ErrDuplicatedKey) {
				return ErrSiteContentConflict
			}
			return fmt.Errorf("create featured article: %w", err)
		}
		return nil
	})
	if err != nil {
		return model.FeaturedArticle{}, err
	}
	logSiteContentOperation("featured_article_created", administratorID, "article", input.ArticleID)
	return featured, nil
}

func (service *FeaturedArticleAdminService) Delete(
	ctx context.Context,
	articleID uint64,
	administratorID uint64,
) error {
	result := service.database.WithContext(ctx).
		Where("article_id = ?", articleID).
		Delete(&model.FeaturedArticle{})
	if result.Error != nil {
		return fmt.Errorf("delete featured article: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return ErrSiteContentNotFound
	}
	logSiteContentOperation("featured_article_deleted", administratorID, "article", articleID)
	return nil
}

func (service *FeaturedArticleAdminService) UpdateVisibility(
	ctx context.Context,
	articleID uint64,
	isVisible bool,
	administratorID uint64,
) (model.FeaturedArticle, error) {
	var featured model.FeaturedArticle
	result := service.database.WithContext(ctx).
		Where("article_id = ?", articleID).
		Limit(1).
		Find(&featured)
	if result.Error != nil {
		return model.FeaturedArticle{}, fmt.Errorf("find featured article: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return model.FeaturedArticle{}, ErrSiteContentNotFound
	}

	featured.IsVisible = isVisible
	featured.UpdatedAt = time.Now().UTC()
	update := service.database.WithContext(ctx).
		Model(&model.FeaturedArticle{}).
		Where("article_id = ?", articleID).
		Updates(map[string]any{"is_visible": isVisible, "updated_at": featured.UpdatedAt})
	if update.Error != nil {
		return model.FeaturedArticle{}, fmt.Errorf("update featured article visibility: %w", update.Error)
	}
	if update.RowsAffected == 0 {
		return model.FeaturedArticle{}, ErrSiteContentNotFound
	}
	logSiteContentOperation("featured_article_visibility_updated", administratorID, "article", articleID)
	return featured, nil
}

func (service *FeaturedArticleAdminService) Reorder(
	ctx context.Context,
	items []ResourceOrder,
	administratorID uint64,
) error {
	if err := validateResourceOrder(items, maximumFeaturedArticles); err != nil {
		return err
	}
	now := time.Now().UTC()
	err := service.database.WithContext(ctx).Transaction(func(transaction *gorm.DB) error {
		for _, item := range items {
			result := transaction.Model(&model.FeaturedArticle{}).
				Where("article_id = ?", item.ID).
				Updates(map[string]any{"sort_order": item.SortOrder, "updated_at": now})
			if result.Error != nil {
				return fmt.Errorf("reorder featured article: %w", result.Error)
			}
			if result.RowsAffected == 0 {
				return ErrSiteContentNotFound
			}
		}
		return nil
	})
	if err != nil {
		return err
	}
	logSiteContentOperation("featured_articles_reordered", administratorID, "featured_articles", 0)
	return nil
}
