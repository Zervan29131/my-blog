package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"sort"
	"time"

	"gorm.io/gorm"

	"personal-blog/backend/internal/model"
)

type PublicConfigService struct {
	database       *gorm.DB
	articleService *ArticleService
}

type PublicSiteConfig struct {
	Site        PublicSite             `json:"site"`
	Navigation  []PublicNavigationItem `json:"navigation"`
	SocialLinks []PublicSocialLink     `json:"social_links"`
	Footer      PublicFooter           `json:"footer"`
}

type PublicSite struct {
	Name                 string `json:"name"`
	ShortName            string `json:"short_name"`
	Description          string `json:"description"`
	TitleSuffix          string `json:"title_suffix"`
	LogoURL              string `json:"logo_url"`
	FaviconURL           string `json:"favicon_url"`
	DefaultShareImageURL string `json:"default_share_image_url"`
}

type PublicNavigationItem struct {
	Name         string `json:"name"`
	URL          string `json:"url"`
	LinkType     string `json:"link_type"`
	OpenInNewTab bool   `json:"open_in_new_tab"`
}

type PublicSocialLink struct {
	Platform    string `json:"platform"`
	DisplayName string `json:"display_name"`
	URL         string `json:"url"`
}

type PublicFooter struct {
	CopyrightName  string `json:"copyright_name"`
	StartYear      *int   `json:"start_year"`
	AdditionalText string `json:"additional_text"`
	FilingNumber   string `json:"filing_number"`
	FilingURL      string `json:"filing_url"`
	ShowTechnology bool   `json:"show_technology"`
	TechnologyText string `json:"technology_text"`
}

type PublicHomepage struct {
	Version uint64                 `json:"version"`
	Modules []PublicHomepageModule `json:"modules"`
}

type PublicHomepageModule struct {
	Type             string                              `json:"type"`
	SortOrder        int                                 `json:"sort_order"`
	Hero             *model.HeroModuleConfig             `json:"-"`
	About            *model.AboutModuleConfig            `json:"-"`
	FeaturedArticles *PublicFeaturedArticlesModuleConfig `json:"-"`
	LatestArticles   *PublicLatestArticlesModuleConfig   `json:"-"`
	TechStack        *PublicTechStackModuleConfig        `json:"-"`
	SocialLinks      *PublicSocialLinksModuleConfig      `json:"-"`
}

type PublicArticleSummary struct {
	ID           uint64     `json:"id"`
	Title        string     `json:"title"`
	Slug         string     `json:"slug"`
	Summary      string     `json:"summary"`
	PublishedAt  *time.Time `json:"published_at"`
	CommentCount int64      `json:"comment_count"`
}

type PublicFeaturedArticlesModuleConfig struct {
	Title       string                 `json:"title"`
	Description string                 `json:"description"`
	Limit       int                    `json:"limit"`
	Articles    []PublicArticleSummary `json:"articles"`
}

type PublicLatestArticlesModuleConfig struct {
	Title            string                 `json:"title"`
	Description      string                 `json:"description"`
	Limit            int                    `json:"limit"`
	ShowSummary      bool                   `json:"show_summary"`
	ShowDate         bool                   `json:"show_date"`
	ShowCommentCount bool                   `json:"show_comment_count"`
	ShowViewAll      bool                   `json:"show_view_all"`
	Articles         []PublicArticleSummary `json:"articles"`
}

type PublicTechStackModuleConfig struct {
	Title       string           `json:"title"`
	Description string           `json:"description"`
	Items       []PublicTechItem `json:"items"`
}

type PublicTechItem struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	IconURL     string `json:"icon_url"`
	URL         string `json:"url"`
}

type PublicSocialLinksModuleConfig struct {
	Title       string             `json:"title"`
	Description string             `json:"description"`
	Links       []PublicSocialLink `json:"links"`
}

func NewPublicConfigService(database *gorm.DB) *PublicConfigService {
	return &PublicConfigService{
		database:       database,
		articleService: NewArticleService(database),
	}
}

func (service *PublicConfigService) GetSiteConfig(ctx context.Context) (PublicSiteConfig, error) {
	setting, err := service.siteSetting(ctx)
	if err != nil {
		return PublicSiteConfig{}, err
	}

	var navigationItems []model.NavigationItem
	if err := service.database.WithContext(ctx).
		Where("is_visible = ?", true).
		Order("sort_order ASC, id ASC").
		Find(&navigationItems).Error; err != nil {
		return PublicSiteConfig{}, fmt.Errorf("list public navigation: %w", err)
	}

	var socialLinks []model.SocialLink
	if err := service.database.WithContext(ctx).
		Where("is_visible = ?", true).
		Order("sort_order ASC, id ASC").
		Find(&socialLinks).Error; err != nil {
		return PublicSiteConfig{}, fmt.Errorf("list public social links: %w", err)
	}

	return publicSiteConfig(setting, validPublicNavigation(navigationItems), validPublicSocialLinks(socialLinks)), nil
}

func (service *PublicConfigService) GetHomepage(ctx context.Context) (PublicHomepage, error) {
	version, config, err := service.publishedHomepageConfig(ctx)
	if err != nil {
		return PublicHomepage{}, err
	}
	return service.composeHomepage(ctx, version, config)
}

func (service *PublicConfigService) composeHomepage(
	ctx context.Context,
	version uint64,
	config model.HomepageConfig,
) (PublicHomepage, error) {
	var socialLinks []PublicSocialLink
	modules := make([]PublicHomepageModule, 0, len(config.Modules))
	for _, module := range config.Modules {
		if !module.Enabled {
			continue
		}

		publicModule, include, err := service.composeHomepageModule(ctx, module, &socialLinks)
		if err != nil {
			return PublicHomepage{}, err
		}
		if include {
			modules = append(modules, publicModule)
		}
	}

	sort.SliceStable(modules, func(left, right int) bool {
		if modules[left].SortOrder != modules[right].SortOrder {
			return modules[left].SortOrder < modules[right].SortOrder
		}
		return homepageModuleRank(modules[left].Type) < homepageModuleRank(modules[right].Type)
	})

	return PublicHomepage{Version: version, Modules: modules}, nil
}

func (service *PublicConfigService) siteSetting(ctx context.Context) (model.SiteSetting, error) {
	var setting model.SiteSetting
	result := service.database.WithContext(ctx).Where("id = ?", uint64(1)).Limit(1).Find(&setting)
	if result.Error != nil {
		return model.SiteSetting{}, fmt.Errorf("find public site setting: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return model.DefaultSiteSetting(time.Now()), nil
	}
	if err := model.ValidateSiteSetting(setting); err != nil {
		slog.Warn("public site setting is invalid; using defaults", "error", err)
		return model.DefaultSiteSetting(time.Now()), nil
	}
	return setting, nil
}

func (service *PublicConfigService) publishedHomepageConfig(
	ctx context.Context,
) (uint64, model.HomepageConfig, error) {
	var version model.HomepageVersion
	result := service.database.WithContext(ctx).
		Where("status = ?", model.HomepageStatusPublished).
		Limit(1).
		Find(&version)
	if result.Error != nil {
		return 0, model.HomepageConfig{}, fmt.Errorf("find published homepage config: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return 0, model.DefaultHomepageConfig(), nil
	}

	config, err := model.DecodeHomepageConfig(version.ConfigJSON)
	if err != nil {
		slog.Warn("published homepage config is invalid; using defaults", "version", version.VersionNumber, "error", err)
		return 0, model.DefaultHomepageConfig(), nil
	}
	return version.VersionNumber, config, nil
}

func (service *PublicConfigService) composeHomepageModule(
	ctx context.Context,
	module model.HomepageModule,
	cachedSocialLinks *[]PublicSocialLink,
) (PublicHomepageModule, bool, error) {
	publicModule := PublicHomepageModule{Type: module.Type, SortOrder: module.SortOrder}
	switch module.Type {
	case model.HomepageModuleHero:
		config := *module.Hero
		publicModule.Hero = &config
	case model.HomepageModuleAbout:
		config := *module.About
		publicModule.About = &config
	case model.HomepageModuleFeaturedArticles:
		articles, err := service.featuredArticles(ctx, module.FeaturedArticles.Limit)
		if err != nil {
			return PublicHomepageModule{}, false, err
		}
		if len(articles) == 0 {
			return PublicHomepageModule{}, false, nil
		}
		publicModule.FeaturedArticles = &PublicFeaturedArticlesModuleConfig{
			Title:       module.FeaturedArticles.Title,
			Description: module.FeaturedArticles.Description,
			Limit:       module.FeaturedArticles.Limit,
			Articles:    articles,
		}
	case model.HomepageModuleLatestArticles:
		page, err := service.articleService.ListPublished(ctx, 1, module.LatestArticles.Limit)
		if err != nil {
			return PublicHomepageModule{}, false, fmt.Errorf("compose latest articles module: %w", err)
		}
		publicModule.LatestArticles = &PublicLatestArticlesModuleConfig{
			Title:            module.LatestArticles.Title,
			Description:      module.LatestArticles.Description,
			Limit:            module.LatestArticles.Limit,
			ShowSummary:      module.LatestArticles.ShowSummary,
			ShowDate:         module.LatestArticles.ShowDate,
			ShowCommentCount: module.LatestArticles.ShowCommentCount,
			ShowViewAll:      module.LatestArticles.ShowViewAll,
			Articles:         publicArticleSummaries(page.Items),
		}
	case model.HomepageModuleTechStack:
		items := make([]model.TechItem, 0, len(module.TechStack.Items))
		for _, item := range module.TechStack.Items {
			if item.IsVisible {
				items = append(items, item)
			}
		}
		sort.SliceStable(items, func(left, right int) bool {
			if items[left].SortOrder != items[right].SortOrder {
				return items[left].SortOrder < items[right].SortOrder
			}
			return items[left].Name < items[right].Name
		})
		publicItems := make([]PublicTechItem, 0, len(items))
		for _, item := range items {
			publicItems = append(publicItems, PublicTechItem{
				Name: item.Name, Description: item.Description, IconURL: item.IconURL, URL: item.URL,
			})
		}
		publicModule.TechStack = &PublicTechStackModuleConfig{
			Title: module.TechStack.Title, Description: module.TechStack.Description, Items: publicItems,
		}
	case model.HomepageModuleSocialLinks:
		if *cachedSocialLinks == nil {
			links, err := service.publicSocialLinks(ctx)
			if err != nil {
				return PublicHomepageModule{}, false, err
			}
			*cachedSocialLinks = links
		}
		publicModule.SocialLinks = &PublicSocialLinksModuleConfig{
			Title: module.SocialLinks.Title, Description: module.SocialLinks.Description, Links: *cachedSocialLinks,
		}
	default:
		slog.Warn("skipping unsupported public homepage module", "type", module.Type)
		return PublicHomepageModule{}, false, nil
	}
	return publicModule, true, nil
}

func (service *PublicConfigService) featuredArticles(
	ctx context.Context,
	limit int,
) ([]PublicArticleSummary, error) {
	articles := make([]PublicArticleSummary, 0)
	result := service.database.WithContext(ctx).
		Table("featured_articles").
		Select(
			"articles.id, articles.title, articles.slug, articles.summary, articles.published_at, "+
				"(SELECT COUNT(*) FROM comments WHERE comments.article_id = articles.id "+
				"AND comments.status = ?) AS comment_count",
			model.CommentStatusApproved,
		).
		Joins("JOIN articles ON articles.id = featured_articles.article_id").
		Where("featured_articles.is_visible = ? AND articles.status = ?", true, model.ArticleStatusPublished).
		Order("featured_articles.sort_order ASC, featured_articles.id ASC").
		Limit(limit).
		Scan(&articles)
	if result.Error != nil {
		return nil, fmt.Errorf("list public featured articles: %w", result.Error)
	}
	return articles, nil
}

func (service *PublicConfigService) publicSocialLinks(ctx context.Context) ([]PublicSocialLink, error) {
	var links []model.SocialLink
	if err := service.database.WithContext(ctx).
		Where("is_visible = ?", true).
		Order("sort_order ASC, id ASC").
		Find(&links).Error; err != nil {
		return nil, fmt.Errorf("list homepage social links: %w", err)
	}
	return validPublicSocialLinks(links), nil
}

func (module PublicHomepageModule) MarshalJSON() ([]byte, error) {
	var config any
	switch module.Type {
	case model.HomepageModuleHero:
		config = module.Hero
	case model.HomepageModuleAbout:
		config = module.About
	case model.HomepageModuleFeaturedArticles:
		config = module.FeaturedArticles
	case model.HomepageModuleLatestArticles:
		config = module.LatestArticles
	case model.HomepageModuleTechStack:
		config = module.TechStack
	case model.HomepageModuleSocialLinks:
		config = module.SocialLinks
	default:
		return nil, fmt.Errorf("unsupported public homepage module type %q", module.Type)
	}
	if config == nil {
		return nil, errors.New("public homepage module config is missing")
	}
	return json.Marshal(struct {
		Type      string `json:"type"`
		SortOrder int    `json:"sort_order"`
		Config    any    `json:"config"`
	}{Type: module.Type, SortOrder: module.SortOrder, Config: config})
}

func publicSiteConfig(
	setting model.SiteSetting,
	navigation []PublicNavigationItem,
	socialLinks []PublicSocialLink,
) PublicSiteConfig {
	return PublicSiteConfig{
		Site: PublicSite{
			Name: setting.SiteName, ShortName: stringValue(setting.SiteShortName), Description: setting.SiteDescription,
			TitleSuffix: stringValue(setting.TitleSuffix), LogoURL: stringValue(setting.LogoURL),
			FaviconURL: stringValue(setting.FaviconURL), DefaultShareImageURL: stringValue(setting.DefaultShareImageURL),
		},
		Navigation:  navigation,
		SocialLinks: socialLinks,
		Footer: PublicFooter{
			CopyrightName: setting.CopyrightName, StartYear: setting.StartYear,
			AdditionalText: stringValue(setting.AdditionalText), FilingNumber: stringValue(setting.FilingNumber),
			FilingURL: stringValue(setting.FilingURL), ShowTechnology: setting.ShowTechnology,
			TechnologyText: stringValue(setting.TechnologyText),
		},
	}
}

func validPublicNavigation(items []model.NavigationItem) []PublicNavigationItem {
	result := make([]PublicNavigationItem, 0, len(items))
	for _, item := range items {
		if len(result) == 10 {
			break
		}
		if err := model.ValidateNavigationItems([]model.NavigationItem{item}); err != nil {
			slog.Warn("skipping invalid public navigation item", "id", item.ID, "error", err)
			continue
		}
		result = append(result, PublicNavigationItem{
			Name: item.Name, URL: item.URL, LinkType: item.LinkType, OpenInNewTab: item.OpenInNewTab,
		})
	}
	return result
}

func validPublicSocialLinks(links []model.SocialLink) []PublicSocialLink {
	result := make([]PublicSocialLink, 0, len(links))
	for _, link := range links {
		if len(result) == 15 {
			break
		}
		if err := model.ValidateSocialLinks([]model.SocialLink{link}); err != nil {
			slog.Warn("skipping invalid public social link", "id", link.ID, "error", err)
			continue
		}
		result = append(result, PublicSocialLink{
			Platform: link.Platform, DisplayName: link.DisplayName, URL: link.URL,
		})
	}
	return result
}

func publicArticleSummaries(articles []model.Article) []PublicArticleSummary {
	result := make([]PublicArticleSummary, 0, len(articles))
	for _, article := range articles {
		result = append(result, PublicArticleSummary{
			ID: article.ID, Title: article.Title, Slug: article.Slug, Summary: article.Summary,
			PublishedAt: article.PublishedAt, CommentCount: article.CommentCount,
		})
	}
	return result
}

func homepageModuleRank(moduleType string) int {
	switch moduleType {
	case model.HomepageModuleHero:
		return 0
	case model.HomepageModuleAbout:
		return 1
	case model.HomepageModuleFeaturedArticles:
		return 2
	case model.HomepageModuleLatestArticles:
		return 3
	case model.HomepageModuleTechStack:
		return 4
	case model.HomepageModuleSocialLinks:
		return 5
	default:
		return 6
	}
}

func stringValue(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}
