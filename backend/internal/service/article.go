package service

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"
	"unicode/utf8"

	"gorm.io/gorm"

	"personal-blog/backend/internal/model"
)

var (
	ErrArticleNotFound = errors.New("article not found")
	ErrInvalidArticle  = errors.New("invalid article")
	ErrSlugConflict    = errors.New("article slug already exists")
)

var slugPattern = regexp.MustCompile(`^[a-z0-9-]+$`)

type ArticleService struct {
	database *gorm.DB
}

type ArticleInput struct {
	Title   string
	Slug    string
	Summary string
	Content string
	Status  string
}

type ArticlePage struct {
	Items      []model.Article
	Page       int
	PageSize   int
	Total      int64
	TotalPages int
}

func NewArticleService(database *gorm.DB) *ArticleService {
	return &ArticleService{database: database}
}

func (service *ArticleService) Create(ctx context.Context, input ArticleInput) (model.Article, error) {
	prepared, err := prepareArticleInput(input, true)
	if err != nil {
		return model.Article{}, err
	}

	exists, err := service.slugExists(ctx, prepared.Slug, 0)
	if err != nil {
		return model.Article{}, err
	}
	if exists {
		return model.Article{}, ErrSlugConflict
	}

	article := model.Article{
		Title:   prepared.Title,
		Slug:    prepared.Slug,
		Summary: prepared.Summary,
		Content: prepared.Content,
		Status:  prepared.Status,
	}
	if article.Status == model.ArticleStatusPublished {
		now := time.Now().UTC()
		article.PublishedAt = &now
	}

	if err := service.database.WithContext(ctx).Create(&article).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return model.Article{}, ErrSlugConflict
		}
		return model.Article{}, fmt.Errorf("create article: %w", err)
	}

	return article, nil
}

func (service *ArticleService) Update(
	ctx context.Context,
	id uint64,
	input ArticleInput,
) (model.Article, error) {
	article, err := service.GetByID(ctx, id)
	if err != nil {
		return model.Article{}, err
	}

	prepared, err := prepareArticleInput(input, false)
	if err != nil {
		return model.Article{}, err
	}

	exists, err := service.slugExists(ctx, prepared.Slug, id)
	if err != nil {
		return model.Article{}, err
	}
	if exists {
		return model.Article{}, ErrSlugConflict
	}

	if prepared.Status == model.ArticleStatusPublished && article.PublishedAt == nil {
		now := time.Now().UTC()
		article.PublishedAt = &now
	}
	article.Title = prepared.Title
	article.Slug = prepared.Slug
	article.Summary = prepared.Summary
	article.Content = prepared.Content
	article.Status = prepared.Status
	article.UpdatedAt = time.Now().UTC()

	result := service.database.WithContext(ctx).
		Model(&model.Article{}).
		Where("id = ?", article.ID).
		Updates(map[string]any{
			"title":        article.Title,
			"slug":         article.Slug,
			"summary":      article.Summary,
			"content":      article.Content,
			"status":       article.Status,
			"published_at": article.PublishedAt,
			"updated_at":   article.UpdatedAt,
		})
	if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
		return model.Article{}, ErrSlugConflict
	}
	if result.Error != nil {
		return model.Article{}, fmt.Errorf("update article: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return model.Article{}, ErrArticleNotFound
	}

	return article, nil
}

func (service *ArticleService) Delete(ctx context.Context, id uint64) error {
	result := service.database.WithContext(ctx).Delete(&model.Article{}, id)
	if result.Error != nil {
		return fmt.Errorf("delete article: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return ErrArticleNotFound
	}
	return nil
}

func (service *ArticleService) GetByID(ctx context.Context, id uint64) (model.Article, error) {
	var article model.Article
	result := service.database.WithContext(ctx).First(&article, id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return model.Article{}, ErrArticleNotFound
	}
	if result.Error != nil {
		return model.Article{}, fmt.Errorf("find article: %w", result.Error)
	}
	return article, nil
}

func (service *ArticleService) GetPublishedBySlug(ctx context.Context, slug string) (model.Article, error) {
	var article model.Article
	result := service.database.WithContext(ctx).
		Where("slug = ? AND status = ?", slug, model.ArticleStatusPublished).
		First(&article)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return model.Article{}, ErrArticleNotFound
	}
	if result.Error != nil {
		return model.Article{}, fmt.Errorf("find published article: %w", result.Error)
	}
	return article, nil
}

func (service *ArticleService) ListAll(ctx context.Context, page, pageSize int) (ArticlePage, error) {
	return service.ListAllByStatus(ctx, page, pageSize, "")
}

func (service *ArticleService) ListAllByStatus(
	ctx context.Context,
	page, pageSize int,
	status string,
) (ArticlePage, error) {
	query := service.database.WithContext(ctx).Model(&model.Article{})
	if status != "" {
		query = query.Where("status = ?", status)
	}
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return ArticlePage{}, fmt.Errorf("count articles: %w", err)
	}

	articles := make([]model.Article, 0)
	result := query.
		Select("id", "title", "slug", "summary", "status", "published_at", "created_at", "updated_at").
		Order("created_at DESC, id DESC").
		Limit(pageSize).
		Offset((page - 1) * pageSize).
		Find(&articles)
	if result.Error != nil {
		return ArticlePage{}, fmt.Errorf("list articles: %w", result.Error)
	}

	return newArticlePage(articles, page, pageSize, total), nil
}

func (service *ArticleService) ListPublished(ctx context.Context, page, pageSize int) (ArticlePage, error) {
	query := service.database.WithContext(ctx).
		Model(&model.Article{}).
		Where("status = ?", model.ArticleStatusPublished)

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return ArticlePage{}, fmt.Errorf("count published articles: %w", err)
	}

	articles := make([]model.Article, 0)
	result := query.
		Select(
			"articles.id, articles.title, articles.slug, articles.summary, articles.published_at, "+
				"(SELECT COUNT(*) FROM comments WHERE comments.article_id = articles.id "+
				"AND comments.status = ?) AS comment_count",
			model.CommentStatusApproved,
		).
		Order("published_at DESC, id DESC").
		Limit(pageSize).
		Offset((page - 1) * pageSize).
		Find(&articles)
	if result.Error != nil {
		return ArticlePage{}, fmt.Errorf("list published articles: %w", result.Error)
	}

	return newArticlePage(articles, page, pageSize, total), nil
}

func (service *ArticleService) slugExists(ctx context.Context, slug string, excludedID uint64) (bool, error) {
	query := service.database.WithContext(ctx).Model(&model.Article{}).Where("slug = ?", slug)
	if excludedID > 0 {
		query = query.Where("id <> ?", excludedID)
	}

	var count int64
	if err := query.Count(&count).Error; err != nil {
		return false, fmt.Errorf("check article slug: %w", err)
	}
	return count > 0, nil
}

func prepareArticleInput(input ArticleInput, create bool) (ArticleInput, error) {
	input.Title = strings.TrimSpace(input.Title)
	input.Slug = strings.TrimSpace(input.Slug)
	input.Summary = strings.TrimSpace(input.Summary)

	if input.Slug == "" && create {
		input.Slug = slugFromTitle(input.Title)
	}
	if input.Status == "" && create {
		input.Status = model.ArticleStatusDraft
	}

	if input.Title == "" || utf8.RuneCountInString(input.Title) > 200 {
		return ArticleInput{}, ErrInvalidArticle
	}
	if input.Slug == "" || utf8.RuneCountInString(input.Slug) > 200 || !slugPattern.MatchString(input.Slug) {
		return ArticleInput{}, ErrInvalidArticle
	}
	if utf8.RuneCountInString(input.Summary) > 500 || strings.TrimSpace(input.Content) == "" {
		return ArticleInput{}, ErrInvalidArticle
	}
	if input.Status != model.ArticleStatusDraft && input.Status != model.ArticleStatusPublished {
		return ArticleInput{}, ErrInvalidArticle
	}

	return input, nil
}

func slugFromTitle(title string) string {
	var builder strings.Builder
	lastWasHyphen := false
	for _, character := range strings.ToLower(title) {
		isLetter := character >= 'a' && character <= 'z'
		isDigit := character >= '0' && character <= '9'
		if isLetter || isDigit {
			builder.WriteRune(character)
			lastWasHyphen = false
			continue
		}
		if builder.Len() > 0 && !lastWasHyphen {
			builder.WriteByte('-')
			lastWasHyphen = true
		}
	}
	return strings.Trim(builder.String(), "-")
}

func newArticlePage(items []model.Article, page, pageSize int, total int64) ArticlePage {
	totalPages := 0
	if total > 0 {
		totalPages = int((total + int64(pageSize) - 1) / int64(pageSize))
	}
	return ArticlePage{
		Items:      items,
		Page:       page,
		PageSize:   pageSize,
		Total:      total,
		TotalPages: totalPages,
	}
}
