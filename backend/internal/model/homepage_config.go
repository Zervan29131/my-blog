package model

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
)

const (
	HomepageModuleHero             = "hero"
	HomepageModuleAbout            = "about"
	HomepageModuleFeaturedArticles = "featured_articles"
	HomepageModuleLatestArticles   = "latest_articles"
	HomepageModuleTechStack        = "tech_stack"
	HomepageModuleSocialLinks      = "social_links"

	HeroLayoutLeft   = "left"
	HeroLayoutCenter = "center"

	AboutImageLeft  = "left"
	AboutImageRight = "right"
	AboutImageNone  = "none"
)

type HomepageConfig struct {
	Modules []HomepageModule `json:"modules"`
}

type HomepageModule struct {
	Type             string                        `json:"type"`
	Enabled          bool                          `json:"enabled"`
	SortOrder        int                           `json:"sort_order"`
	Hero             *HeroModuleConfig             `json:"-"`
	About            *AboutModuleConfig            `json:"-"`
	FeaturedArticles *FeaturedArticlesModuleConfig `json:"-"`
	LatestArticles   *LatestArticlesModuleConfig   `json:"-"`
	TechStack        *TechStackModuleConfig        `json:"-"`
	SocialLinks      *SocialLinksModuleConfig      `json:"-"`
}

type HeroModuleConfig struct {
	Eyebrow            string         `json:"eyebrow"`
	Title              string         `json:"title"`
	HighlightText      string         `json:"highlight_text"`
	Description        string         `json:"description"`
	ImageURL           string         `json:"image_url"`
	BackgroundImageURL string         `json:"background_image_url"`
	Layout             string         `json:"layout"`
	PrimaryButton      HomepageButton `json:"primary_button"`
	SecondaryButton    HomepageButton `json:"secondary_button"`
}

type HomepageButton struct {
	Enabled      bool   `json:"enabled"`
	Text         string `json:"text"`
	URL          string `json:"url"`
	LinkType     string `json:"link_type"`
	OpenInNewTab bool   `json:"open_in_new_tab"`
}

type AboutModuleConfig struct {
	Title         string `json:"title"`
	Description   string `json:"description"`
	Content       string `json:"content"`
	ImageURL      string `json:"image_url"`
	ImagePosition string `json:"image_position"`
}

type FeaturedArticlesModuleConfig struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Limit       int    `json:"limit"`
}

type LatestArticlesModuleConfig struct {
	Title            string `json:"title"`
	Description      string `json:"description"`
	Limit            int    `json:"limit"`
	ShowSummary      bool   `json:"show_summary"`
	ShowDate         bool   `json:"show_date"`
	ShowCommentCount bool   `json:"show_comment_count"`
	ShowViewAll      bool   `json:"show_view_all"`
}

type TechStackModuleConfig struct {
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Items       []TechItem `json:"items"`
}

type TechItem struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	IconURL     string `json:"icon_url"`
	URL         string `json:"url"`
	IsVisible   bool   `json:"is_visible"`
	SortOrder   int    `json:"sort_order"`
}

type SocialLinksModuleConfig struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

func (module HomepageModule) MarshalJSON() ([]byte, error) {
	var config any
	switch module.Type {
	case HomepageModuleHero:
		config = module.Hero
	case HomepageModuleAbout:
		config = module.About
	case HomepageModuleFeaturedArticles:
		config = module.FeaturedArticles
	case HomepageModuleLatestArticles:
		config = module.LatestArticles
	case HomepageModuleTechStack:
		config = module.TechStack
	case HomepageModuleSocialLinks:
		config = module.SocialLinks
	default:
		return nil, fmt.Errorf("unsupported homepage module type %q", module.Type)
	}
	if config == nil {
		return nil, fmt.Errorf("homepage module %q has no config", module.Type)
	}
	return json.Marshal(struct {
		Type      string `json:"type"`
		Enabled   bool   `json:"enabled"`
		SortOrder int    `json:"sort_order"`
		Config    any    `json:"config"`
	}{
		Type: module.Type, Enabled: module.Enabled, SortOrder: module.SortOrder, Config: config,
	})
}

func (module *HomepageModule) UnmarshalJSON(data []byte) error {
	var raw struct {
		Type      string          `json:"type"`
		Enabled   bool            `json:"enabled"`
		SortOrder int             `json:"sort_order"`
		Config    json.RawMessage `json:"config"`
	}
	if err := decodeStrictJSON(data, &raw); err != nil {
		return fmt.Errorf("decode homepage module: %w", err)
	}
	if len(raw.Config) == 0 || bytes.Equal(raw.Config, []byte("null")) {
		return fmt.Errorf("homepage module %q config is required", raw.Type)
	}

	*module = HomepageModule{Type: raw.Type, Enabled: raw.Enabled, SortOrder: raw.SortOrder}
	switch raw.Type {
	case HomepageModuleHero:
		module.Hero = &HeroModuleConfig{}
		return decodeStrictJSON(raw.Config, module.Hero)
	case HomepageModuleAbout:
		module.About = &AboutModuleConfig{}
		return decodeStrictJSON(raw.Config, module.About)
	case HomepageModuleFeaturedArticles:
		module.FeaturedArticles = &FeaturedArticlesModuleConfig{}
		return decodeStrictJSON(raw.Config, module.FeaturedArticles)
	case HomepageModuleLatestArticles:
		module.LatestArticles = &LatestArticlesModuleConfig{}
		return decodeStrictJSON(raw.Config, module.LatestArticles)
	case HomepageModuleTechStack:
		module.TechStack = &TechStackModuleConfig{}
		return decodeStrictJSON(raw.Config, module.TechStack)
	case HomepageModuleSocialLinks:
		module.SocialLinks = &SocialLinksModuleConfig{}
		return decodeStrictJSON(raw.Config, module.SocialLinks)
	default:
		return fmt.Errorf("unsupported homepage module type %q", raw.Type)
	}
}

func DecodeHomepageConfig(data []byte) (HomepageConfig, error) {
	var config HomepageConfig
	if err := decodeStrictJSON(data, &config); err != nil {
		return HomepageConfig{}, fmt.Errorf("decode homepage config: %w", err)
	}
	if err := ValidateHomepageConfig(config); err != nil {
		return HomepageConfig{}, err
	}
	return config, nil
}

func EncodeHomepageConfig(config HomepageConfig) (JSONDocument, error) {
	if err := ValidateHomepageConfig(config); err != nil {
		return nil, err
	}
	data, err := json.Marshal(config)
	if err != nil {
		return nil, fmt.Errorf("encode homepage config: %w", err)
	}
	return JSONDocument(data), nil
}

func decodeStrictJSON(data []byte, target any) error {
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(target); err != nil {
		return err
	}
	var trailing any
	if err := decoder.Decode(&trailing); err != io.EOF {
		if err != nil {
			return err
		}
		return fmt.Errorf("multiple JSON values are not allowed")
	}
	return nil
}
