package model

import "time"

type SiteSetting struct {
	ID                   uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	SiteName             string    `gorm:"type:varchar(50);not null" json:"site_name"`
	SiteShortName        *string   `gorm:"type:varchar(20)" json:"site_short_name"`
	SiteDescription      string    `gorm:"type:varchar(200);not null" json:"site_description"`
	TitleSuffix          *string   `gorm:"type:varchar(50)" json:"title_suffix"`
	LogoURL              *string   `gorm:"type:varchar(500)" json:"logo_url"`
	FaviconURL           *string   `gorm:"type:varchar(500)" json:"favicon_url"`
	DefaultShareImageURL *string   `gorm:"type:varchar(500)" json:"default_share_image_url"`
	CopyrightName        string    `gorm:"type:varchar(50);not null" json:"copyright_name"`
	StartYear            *int      `json:"start_year"`
	AdditionalText       *string   `gorm:"type:varchar(200)" json:"additional_text"`
	FilingNumber         *string   `gorm:"type:varchar(100)" json:"filing_number"`
	FilingURL            *string   `gorm:"type:varchar(500)" json:"filing_url"`
	ShowTechnology       bool      `gorm:"not null;default:true" json:"show_technology"`
	TechnologyText       *string   `gorm:"type:varchar(100)" json:"technology_text"`
	CreatedAt            time.Time `gorm:"not null" json:"created_at"`
	UpdatedAt            time.Time `gorm:"not null" json:"updated_at"`
}

func (SiteSetting) TableName() string {
	return "site_settings"
}
