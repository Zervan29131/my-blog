package service

import (
	"context"
	"fmt"
	"time"

	"gorm.io/gorm"

	"personal-blog/backend/internal/model"
)

const maximumSocialLinks = 15

type SocialLinkAdminService struct {
	database *gorm.DB
}

type SocialLinkInput struct {
	Platform    string
	DisplayName string
	URL         string
	IsVisible   bool
	SortOrder   int
}

func NewSocialLinkAdminService(database *gorm.DB) *SocialLinkAdminService {
	return &SocialLinkAdminService{database: database}
}

func (service *SocialLinkAdminService) List(ctx context.Context) ([]model.SocialLink, error) {
	links := make([]model.SocialLink, 0)
	if err := service.database.WithContext(ctx).
		Order("sort_order ASC, id ASC").
		Find(&links).Error; err != nil {
		return nil, fmt.Errorf("list administrator social links: %w", err)
	}
	return links, nil
}

func (service *SocialLinkAdminService) Create(
	ctx context.Context,
	input SocialLinkInput,
	administratorID uint64,
) (model.SocialLink, error) {
	link, err := socialLinkFromInput(input)
	if err != nil {
		return model.SocialLink{}, err
	}
	err = service.database.WithContext(ctx).Transaction(func(transaction *gorm.DB) error {
		var count int64
		if err := transaction.Model(&model.SocialLink{}).Count(&count).Error; err != nil {
			return fmt.Errorf("count social links: %w", err)
		}
		if count >= maximumSocialLinks {
			return ErrSiteContentLimit
		}
		if err := transaction.Create(&link).Error; err != nil {
			return fmt.Errorf("create social link: %w", err)
		}
		return nil
	})
	if err != nil {
		return model.SocialLink{}, err
	}
	logSiteContentOperation("social_link_created", administratorID, "social_link", link.ID)
	return link, nil
}

func (service *SocialLinkAdminService) Update(
	ctx context.Context,
	id uint64,
	input SocialLinkInput,
	administratorID uint64,
) (model.SocialLink, error) {
	link, err := socialLinkFromInput(input)
	if err != nil {
		return model.SocialLink{}, err
	}
	var existing model.SocialLink
	findResult := service.database.WithContext(ctx).Where("id = ?", id).Limit(1).Find(&existing)
	if findResult.Error != nil {
		return model.SocialLink{}, fmt.Errorf("find social link: %w", findResult.Error)
	}
	if findResult.RowsAffected == 0 {
		return model.SocialLink{}, ErrSiteContentNotFound
	}
	link.ID = id
	link.CreatedAt = existing.CreatedAt
	link.UpdatedAt = time.Now().UTC()
	result := service.database.WithContext(ctx).
		Model(&model.SocialLink{}).
		Where("id = ?", id).
		Updates(map[string]any{
			"platform": link.Platform, "display_name": link.DisplayName, "url": link.URL,
			"is_visible": link.IsVisible, "sort_order": link.SortOrder, "updated_at": link.UpdatedAt,
		})
	if result.Error != nil {
		return model.SocialLink{}, fmt.Errorf("update social link: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return model.SocialLink{}, ErrSiteContentNotFound
	}
	logSiteContentOperation("social_link_updated", administratorID, "social_link", id)
	return link, nil
}

func (service *SocialLinkAdminService) Delete(
	ctx context.Context,
	id uint64,
	administratorID uint64,
) error {
	result := service.database.WithContext(ctx).Delete(&model.SocialLink{}, id)
	if result.Error != nil {
		return fmt.Errorf("delete social link: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return ErrSiteContentNotFound
	}
	logSiteContentOperation("social_link_deleted", administratorID, "social_link", id)
	return nil
}

func (service *SocialLinkAdminService) Reorder(
	ctx context.Context,
	items []ResourceOrder,
	administratorID uint64,
) error {
	if err := validateResourceOrder(items, maximumSocialLinks); err != nil {
		return err
	}
	now := time.Now().UTC()
	err := service.database.WithContext(ctx).Transaction(func(transaction *gorm.DB) error {
		for _, item := range items {
			result := transaction.Model(&model.SocialLink{}).
				Where("id = ?", item.ID).
				Updates(map[string]any{"sort_order": item.SortOrder, "updated_at": now})
			if result.Error != nil {
				return fmt.Errorf("reorder social link: %w", result.Error)
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
	logSiteContentOperation("social_links_reordered", administratorID, "social_links", 0)
	return nil
}

func socialLinkFromInput(input SocialLinkInput) (model.SocialLink, error) {
	link := model.SocialLink{
		Platform: input.Platform, DisplayName: input.DisplayName, URL: input.URL,
		IsVisible: input.IsVisible, SortOrder: input.SortOrder,
	}
	if err := model.ValidateSocialLinks([]model.SocialLink{link}); err != nil {
		return model.SocialLink{}, fmt.Errorf("%w: %v", ErrInvalidSiteContent, err)
	}
	return link, nil
}
