package service

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"gorm.io/gorm"

	"personal-blog/backend/internal/model"
)

const maximumNavigationItems = 10

type NavigationAdminService struct {
	database *gorm.DB
}

type NavigationInput struct {
	Name         string
	URL          string
	LinkType     string
	OpenInNewTab bool
	IsVisible    bool
	SortOrder    int
}

func NewNavigationAdminService(database *gorm.DB) *NavigationAdminService {
	return &NavigationAdminService{database: database}
}

func (service *NavigationAdminService) List(ctx context.Context) ([]model.NavigationItem, error) {
	items := make([]model.NavigationItem, 0)
	if err := service.database.WithContext(ctx).
		Order("sort_order ASC, id ASC").
		Find(&items).Error; err != nil {
		return nil, fmt.Errorf("list administrator navigation: %w", err)
	}
	return items, nil
}

func (service *NavigationAdminService) Create(
	ctx context.Context,
	input NavigationInput,
	administratorID uint64,
) (model.NavigationItem, error) {
	item, err := navigationFromInput(input)
	if err != nil {
		return model.NavigationItem{}, err
	}

	err = service.database.WithContext(ctx).Transaction(func(transaction *gorm.DB) error {
		var count int64
		if err := transaction.Model(&model.NavigationItem{}).Count(&count).Error; err != nil {
			return fmt.Errorf("count navigation items: %w", err)
		}
		if count >= maximumNavigationItems {
			return ErrSiteContentLimit
		}
		if err := transaction.Create(&item).Error; err != nil {
			return fmt.Errorf("create navigation item: %w", err)
		}
		return nil
	})
	if err != nil {
		return model.NavigationItem{}, err
	}
	logSiteContentOperation("navigation_created", administratorID, "navigation_item", item.ID)
	return item, nil
}

func (service *NavigationAdminService) Update(
	ctx context.Context,
	id uint64,
	input NavigationInput,
	administratorID uint64,
) (model.NavigationItem, error) {
	item, err := navigationFromInput(input)
	if err != nil {
		return model.NavigationItem{}, err
	}
	var existing model.NavigationItem
	findResult := service.database.WithContext(ctx).Where("id = ?", id).Limit(1).Find(&existing)
	if findResult.Error != nil {
		return model.NavigationItem{}, fmt.Errorf("find navigation item: %w", findResult.Error)
	}
	if findResult.RowsAffected == 0 {
		return model.NavigationItem{}, ErrSiteContentNotFound
	}
	item.ID = id
	item.CreatedAt = existing.CreatedAt
	item.UpdatedAt = time.Now().UTC()
	result := service.database.WithContext(ctx).
		Model(&model.NavigationItem{}).
		Where("id = ?", id).
		Updates(map[string]any{
			"name": item.Name, "url": item.URL, "link_type": item.LinkType,
			"open_in_new_tab": item.OpenInNewTab, "is_visible": item.IsVisible,
			"sort_order": item.SortOrder, "updated_at": item.UpdatedAt,
		})
	if result.Error != nil {
		return model.NavigationItem{}, fmt.Errorf("update navigation item: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return model.NavigationItem{}, ErrSiteContentNotFound
	}
	logSiteContentOperation("navigation_updated", administratorID, "navigation_item", id)
	return item, nil
}

func (service *NavigationAdminService) Delete(
	ctx context.Context,
	id uint64,
	administratorID uint64,
) error {
	result := service.database.WithContext(ctx).Delete(&model.NavigationItem{}, id)
	if result.Error != nil {
		return fmt.Errorf("delete navigation item: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return ErrSiteContentNotFound
	}
	logSiteContentOperation("navigation_deleted", administratorID, "navigation_item", id)
	return nil
}

func (service *NavigationAdminService) Reorder(
	ctx context.Context,
	items []ResourceOrder,
	administratorID uint64,
) error {
	if err := validateResourceOrder(items, maximumNavigationItems); err != nil {
		return err
	}
	now := time.Now().UTC()
	err := service.database.WithContext(ctx).Transaction(func(transaction *gorm.DB) error {
		for _, item := range items {
			result := transaction.Model(&model.NavigationItem{}).
				Where("id = ?", item.ID).
				Updates(map[string]any{"sort_order": item.SortOrder, "updated_at": now})
			if result.Error != nil {
				return fmt.Errorf("reorder navigation item: %w", result.Error)
			}
			if result.RowsAffected == 0 {
				return ErrSiteContentNotFound
			}
		}
		return nil
	})
	if err != nil {
		return err
	}
	logSiteContentOperation("navigation_reordered", administratorID, "navigation", 0)
	return nil
}

func navigationFromInput(input NavigationInput) (model.NavigationItem, error) {
	item := model.NavigationItem{
		Name: input.Name, URL: input.URL, LinkType: input.LinkType,
		OpenInNewTab: input.OpenInNewTab, IsVisible: input.IsVisible, SortOrder: input.SortOrder,
	}
	if err := model.ValidateNavigationItems([]model.NavigationItem{item}); err != nil {
		return model.NavigationItem{}, fmt.Errorf("%w: %v", ErrInvalidSiteContent, err)
	}
	return item, nil
}

func logSiteContentOperation(operation string, administratorID uint64, targetType string, targetID uint64) {
	slog.Info(
		"site content changed",
		"operation", operation,
		"administrator_id", administratorID,
		"target_type", targetType,
		"target_id", targetID,
		"operation_time", time.Now().UTC(),
	)
}
