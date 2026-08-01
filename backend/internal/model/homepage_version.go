package model

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"
)

const (
	HomepageStatusDraft     = "draft"
	HomepageStatusPublished = "published"
)

type JSONDocument json.RawMessage

func (document JSONDocument) Value() (driver.Value, error) {
	if len(document) == 0 {
		return nil, nil
	}
	if !json.Valid(document) {
		return nil, fmt.Errorf("invalid JSON document")
	}
	return string(document), nil
}

func (document *JSONDocument) Scan(value any) error {
	if value == nil {
		*document = nil
		return nil
	}

	var raw []byte
	switch typed := value.(type) {
	case []byte:
		raw = typed
	case string:
		raw = []byte(typed)
	default:
		return fmt.Errorf("scan JSON document from %T", value)
	}
	if !json.Valid(raw) {
		return fmt.Errorf("scan invalid JSON document")
	}
	*document = append((*document)[:0], raw...)
	return nil
}

type HomepageVersion struct {
	ID            uint64       `gorm:"primaryKey;autoIncrement" json:"id"`
	Status        string       `gorm:"type:varchar(20);uniqueIndex;not null;check:homepage_versions_status_check,status IN ('draft','published')" json:"status"`
	VersionNumber uint64       `gorm:"not null" json:"version_number"`
	ConfigJSON    JSONDocument `gorm:"type:jsonb;not null" json:"-"`
	PublishedAt   *time.Time   `json:"published_at"`
	CreatedAt     time.Time    `gorm:"not null" json:"created_at"`
	UpdatedAt     time.Time    `gorm:"not null" json:"updated_at"`
}

func (HomepageVersion) TableName() string {
	return "homepage_versions"
}
