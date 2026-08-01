package model

import (
	"fmt"
	"net/url"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"
)

var homepageModuleTypes = []string{
	HomepageModuleHero,
	HomepageModuleAbout,
	HomepageModuleFeaturedArticles,
	HomepageModuleLatestArticles,
	HomepageModuleTechStack,
	HomepageModuleSocialLinks,
}

func ValidateHomepageConfig(config HomepageConfig) error {
	if len(config.Modules) != len(homepageModuleTypes) {
		return fmt.Errorf("homepage config must contain exactly %d modules", len(homepageModuleTypes))
	}

	seen := make(map[string]struct{}, len(config.Modules))
	for index, module := range config.Modules {
		if _, exists := seen[module.Type]; exists {
			return fmt.Errorf("homepage module %q is duplicated", module.Type)
		}
		seen[module.Type] = struct{}{}
		if err := validateHomepageModule(module); err != nil {
			return fmt.Errorf("modules[%d]: %w", index, err)
		}
	}
	for _, moduleType := range homepageModuleTypes {
		if _, exists := seen[moduleType]; !exists {
			return fmt.Errorf("homepage module %q is required", moduleType)
		}
	}
	return nil
}

func validateHomepageModule(module HomepageModule) error {
	switch module.Type {
	case HomepageModuleHero:
		return validateHeroConfig(module.Hero)
	case HomepageModuleAbout:
		return validateAboutConfig(module.About)
	case HomepageModuleFeaturedArticles:
		return validateFeaturedArticlesConfig(module.FeaturedArticles)
	case HomepageModuleLatestArticles:
		return validateLatestArticlesConfig(module.LatestArticles)
	case HomepageModuleTechStack:
		return validateTechStackConfig(module.TechStack)
	case HomepageModuleSocialLinks:
		return validateSocialLinksConfig(module.SocialLinks)
	default:
		return fmt.Errorf("unsupported module type %q", module.Type)
	}
}

func validateHeroConfig(config *HeroModuleConfig) error {
	if config == nil {
		return fmt.Errorf("hero config is required")
	}
	if err := validateOptionalText("hero eyebrow", config.Eyebrow, 50, false); err != nil {
		return err
	}
	if err := validateRequiredText("hero title", config.Title, 100, false); err != nil {
		return err
	}
	if err := validateOptionalText("hero highlight text", config.HighlightText, 50, false); err != nil {
		return err
	}
	if err := validateRequiredText("hero description", config.Description, 300, true); err != nil {
		return err
	}
	if err := validateOptionalImageURL("hero image URL", config.ImageURL); err != nil {
		return err
	}
	if err := validateOptionalImageURL("hero background image URL", config.BackgroundImageURL); err != nil {
		return err
	}
	if config.Layout != HeroLayoutLeft && config.Layout != HeroLayoutCenter {
		return fmt.Errorf("hero layout must be left or center")
	}
	if err := validateHomepageButton("primary button", config.PrimaryButton); err != nil {
		return err
	}
	return validateHomepageButton("secondary button", config.SecondaryButton)
}

func validateHomepageButton(name string, button HomepageButton) error {
	if button.LinkType != LinkTypeInternal && button.LinkType != LinkTypeExternal {
		return fmt.Errorf("%s link type must be internal or external", name)
	}
	if !button.Enabled {
		if button.Text != "" {
			return validateOptionalText(name+" text", button.Text, 20, false)
		}
		return nil
	}
	if err := validateRequiredText(name+" text", button.Text, 20, false); err != nil {
		return err
	}
	return validateLinkURL(name+" URL", button.URL, button.LinkType)
}

func validateAboutConfig(config *AboutModuleConfig) error {
	if config == nil {
		return fmt.Errorf("about config is required")
	}
	if err := validateRequiredText("about title", config.Title, 100, false); err != nil {
		return err
	}
	if err := validateOptionalText("about description", config.Description, 200, true); err != nil {
		return err
	}
	if err := validateRequiredText("about content", config.Content, 2000, true); err != nil {
		return err
	}
	if err := validateOptionalImageURL("about image URL", config.ImageURL); err != nil {
		return err
	}
	if config.ImagePosition != AboutImageLeft && config.ImagePosition != AboutImageRight && config.ImagePosition != AboutImageNone {
		return fmt.Errorf("about image position must be left, right, or none")
	}
	return nil
}

func validateFeaturedArticlesConfig(config *FeaturedArticlesModuleConfig) error {
	if config == nil {
		return fmt.Errorf("featured articles config is required")
	}
	if err := validateRequiredText("featured articles title", config.Title, 100, false); err != nil {
		return err
	}
	if err := validateOptionalText("featured articles description", config.Description, 200, true); err != nil {
		return err
	}
	if config.Limit < 1 || config.Limit > 10 {
		return fmt.Errorf("featured articles limit must be between 1 and 10")
	}
	return nil
}

func validateLatestArticlesConfig(config *LatestArticlesModuleConfig) error {
	if config == nil {
		return fmt.Errorf("latest articles config is required")
	}
	if err := validateRequiredText("latest articles title", config.Title, 100, false); err != nil {
		return err
	}
	if err := validateOptionalText("latest articles description", config.Description, 200, true); err != nil {
		return err
	}
	if config.Limit < 3 || config.Limit > 20 {
		return fmt.Errorf("latest articles limit must be between 3 and 20")
	}
	return nil
}

func validateTechStackConfig(config *TechStackModuleConfig) error {
	if config == nil {
		return fmt.Errorf("tech stack config is required")
	}
	if err := validateRequiredText("tech stack title", config.Title, 100, false); err != nil {
		return err
	}
	if err := validateOptionalText("tech stack description", config.Description, 200, true); err != nil {
		return err
	}
	if len(config.Items) > 20 {
		return fmt.Errorf("tech stack must not contain more than 20 items")
	}
	for index, item := range config.Items {
		if err := validateRequiredText("tech item name", item.Name, 30, false); err != nil {
			return fmt.Errorf("tech items[%d]: %w", index, err)
		}
		if err := validateOptionalText("tech item description", item.Description, 100, true); err != nil {
			return fmt.Errorf("tech items[%d]: %w", index, err)
		}
		if err := validateOptionalImageURL("tech item icon URL", item.IconURL); err != nil {
			return fmt.Errorf("tech items[%d]: %w", index, err)
		}
		if item.URL != "" {
			if err := validateExternalURL("tech item URL", item.URL); err != nil {
				return fmt.Errorf("tech items[%d]: %w", index, err)
			}
		}
	}
	return nil
}

func validateSocialLinksConfig(config *SocialLinksModuleConfig) error {
	if config == nil {
		return fmt.Errorf("social links config is required")
	}
	if err := validateRequiredText("social links title", config.Title, 100, false); err != nil {
		return err
	}
	return validateOptionalText("social links description", config.Description, 200, true)
}

func ValidateSiteSetting(setting SiteSetting) error {
	if err := validateRequiredText("site name", setting.SiteName, 50, false); err != nil {
		return err
	}
	if err := validateOptionalStringPointer("site short name", setting.SiteShortName, 20, false); err != nil {
		return err
	}
	if err := validateRequiredText("site description", setting.SiteDescription, 200, true); err != nil {
		return err
	}
	if err := validateOptionalStringPointer("title suffix", setting.TitleSuffix, 50, false); err != nil {
		return err
	}
	for name, value := range map[string]*string{
		"logo URL": setting.LogoURL, "favicon URL": setting.FaviconURL, "default share image URL": setting.DefaultShareImageURL,
	} {
		if value != nil && *value != "" {
			if err := validateOptionalImageURL(name, *value); err != nil {
				return err
			}
		}
	}
	if err := validateRequiredText("copyright name", setting.CopyrightName, 50, false); err != nil {
		return err
	}
	if setting.StartYear != nil && (*setting.StartYear < 1900 || *setting.StartYear > time.Now().UTC().Year()) {
		return fmt.Errorf("start year must be between 1900 and the current year")
	}
	if err := validateOptionalStringPointer("additional text", setting.AdditionalText, 200, true); err != nil {
		return err
	}
	if err := validateOptionalStringPointer("filing number", setting.FilingNumber, 100, false); err != nil {
		return err
	}
	if setting.FilingURL != nil && *setting.FilingURL != "" {
		if err := validateExternalURL("filing URL", *setting.FilingURL); err != nil {
			return err
		}
	}
	return validateOptionalStringPointer("technology text", setting.TechnologyText, 100, false)
}

func ValidateNavigationItems(items []NavigationItem) error {
	if len(items) > 10 {
		return fmt.Errorf("navigation must not contain more than 10 items")
	}
	for index, item := range items {
		if err := validateRequiredText("navigation name", item.Name, 20, false); err != nil {
			return fmt.Errorf("navigation items[%d]: %w", index, err)
		}
		if err := validateLinkURL("navigation URL", item.URL, item.LinkType); err != nil {
			return fmt.Errorf("navigation items[%d]: %w", index, err)
		}
	}
	return nil
}

func ValidateSocialLinks(links []SocialLink) error {
	if len(links) > 15 {
		return fmt.Errorf("social links must not contain more than 15 items")
	}
	allowedPlatforms := map[string]struct{}{
		SocialPlatformGitHub: {}, SocialPlatformEmail: {}, SocialPlatformLinkedIn: {}, SocialPlatformX: {},
		SocialPlatformWeibo: {}, SocialPlatformBilibili: {}, SocialPlatformZhihu: {}, SocialPlatformCustom: {},
	}
	for index, link := range links {
		if _, allowed := allowedPlatforms[link.Platform]; !allowed {
			return fmt.Errorf("social links[%d]: unsupported platform %q", index, link.Platform)
		}
		if err := validateRequiredText("social link display name", link.DisplayName, 30, false); err != nil {
			return fmt.Errorf("social links[%d]: %w", index, err)
		}
		if link.Platform == SocialPlatformEmail && strings.HasPrefix(strings.ToLower(link.URL), "mailto:") {
			if err := validateMailtoURL("social link URL", link.URL); err != nil {
				return fmt.Errorf("social links[%d]: %w", index, err)
			}
			continue
		}
		if err := validateExternalURL("social link URL", link.URL); err != nil {
			return fmt.Errorf("social links[%d]: %w", index, err)
		}
	}
	return nil
}

func validateRequiredText(name, value string, maximum int, allowLineBreaks bool) error {
	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("%s is required", name)
	}
	return validateOptionalText(name, value, maximum, allowLineBreaks)
}

func validateOptionalText(name, value string, maximum int, allowLineBreaks bool) error {
	if utf8.RuneCountInString(value) > maximum {
		return fmt.Errorf("%s must not exceed %d characters", name, maximum)
	}
	if strings.TrimSpace(value) != value {
		return fmt.Errorf("%s must not have surrounding whitespace", name)
	}
	if strings.IndexFunc(value, func(character rune) bool {
		if !unicode.IsControl(character) {
			return false
		}
		return !allowLineBreaks || (character != '\n' && character != '\r' && character != '\t')
	}) >= 0 {
		return fmt.Errorf("%s contains control characters", name)
	}
	return nil
}

func validateOptionalStringPointer(name string, value *string, maximum int, allowLineBreaks bool) error {
	if value == nil {
		return nil
	}
	return validateOptionalText(name, *value, maximum, allowLineBreaks)
}

func validateOptionalImageURL(name, value string) error {
	if value == "" {
		return nil
	}
	return validateExternalURL(name, value)
}

func validateLinkURL(name, value, linkType string) error {
	if linkType == LinkTypeInternal {
		if len(value) > 500 || !strings.HasPrefix(value, "/") || strings.HasPrefix(value, "//") || strings.ContainsAny(value, "\r\n") {
			return fmt.Errorf("%s must be a safe internal path", name)
		}
		parsed, err := url.Parse(value)
		if err != nil || parsed.Scheme != "" || parsed.Host != "" {
			return fmt.Errorf("%s must be a safe internal path", name)
		}
		return nil
	}
	if linkType == LinkTypeExternal {
		return validateExternalURL(name, value)
	}
	return fmt.Errorf("%s link type must be internal or external", name)
}

func validateExternalURL(name, value string) error {
	if len(value) == 0 || len(value) > 500 || strings.ContainsAny(value, "\r\n") {
		return fmt.Errorf("%s must be a valid HTTP or HTTPS URL", name)
	}
	parsed, err := url.ParseRequestURI(value)
	if err != nil || (parsed.Scheme != "http" && parsed.Scheme != "https") || parsed.Host == "" {
		return fmt.Errorf("%s must be a valid HTTP or HTTPS URL", name)
	}
	return nil
}

func validateMailtoURL(name, value string) error {
	if len(value) == 0 || len(value) > 500 || strings.ContainsAny(value, "\r\n") {
		return fmt.Errorf("%s must be a valid mailto URL", name)
	}
	parsed, err := url.Parse(value)
	if err != nil || parsed.Scheme != "mailto" || parsed.Opaque == "" || strings.Contains(parsed.Opaque, " ") {
		return fmt.Errorf("%s must be a valid mailto URL", name)
	}
	return nil
}
