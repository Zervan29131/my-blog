package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"personal-blog/backend/internal/model"
)

func Open(ctx context.Context, databaseURL string) (*gorm.DB, error) {
	database, err := gorm.Open(postgres.Open(databaseURL), &gorm.Config{
		DisableAutomaticPing: true,
		Logger:               logger.Default.LogMode(logger.Warn),
		TranslateError:       true,
	})
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}

	sqlDatabase, err := database.DB()
	if err != nil {
		return nil, fmt.Errorf("get database connection: %w", err)
	}
	sqlDatabase.SetMaxOpenConns(10)
	sqlDatabase.SetMaxIdleConns(5)
	sqlDatabase.SetConnMaxLifetime(time.Hour)

	if err := sqlDatabase.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("ping database: %w", err)
	}

	return database, nil
}

func Migrate(database *gorm.DB) error {
	if err := database.AutoMigrate(
		&model.Administrator{},
		&model.Article{},
		&model.Comment{},
		&model.SiteSetting{},
		&model.HomepageVersion{},
		&model.NavigationItem{},
		&model.SocialLink{},
		&model.FeaturedArticle{},
	); err != nil {
		return fmt.Errorf("migrate database: %w", err)
	}
	return nil
}

func InitializeHomepageCMS(database *gorm.DB) error {
	now := time.Now().UTC()
	if err := database.Transaction(func(transaction *gorm.DB) error {
		if err := initializeSiteSetting(transaction, now); err != nil {
			return err
		}
		if err := initializeNavigation(transaction); err != nil {
			return err
		}
		if err := initializeSocialLinks(transaction); err != nil {
			return err
		}
		return initializeHomepageVersions(transaction, now)
	}); err != nil {
		return fmt.Errorf("initialize homepage CMS: %w", err)
	}
	return nil
}

func initializeSiteSetting(database *gorm.DB, now time.Time) error {
	var count int64
	if err := database.Model(&model.SiteSetting{}).Count(&count).Error; err != nil {
		return fmt.Errorf("count site settings: %w", err)
	}
	if count > 0 {
		return nil
	}
	setting := model.DefaultSiteSetting(now)
	if err := model.ValidateSiteSetting(setting); err != nil {
		return fmt.Errorf("validate default site setting: %w", err)
	}
	if err := database.Create(&setting).Error; err != nil {
		return fmt.Errorf("create default site setting: %w", err)
	}
	return nil
}

func initializeNavigation(database *gorm.DB) error {
	var count int64
	if err := database.Model(&model.NavigationItem{}).Count(&count).Error; err != nil {
		return fmt.Errorf("count navigation items: %w", err)
	}
	if count > 0 {
		return nil
	}
	items := model.DefaultNavigationItems()
	if err := model.ValidateNavigationItems(items); err != nil {
		return fmt.Errorf("validate default navigation: %w", err)
	}
	if err := database.Create(&items).Error; err != nil {
		return fmt.Errorf("create default navigation: %w", err)
	}
	return nil
}

func initializeSocialLinks(database *gorm.DB) error {
	var count int64
	if err := database.Model(&model.SocialLink{}).Count(&count).Error; err != nil {
		return fmt.Errorf("count social links: %w", err)
	}
	if count > 0 {
		return nil
	}
	links := model.DefaultSocialLinks()
	if err := model.ValidateSocialLinks(links); err != nil {
		return fmt.Errorf("validate default social links: %w", err)
	}
	if err := database.Create(&links).Error; err != nil {
		return fmt.Errorf("create default social links: %w", err)
	}
	return nil
}

func initializeHomepageVersions(database *gorm.DB, now time.Time) error {
	var published model.HomepageVersion
	result := database.Where("status = ?", model.HomepageStatusPublished).Limit(1).Find(&published)
	if result.Error != nil {
		return fmt.Errorf("find published homepage version: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		configJSON, err := model.EncodeHomepageConfig(model.DefaultHomepageConfig())
		if err != nil {
			return fmt.Errorf("encode default homepage config: %w", err)
		}
		published = model.HomepageVersion{
			Status: model.HomepageStatusPublished, VersionNumber: 1,
			ConfigJSON: configJSON, PublishedAt: &now,
		}
		if err := database.Create(&published).Error; err != nil {
			return fmt.Errorf("create published homepage version: %w", err)
		}
	}

	var draft model.HomepageVersion
	result = database.Where("status = ?", model.HomepageStatusDraft).Limit(1).Find(&draft)
	if result.Error != nil {
		return fmt.Errorf("find draft homepage version: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		draft = model.HomepageVersion{
			Status: model.HomepageStatusDraft, VersionNumber: published.VersionNumber,
			ConfigJSON: append(model.JSONDocument(nil), published.ConfigJSON...),
		}
		if err := database.Create(&draft).Error; err != nil {
			return fmt.Errorf("create draft homepage version: %w", err)
		}
	}
	return nil
}

func SQL(database *gorm.DB) (*sql.DB, error) {
	sqlDatabase, err := database.DB()
	if err != nil {
		return nil, fmt.Errorf("get database connection: %w", err)
	}
	return sqlDatabase, nil
}
