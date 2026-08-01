package model

import "time"

const (
	ArticleStatusDraft     = "draft"
	ArticleStatusPublished = "published"
)

type Article struct {
	ID           uint64     `gorm:"primaryKey;autoIncrement" json:"id"`
	Title        string     `gorm:"type:varchar(200);not null" json:"title"`
	Slug         string     `gorm:"type:varchar(200);uniqueIndex;not null" json:"slug"`
	Summary      string     `gorm:"type:varchar(500);not null;default:''" json:"summary"`
	Content      string     `gorm:"type:text;not null" json:"content"`
	Status       string     `gorm:"type:varchar(20);not null;default:draft;check:articles_status_check,status IN ('draft','published')" json:"status"`
	PublishedAt  *time.Time `json:"published_at"`
	CreatedAt    time.Time  `gorm:"not null" json:"created_at"`
	UpdatedAt    time.Time  `gorm:"not null" json:"updated_at"`
	CommentCount int64      `gorm:"->;-:migration" json:"comment_count"`
}

func (Article) TableName() string {
	return "articles"
}
