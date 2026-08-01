package model

import "time"

const (
	SocialPlatformGitHub   = "github"
	SocialPlatformEmail    = "email"
	SocialPlatformLinkedIn = "linkedin"
	SocialPlatformX        = "x"
	SocialPlatformWeibo    = "weibo"
	SocialPlatformBilibili = "bilibili"
	SocialPlatformZhihu    = "zhihu"
	SocialPlatformCustom   = "custom"
)

type SocialLink struct {
	ID          uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	Platform    string    `gorm:"type:varchar(30);not null" json:"platform"`
	DisplayName string    `gorm:"type:varchar(30);not null" json:"display_name"`
	URL         string    `gorm:"type:varchar(500);not null" json:"url"`
	IsVisible   bool      `gorm:"not null;default:true" json:"is_visible"`
	SortOrder   int       `gorm:"not null;index" json:"sort_order"`
	CreatedAt   time.Time `gorm:"not null" json:"created_at"`
	UpdatedAt   time.Time `gorm:"not null" json:"updated_at"`
}

func (SocialLink) TableName() string {
	return "social_links"
}
