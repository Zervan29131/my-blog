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
	ErrHomepageVersionNotFound = errors.New("homepage version not found")
	ErrInvalidHomepageConfig   = errors.New("invalid homepage config")
)

type HomepageAdminService struct {
	database            *gorm.DB
	publicConfigService *PublicConfigService
}

type AdminHomepageConfig struct {
	Status      string                 `json:"status"`
	Version     uint64                 `json:"version"`
	Modules     []model.HomepageModule `json:"modules"`
	UpdatedAt   time.Time              `json:"updated_at"`
	PublishedAt *time.Time             `json:"published_at"`
}

type HomepagePublishResult struct {
	Version     uint64    `json:"version"`
	PublishedAt time.Time `json:"published_at"`
}

func NewHomepageAdminService(database *gorm.DB) *HomepageAdminService {
	return &HomepageAdminService{
		database:            database,
		publicConfigService: NewPublicConfigService(database),
	}
}

func (service *HomepageAdminService) GetDraft(ctx context.Context) (AdminHomepageConfig, error) {
	return service.getConfig(ctx, model.HomepageStatusDraft)
}

func (service *HomepageAdminService) GetPublished(ctx context.Context) (AdminHomepageConfig, error) {
	return service.getConfig(ctx, model.HomepageStatusPublished)
}

func (service *HomepageAdminService) SaveDraft(
	ctx context.Context,
	config model.HomepageConfig,
	administratorID uint64,
) (AdminHomepageConfig, error) {
	configJSON, err := encodeAdminHomepageConfig(config)
	if err != nil {
		return AdminHomepageConfig{}, err
	}

	now := time.Now().UTC()
	var saved AdminHomepageConfig
	var targetID uint64
	err = service.database.WithContext(ctx).Transaction(func(transaction *gorm.DB) error {
		draft, err := lockedHomepageVersion(transaction, model.HomepageStatusDraft)
		if err != nil {
			return err
		}
		result := transaction.Model(&model.HomepageVersion{}).
			Where("id = ?", draft.ID).
			Updates(map[string]any{"config_json": configJSON, "updated_at": now})
		if result.Error != nil {
			return fmt.Errorf("save homepage draft: %w", result.Error)
		}
		if result.RowsAffected == 0 {
			return ErrHomepageVersionNotFound
		}
		draft.ConfigJSON = configJSON
		draft.UpdatedAt = now
		targetID = draft.ID
		saved = adminHomepageConfig(draft, config)
		return nil
	})
	if err != nil {
		return AdminHomepageConfig{}, err
	}
	slog.Info(
		"homepage draft saved",
		"operation", "homepage_draft_save",
		"administrator_id", administratorID,
		"target_type", "homepage_version",
		"target_id", targetID,
		"version", saved.Version,
		"operation_time", now,
	)
	return saved, nil
}

func (service *HomepageAdminService) Publish(
	ctx context.Context,
	administratorID uint64,
) (HomepagePublishResult, error) {
	var publishedResult HomepagePublishResult
	var targetID uint64
	err := service.database.WithContext(ctx).Transaction(func(transaction *gorm.DB) error {
		draft, err := lockedHomepageVersion(transaction, model.HomepageStatusDraft)
		if err != nil {
			return err
		}
		published, err := lockedHomepageVersion(transaction, model.HomepageStatusPublished)
		if err != nil {
			return err
		}

		config, err := decodeStoredHomepageConfig(draft)
		if err != nil {
			return err
		}
		configJSON, err := encodeAdminHomepageConfig(config)
		if err != nil {
			return err
		}

		now := time.Now().UTC()
		newVersion := published.VersionNumber + 1
		result := transaction.Model(&model.HomepageVersion{}).
			Where("id = ?", published.ID).
			Updates(map[string]any{
				"config_json":    configJSON,
				"version_number": newVersion,
				"published_at":   now,
				"updated_at":     now,
			})
		if result.Error != nil {
			return fmt.Errorf("publish homepage config: %w", result.Error)
		}
		if result.RowsAffected == 0 {
			return ErrHomepageVersionNotFound
		}

		result = transaction.Model(&model.HomepageVersion{}).
			Where("id = ?", draft.ID).
			Updates(map[string]any{"version_number": newVersion, "updated_at": now})
		if result.Error != nil {
			return fmt.Errorf("synchronize homepage draft version: %w", result.Error)
		}
		if result.RowsAffected == 0 {
			return ErrHomepageVersionNotFound
		}

		publishedResult = HomepagePublishResult{Version: newVersion, PublishedAt: now}
		targetID = published.ID
		return nil
	})
	if err != nil {
		return HomepagePublishResult{}, err
	}

	slog.Info(
		"homepage config published",
		"operation", "homepage_publish",
		"administrator_id", administratorID,
		"target_type", "homepage_version",
		"target_id", targetID,
		"version", publishedResult.Version,
		"operation_time", publishedResult.PublishedAt,
	)
	return publishedResult, nil
}

func (service *HomepageAdminService) ResetDraft(
	ctx context.Context,
	administratorID uint64,
) (AdminHomepageConfig, error) {
	var resetDraft AdminHomepageConfig
	var targetID uint64
	now := time.Now().UTC()
	err := service.database.WithContext(ctx).Transaction(func(transaction *gorm.DB) error {
		draft, err := lockedHomepageVersion(transaction, model.HomepageStatusDraft)
		if err != nil {
			return err
		}
		published, err := lockedHomepageVersion(transaction, model.HomepageStatusPublished)
		if err != nil {
			return err
		}
		config, err := decodeStoredHomepageConfig(published)
		if err != nil {
			return err
		}
		configJSON, err := encodeAdminHomepageConfig(config)
		if err != nil {
			return err
		}

		result := transaction.Model(&model.HomepageVersion{}).
			Where("id = ?", draft.ID).
			Updates(map[string]any{
				"config_json":    configJSON,
				"version_number": published.VersionNumber,
				"updated_at":     now,
			})
		if result.Error != nil {
			return fmt.Errorf("reset homepage draft: %w", result.Error)
		}
		if result.RowsAffected == 0 {
			return ErrHomepageVersionNotFound
		}

		resetDraft = adminHomepageConfig(model.HomepageVersion{
			Status:        model.HomepageStatusDraft,
			VersionNumber: published.VersionNumber,
			ConfigJSON:    configJSON,
			UpdatedAt:     now,
		}, config)
		targetID = draft.ID
		return nil
	})
	if err != nil {
		return AdminHomepageConfig{}, err
	}

	slog.Info(
		"homepage draft reset",
		"operation", "homepage_draft_reset",
		"administrator_id", administratorID,
		"target_type", "homepage_version",
		"target_id", targetID,
		"version", resetDraft.Version,
		"operation_time", now,
	)
	return resetDraft, nil
}

func (service *HomepageAdminService) PreviewDraft(ctx context.Context) (PublicHomepage, error) {
	draft, err := service.homepageVersion(ctx, model.HomepageStatusDraft)
	if err != nil {
		return PublicHomepage{}, err
	}
	config, err := decodeStoredHomepageConfig(draft)
	if err != nil {
		return PublicHomepage{}, err
	}
	return service.publicConfigService.composeHomepage(ctx, draft.VersionNumber, config)
}

func (service *HomepageAdminService) getConfig(
	ctx context.Context,
	status string,
) (AdminHomepageConfig, error) {
	version, err := service.homepageVersion(ctx, status)
	if err != nil {
		return AdminHomepageConfig{}, err
	}
	config, err := decodeStoredHomepageConfig(version)
	if err != nil {
		return AdminHomepageConfig{}, err
	}
	return adminHomepageConfig(version, config), nil
}

func (service *HomepageAdminService) homepageVersion(
	ctx context.Context,
	status string,
) (model.HomepageVersion, error) {
	var version model.HomepageVersion
	result := service.database.WithContext(ctx).Where("status = ?", status).Limit(1).Find(&version)
	if result.Error != nil {
		return model.HomepageVersion{}, fmt.Errorf("find %s homepage version: %w", status, result.Error)
	}
	if result.RowsAffected == 0 {
		return model.HomepageVersion{}, ErrHomepageVersionNotFound
	}
	return version, nil
}

func lockedHomepageVersion(database *gorm.DB, status string) (model.HomepageVersion, error) {
	var version model.HomepageVersion
	result := database.
		Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("status = ?", status).
		Limit(1).
		Find(&version)
	if result.Error != nil {
		return model.HomepageVersion{}, fmt.Errorf("lock %s homepage version: %w", status, result.Error)
	}
	if result.RowsAffected == 0 {
		return model.HomepageVersion{}, ErrHomepageVersionNotFound
	}
	return version, nil
}

func decodeStoredHomepageConfig(version model.HomepageVersion) (model.HomepageConfig, error) {
	config, err := model.DecodeHomepageConfig(version.ConfigJSON)
	if err != nil {
		return model.HomepageConfig{}, fmt.Errorf("decode stored %s homepage config: %w", version.Status, err)
	}
	return config, nil
}

func encodeAdminHomepageConfig(config model.HomepageConfig) (model.JSONDocument, error) {
	configJSON, err := model.EncodeHomepageConfig(config)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidHomepageConfig, err)
	}
	return configJSON, nil
}

func adminHomepageConfig(
	version model.HomepageVersion,
	config model.HomepageConfig,
) AdminHomepageConfig {
	return AdminHomepageConfig{
		Status: version.Status, Version: version.VersionNumber, Modules: config.Modules,
		UpdatedAt: version.UpdatedAt, PublishedAt: version.PublishedAt,
	}
}
