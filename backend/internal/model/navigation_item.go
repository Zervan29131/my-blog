package model

import "time"

const (
	LinkTypeInternal = "internal"
	LinkTypeExternal = "external"
)

type NavigationItem struct {
	ID           uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	Name         string    `gorm:"type:varchar(20);not null" json:"name"`
	URL          string    `gorm:"type:varchar(500);not null" json:"url"`
	LinkType     string    `gorm:"type:varchar(20);not null;check:navigation_items_link_type_check,link_type IN ('internal','external')" json:"link_type"`
	OpenInNewTab bool      `gorm:"not null;default:false" json:"open_in_new_tab"`
	IsVisible    bool      `gorm:"not null;default:true" json:"is_visible"`
	SortOrder    int       `gorm:"not null;index" json:"sort_order"`
	CreatedAt    time.Time `gorm:"not null" json:"created_at"`
	UpdatedAt    time.Time `gorm:"not null" json:"updated_at"`
}

func (NavigationItem) TableName() string {
	return "navigation_items"
}
