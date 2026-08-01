package model

import (
	"strings"
	"testing"
	"time"
)

func TestDefaultCMSConfigurationIsValidAndRoundTrips(t *testing.T) {
	now := time.Date(2026, time.August, 1, 0, 0, 0, 0, time.UTC)
	if err := ValidateSiteSetting(DefaultSiteSetting(now)); err != nil {
		t.Fatalf("validate default site setting: %v", err)
	}
	if err := ValidateNavigationItems(DefaultNavigationItems()); err != nil {
		t.Fatalf("validate default navigation: %v", err)
	}
	if err := ValidateSocialLinks(DefaultSocialLinks()); err != nil {
		t.Fatalf("validate default social links: %v", err)
	}

	encoded, err := EncodeHomepageConfig(DefaultHomepageConfig())
	if err != nil {
		t.Fatalf("encode default homepage config: %v", err)
	}
	decoded, err := DecodeHomepageConfig(encoded)
	if err != nil {
		t.Fatalf("decode default homepage config: %v", err)
	}
	if len(decoded.Modules) != 6 {
		t.Fatalf("expected six fixed modules, got %d", len(decoded.Modules))
	}
	if decoded.Modules[0].Type != HomepageModuleHero || decoded.Modules[0].Hero == nil {
		t.Fatalf("expected a typed hero module, got %#v", decoded.Modules[0])
	}
}

func TestDecodeHomepageConfigRejectsUnknownFields(t *testing.T) {
	encoded, err := EncodeHomepageConfig(DefaultHomepageConfig())
	if err != nil {
		t.Fatalf("encode default homepage config: %v", err)
	}
	withUnknownField := strings.Replace(
		string(encoded),
		`"eyebrow":`,
		`"arbitrary_html":"<script>alert(1)</script>","eyebrow":`,
		1,
	)
	if _, err := DecodeHomepageConfig([]byte(withUnknownField)); err == nil {
		t.Fatal("expected an unknown module config field to be rejected")
	}
}

func TestHomepageConfigRejectsUnknownOrMissingModules(t *testing.T) {
	config := DefaultHomepageConfig()
	config.Modules[0].Type = "custom_html"
	if err := ValidateHomepageConfig(config); err == nil {
		t.Fatal("expected an unknown module type to be rejected")
	}

	config = DefaultHomepageConfig()
	config.Modules = config.Modules[:len(config.Modules)-1]
	if err := ValidateHomepageConfig(config); err == nil {
		t.Fatal("expected a missing fixed module to be rejected")
	}
}

func TestHomepageConfigRejectsDangerousURLsAndInvalidLimits(t *testing.T) {
	config := DefaultHomepageConfig()
	config.Modules[0].Hero.PrimaryButton.LinkType = LinkTypeExternal
	config.Modules[0].Hero.PrimaryButton.URL = "javascript:alert(1)"
	if err := ValidateHomepageConfig(config); err == nil {
		t.Fatal("expected a dangerous hero button URL to be rejected")
	}

	config = DefaultHomepageConfig()
	config.Modules[3].LatestArticles.Limit = 21
	if err := ValidateHomepageConfig(config); err == nil {
		t.Fatal("expected an invalid latest article limit to be rejected")
	}
}

func TestCMSResourceValidation(t *testing.T) {
	now := time.Date(2026, time.August, 1, 0, 0, 0, 0, time.UTC)
	setting := DefaultSiteSetting(now)
	unsafeImage := "data:text/html,<script>alert(1)</script>"
	setting.LogoURL = &unsafeImage
	if err := ValidateSiteSetting(setting); err == nil {
		t.Fatal("expected an unsafe image URL to be rejected")
	}

	if err := ValidateNavigationItems([]NavigationItem{{
		Name: "危险链接", URL: "javascript:alert(1)", LinkType: LinkTypeExternal,
	}}); err == nil {
		t.Fatal("expected a dangerous navigation URL to be rejected")
	}

	if err := ValidateSocialLinks([]SocialLink{{
		Platform: SocialPlatformEmail, DisplayName: "邮箱", URL: "mailto:hello@example.com",
	}}); err != nil {
		t.Fatalf("expected a mailto social link to be accepted: %v", err)
	}
	if err := ValidateSocialLinks([]SocialLink{{
		Platform: SocialPlatformCustom, DisplayName: "危险链接", URL: "javascript:alert(1)",
	}}); err == nil {
		t.Fatal("expected a dangerous social link URL to be rejected")
	}
}

func TestJSONDocumentValueAndScan(t *testing.T) {
	document := JSONDocument(`{"modules":[]}`)
	value, err := document.Value()
	if err != nil {
		t.Fatalf("get JSON value: %v", err)
	}
	if value != `{"modules":[]}` {
		t.Fatalf("unexpected JSON value: %v", value)
	}

	var scanned JSONDocument
	if err := scanned.Scan([]byte(`{"version":1}`)); err != nil {
		t.Fatalf("scan JSON document: %v", err)
	}
	if string(scanned) != `{"version":1}` {
		t.Fatalf("unexpected scanned JSON: %s", scanned)
	}
	if err := scanned.Scan([]byte(`not-json`)); err == nil {
		t.Fatal("expected invalid JSON to be rejected")
	}
}
