package service

import (
	"context"
	"errors"
	"fmt"
	"net/mail"
	"strings"
	"time"
	"unicode/utf8"

	"gorm.io/gorm"

	"personal-blog/backend/internal/model"
)

var (
	ErrCommentNotFound = errors.New("comment not found")
	ErrInvalidComment  = errors.New("invalid comment")
)

type CommentService struct {
	database *gorm.DB
}

type CommentInput struct {
	Nickname string
	Email    string
	Content  string
}

type PublicComment struct {
	ID        uint64    `gorm:"column:id"`
	Nickname  string    `gorm:"column:nickname"`
	Content   string    `gorm:"column:content"`
	CreatedAt time.Time `gorm:"column:created_at"`
}

type AdminComment struct {
	ID           uint64    `gorm:"column:id"`
	ArticleID    uint64    `gorm:"column:article_id"`
	ArticleTitle string    `gorm:"column:article_title"`
	ArticleSlug  string    `gorm:"column:article_slug"`
	Nickname     string    `gorm:"column:nickname"`
	Email        *string   `gorm:"column:email"`
	Content      string    `gorm:"column:content"`
	Status       string    `gorm:"column:status"`
	CreatedAt    time.Time `gorm:"column:created_at"`
	UpdatedAt    time.Time `gorm:"column:updated_at"`
}

type PublicCommentPage struct {
	Items      []PublicComment
	Page       int
	PageSize   int
	Total      int64
	TotalPages int
}

type AdminCommentPage struct {
	Items      []AdminComment
	Page       int
	PageSize   int
	Total      int64
	TotalPages int
}

func NewCommentService(database *gorm.DB) *CommentService {
	return &CommentService{database: database}
}

func (service *CommentService) Submit(
	ctx context.Context,
	articleSlug string,
	input CommentInput,
) (model.Comment, error) {
	prepared, email, err := prepareCommentInput(input)
	if err != nil {
		return model.Comment{}, err
	}

	article, err := service.findPublishedArticle(ctx, articleSlug)
	if err != nil {
		return model.Comment{}, err
	}

	comment := model.Comment{
		ArticleID: article.ID,
		Nickname:  prepared.Nickname,
		Email:     email,
		Content:   prepared.Content,
		Status:    model.CommentStatusPending,
	}
	if err := service.database.WithContext(ctx).Create(&comment).Error; err != nil {
		return model.Comment{}, fmt.Errorf("create comment: %w", err)
	}
	return comment, nil
}

func (service *CommentService) ListApproved(
	ctx context.Context,
	articleSlug string,
	page, pageSize int,
) (PublicCommentPage, error) {
	article, err := service.findPublishedArticle(ctx, articleSlug)
	if err != nil {
		return PublicCommentPage{}, err
	}

	query := service.database.WithContext(ctx).
		Model(&model.Comment{}).
		Where("article_id = ? AND status = ?", article.ID, model.CommentStatusApproved)

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return PublicCommentPage{}, fmt.Errorf("count approved comments: %w", err)
	}

	comments := make([]PublicComment, 0)
	if err := query.
		Select("id", "nickname", "content", "created_at").
		Order("created_at ASC, id ASC").
		Limit(pageSize).
		Offset((page - 1) * pageSize).
		Scan(&comments).Error; err != nil {
		return PublicCommentPage{}, fmt.Errorf("list approved comments: %w", err)
	}

	return PublicCommentPage{
		Items:      comments,
		Page:       page,
		PageSize:   pageSize,
		Total:      total,
		TotalPages: totalPages(total, pageSize),
	}, nil
}

func (service *CommentService) ListAdmin(
	ctx context.Context,
	status string,
	page, pageSize int,
) (AdminCommentPage, error) {
	if status != "" && !IsValidCommentStatus(status) {
		return AdminCommentPage{}, ErrInvalidComment
	}

	query := service.database.WithContext(ctx).Model(&model.Comment{})
	if status != "" {
		query = query.Where("comments.status = ?", status)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return AdminCommentPage{}, fmt.Errorf("count administrator comments: %w", err)
	}

	comments := make([]AdminComment, 0)
	if err := query.
		Select(
			"comments.id, comments.article_id, articles.title AS article_title, " +
				"articles.slug AS article_slug, comments.nickname, comments.email, " +
				"comments.content, comments.status, comments.created_at, comments.updated_at",
		).
		Joins("JOIN articles ON articles.id = comments.article_id").
		Order("comments.created_at DESC, comments.id DESC").
		Limit(pageSize).
		Offset((page - 1) * pageSize).
		Scan(&comments).Error; err != nil {
		return AdminCommentPage{}, fmt.Errorf("list administrator comments: %w", err)
	}

	return AdminCommentPage{
		Items:      comments,
		Page:       page,
		PageSize:   pageSize,
		Total:      total,
		TotalPages: totalPages(total, pageSize),
	}, nil
}

func (service *CommentService) UpdateStatus(ctx context.Context, id uint64, status string) error {
	if !IsValidCommentStatus(status) {
		return ErrInvalidComment
	}

	result := service.database.WithContext(ctx).
		Model(&model.Comment{}).
		Where("id = ?", id).
		Updates(map[string]any{
			"status":     status,
			"updated_at": time.Now().UTC(),
		})
	if result.Error != nil {
		return fmt.Errorf("update comment status: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return ErrCommentNotFound
	}
	return nil
}

func (service *CommentService) Delete(ctx context.Context, id uint64) error {
	result := service.database.WithContext(ctx).Delete(&model.Comment{}, id)
	if result.Error != nil {
		return fmt.Errorf("delete comment: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return ErrCommentNotFound
	}
	return nil
}

func (service *CommentService) findPublishedArticle(ctx context.Context, slug string) (model.Article, error) {
	var article model.Article
	result := service.database.WithContext(ctx).
		Select("id", "title", "slug").
		Where("slug = ? AND status = ?", slug, model.ArticleStatusPublished).
		First(&article)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return model.Article{}, ErrArticleNotFound
	}
	if result.Error != nil {
		return model.Article{}, fmt.Errorf("find published article for comments: %w", result.Error)
	}
	return article, nil
}

func prepareCommentInput(input CommentInput) (CommentInput, *string, error) {
	input.Nickname = strings.TrimSpace(input.Nickname)
	input.Email = strings.TrimSpace(input.Email)
	input.Content = strings.TrimSpace(input.Content)

	nicknameLength := utf8.RuneCountInString(input.Nickname)
	contentLength := utf8.RuneCountInString(input.Content)
	if nicknameLength < 2 || nicknameLength > 50 || contentLength < 2 || contentLength > 1000 {
		return CommentInput{}, nil, ErrInvalidComment
	}

	var email *string
	if input.Email != "" {
		if utf8.RuneCountInString(input.Email) > 255 {
			return CommentInput{}, nil, ErrInvalidComment
		}
		parsed, err := mail.ParseAddress(input.Email)
		if err != nil || parsed.Name != "" || !strings.EqualFold(parsed.Address, input.Email) {
			return CommentInput{}, nil, ErrInvalidComment
		}
		email = &input.Email
	}

	return input, email, nil
}

func IsValidCommentStatus(status string) bool {
	return status == model.CommentStatusPending ||
		status == model.CommentStatusApproved ||
		status == model.CommentStatusRejected
}

func totalPages(total int64, pageSize int) int {
	if total == 0 {
		return 0
	}
	return int((total + int64(pageSize) - 1) / int64(pageSize))
}
