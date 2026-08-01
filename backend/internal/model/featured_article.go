package model

import "time"

type FeaturedArticle struct {
	ID        uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	ArticleID uint64    `gorm:"not null;uniqueIndex" json:"article_id"`
	Article   Article   `gorm:"foreignKey:ArticleID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"-"`
	SortOrder int       `gorm:"not null;index" json:"sort_order"`
	IsVisible bool      `gorm:"not null;default:true" json:"is_visible"`
	CreatedAt time.Time `gorm:"not null" json:"created_at"`
	UpdatedAt time.Time `gorm:"not null" json:"updated_at"`
}

func (FeaturedArticle) TableName() string {
	return "featured_articles"
}
