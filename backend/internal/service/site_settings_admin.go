package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"personal-blog/backend/internal/model"
)

var (
	ErrSiteSettingNotFound = errors.New("site setting not found")
	ErrInvalidSiteSetting  = errors.New("invalid site setting")
)

type SiteSettingsAdminService struct {
	database *gorm.DB
}

type SiteSettingInput struct {
	SiteName             string
	SiteShortName        *string
	SiteDescription      string
	TitleSuffix          *string
	LogoURL              *string
	FaviconURL           *string
	DefaultShareImageURL *string
	CopyrightName        string
	StartYear            *int
	AdditionalText       *string
	FilingNumber         *string
	FilingURL            *string
	ShowTechnology       bool
	TechnologyText       *string
}

func NewSiteSettingsAdminService(database *gorm.DB) *SiteSettingsAdminService {
	return &SiteSettingsAdminService{database: database}
}

func (service *SiteSettingsAdminService) Get(ctx context.Context) (model.SiteSetting, error) {
	var setting model.SiteSetting
	result := service.database.WithContext(ctx).Where("id = ?", uint64(1)).Limit(1).Find(&setting)
	if result.Error != nil {
		return model.SiteSetting{}, fmt.Errorf("find administrator site setting: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return model.SiteSetting{}, ErrSiteSettingNotFound
	}
	return setting, nil
}

func (service *SiteSettingsAdminService) Update(
	ctx context.Context,
	input SiteSettingInput,
	administratorID uint64,
) (model.SiteSetting, error) {
	prepared, err := siteSettingFromInput(input)
	if err != nil {
		return model.SiteSetting{}, err
	}

	now := time.Now().UTC()
	var updated model.SiteSetting
	err = service.database.WithContext(ctx).Transaction(func(transaction *gorm.DB) error {
		result := transaction.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("id = ?", uint64(1)).
			Limit(1).
			Find(&updated)
		if result.Error != nil {
			return fmt.Errorf("lock site setting: %w", result.Error)
		}
		if result.RowsAffected == 0 {
			return ErrSiteSettingNotFound
		}

		result = transaction.Model(&model.SiteSetting{}).
			Where("id = ?", updated.ID).
			Updates(map[string]any{
				"site_name": prepared.SiteName, "site_short_name": prepared.SiteShortName,
				"site_description": prepared.SiteDescription, "title_suffix": prepared.TitleSuffix,
				"logo_url": prepared.LogoURL, "favicon_url": prepared.FaviconURL,
				"default_share_image_url": prepared.DefaultShareImageURL,
				"copyright_name":          prepared.CopyrightName, "start_year": prepared.StartYear,
				"additional_text": prepared.AdditionalText, "filing_number": prepared.FilingNumber,
				"filing_url": prepared.FilingURL, "show_technology": prepared.ShowTechnology,
				"technology_text": prepared.TechnologyText, "updated_at": now,
			})
		if result.Error != nil {
			return fmt.Errorf("update site setting: %w", result.Error)
		}
		if result.RowsAffected == 0 {
			return ErrSiteSettingNotFound
		}

		createdAt := updated.CreatedAt
		updated = prepared
		updated.ID = 1
		updated.CreatedAt = createdAt
		updated.UpdatedAt = now
		return nil
	})
	if err != nil {
		return model.SiteSetting{}, err
	}

	slog.Info(
		"site setting updated",
		"operation", "site_setting_update",
		"administrator_id", administratorID,
		"target_type", "site_setting",
		"target_id", updated.ID,
		"operation_time", now,
	)
	return updated, nil
}

func siteSettingFromInput(input SiteSettingInput) (model.SiteSetting, error) {
	setting := model.SiteSetting{
		ID: 1, SiteName: input.SiteName, SiteShortName: emptyStringAsNil(input.SiteShortName),
		SiteDescription: input.SiteDescription, TitleSuffix: emptyStringAsNil(input.TitleSuffix),
		LogoURL: emptyStringAsNil(input.LogoURL), FaviconURL: emptyStringAsNil(input.FaviconURL),
		DefaultShareImageURL: emptyStringAsNil(input.DefaultShareImageURL),
		CopyrightName:        input.CopyrightName, StartYear: input.StartYear,
		AdditionalText: emptyStringAsNil(input.AdditionalText), FilingNumber: emptyStringAsNil(input.FilingNumber),
		FilingURL: emptyStringAsNil(input.FilingURL), ShowTechnology: input.ShowTechnology,
		TechnologyText: emptyStringAsNil(input.TechnologyText),
	}
	if err := model.ValidateSiteSetting(setting); err != nil {
		return model.SiteSetting{}, fmt.Errorf("%w: %v", ErrInvalidSiteSetting, err)
	}
	return setting, nil
}

func emptyStringAsNil(value *string) *string {
	if value == nil || *value == "" {
		return nil
	}
	copy := *value
	return &copy
}
