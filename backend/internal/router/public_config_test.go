package router

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"personal-blog/backend/internal/handler"
	"personal-blog/backend/internal/model"
	"personal-blog/backend/internal/service"
)

func TestPublicSiteConfigReturnsOnlyPublicFields(t *testing.T) {
	engine, mock := newPublicConfigTestRouter(t)
	now := time.Now().UTC()
	shortName := "Journal"
	titleSuffix := "Example Blog"
	technology := "Built with Vue 3 and Go"
	startYear := 2024

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "site_settings" WHERE id = $1 LIMIT $2`)).
		WithArgs(uint64(1), 1).
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "site_name", "site_short_name", "site_description", "title_suffix", "logo_url",
			"favicon_url", "default_share_image_url", "copyright_name", "start_year", "additional_text",
			"filing_number", "filing_url", "show_technology", "technology_text", "created_at", "updated_at",
		}).AddRow(
			uint64(1), "Example Blog", shortName, "A public description", titleSuffix, nil,
			nil, nil, "Example Blog", startYear, nil, nil, nil, true, technology, now, now,
		))
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "navigation_items" WHERE is_visible = $1 ORDER BY sort_order ASC, id ASC`)).
		WithArgs(true).
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "name", "url", "link_type", "open_in_new_tab", "is_visible", "sort_order", "created_at", "updated_at",
		}).AddRow(uint64(2), "Archive", "/archive", model.LinkTypeInternal, false, true, 20, now, now))
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "social_links" WHERE is_visible = $1 ORDER BY sort_order ASC, id ASC`)).
		WithArgs(true).
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "platform", "display_name", "url", "is_visible", "sort_order", "created_at", "updated_at",
		}).AddRow(uint64(3), model.SocialPlatformGitHub, "GitHub", "https://github.com/example", true, 10, now, now))

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/api/v1/site/config", nil)
	engine.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d: %s", http.StatusOK, recorder.Code, recorder.Body.String())
	}
	var response struct {
		Data struct {
			Site        map[string]any   `json:"site"`
			Navigation  []map[string]any `json:"navigation"`
			SocialLinks []map[string]any `json:"social_links"`
			Footer      map[string]any   `json:"footer"`
		} `json:"data"`
	}
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("decode site config response: %v", err)
	}
	if response.Data.Site["name"] != "Example Blog" || response.Data.Footer["copyright_name"] != "Example Blog" {
		t.Fatalf("unexpected site config: %+v", response.Data)
	}
	if len(response.Data.Navigation) != 1 || response.Data.Navigation[0]["name"] != "Archive" {
		t.Fatalf("unexpected navigation: %+v", response.Data.Navigation)
	}
	if len(response.Data.SocialLinks) != 1 || response.Data.SocialLinks[0]["platform"] != model.SocialPlatformGitHub {
		t.Fatalf("unexpected social links: %+v", response.Data.SocialLinks)
	}
	for _, item := range []map[string]any{response.Data.Site, response.Data.Navigation[0], response.Data.SocialLinks[0], response.Data.Footer} {
		for _, internalField := range []string{"id", "is_visible", "sort_order", "created_at", "updated_at"} {
			if _, exists := item[internalField]; exists {
				t.Fatalf("public site response exposed internal field %q: %+v", internalField, item)
			}
		}
	}
}

func TestPublicHomepageSortsModulesAndComposesPublishedFeaturedArticles(t *testing.T) {
	engine, mock := newPublicConfigTestRouter(t)
	config := model.DefaultHomepageConfig()
	for index := range config.Modules {
		config.Modules[index].Enabled = false
	}
	moduleByType(config.Modules, model.HomepageModuleFeaturedArticles).Enabled = true
	moduleByType(config.Modules, model.HomepageModuleFeaturedArticles).SortOrder = 10
	moduleByType(config.Modules, model.HomepageModuleTechStack).Enabled = true
	moduleByType(config.Modules, model.HomepageModuleTechStack).SortOrder = 20
	moduleByType(config.Modules, model.HomepageModuleHero).Enabled = true
	moduleByType(config.Modules, model.HomepageModuleHero).SortOrder = 30
	document, err := model.EncodeHomepageConfig(config)
	if err != nil {
		t.Fatalf("encode homepage fixture: %v", err)
	}
	publishedAt := time.Now().UTC()

	expectPublishedHomepageVersion(mock, 7, document)
	mock.ExpectQuery(`SELECT articles\.id, articles\.title, articles\.slug, articles\.summary, articles\.published_at, \(SELECT COUNT\(\*\) FROM comments WHERE comments\.article_id = articles\.id AND comments\.status = \$1\) AS comment_count FROM "featured_articles" JOIN articles ON articles\.id = featured_articles\.article_id WHERE featured_articles\.is_visible = \$2 AND articles\.status = \$3 ORDER BY featured_articles\.sort_order ASC, featured_articles\.id ASC LIMIT \$4`).
		WithArgs(model.CommentStatusApproved, true, model.ArticleStatusPublished, 3).
		WillReturnRows(sqlmock.NewRows([]string{"id", "title", "slug", "summary", "published_at", "comment_count"}).
			AddRow(uint64(9), "Published feature", "published-feature", "Summary", publishedAt, int64(2)))

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/api/v1/homepage", nil)
	engine.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d: %s", http.StatusOK, recorder.Code, recorder.Body.String())
	}
	var response homepageAPIResponse
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("decode homepage response: %v", err)
	}
	if response.Data.Version != 7 || len(response.Data.Modules) != 3 {
		t.Fatalf("unexpected homepage metadata: %+v", response.Data)
	}
	wantTypes := []string{model.HomepageModuleFeaturedArticles, model.HomepageModuleTechStack, model.HomepageModuleHero}
	for index, wantType := range wantTypes {
		if response.Data.Modules[index].Type != wantType {
			t.Fatalf("module %d: expected %q, got %q", index, wantType, response.Data.Modules[index].Type)
		}
	}
	var featured struct {
		Articles []service.PublicArticleSummary `json:"articles"`
	}
	if err := json.Unmarshal(response.Data.Modules[0].Config, &featured); err != nil {
		t.Fatalf("decode featured module: %v", err)
	}
	if len(featured.Articles) != 1 || featured.Articles[0].Slug != "published-feature" {
		t.Fatalf("unexpected featured articles: %+v", featured.Articles)
	}
	for _, module := range response.Data.Modules {
		if module.Enabled != nil {
			t.Fatalf("public homepage exposed module enabled state: %s", recorder.Body.String())
		}
	}
	if strings.Contains(recorder.Body.String(), `"status"`) {
		t.Fatalf("public homepage exposed a private field: %s", recorder.Body.String())
	}
}

func TestPublicHomepageFallsBackWhenPublishedConfigIsMissing(t *testing.T) {
	engine, mock := newPublicConfigTestRouter(t)
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "homepage_versions" WHERE status = $1 LIMIT $2`)).
		WithArgs(model.HomepageStatusPublished, 1).
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "status", "version_number", "config_json", "published_at", "created_at", "updated_at",
		}))
	mock.ExpectQuery(`SELECT articles\..*FROM "featured_articles".*WHERE featured_articles\.is_visible = \$2 AND articles\.status = \$3.*LIMIT \$4`).
		WithArgs(model.CommentStatusApproved, true, model.ArticleStatusPublished, 3).
		WillReturnRows(sqlmock.NewRows([]string{"id", "title", "slug", "summary", "published_at", "comment_count"}))
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "articles" WHERE status = $1`)).
		WithArgs(model.ArticleStatusPublished).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))
	mock.ExpectQuery(`SELECT articles\.id, articles\.title, articles\.slug, articles\.summary, articles\.published_at, .* FROM "articles" WHERE status = \$2 ORDER BY published_at DESC, id DESC LIMIT \$3`).
		WithArgs(model.CommentStatusApproved, model.ArticleStatusPublished, 10).
		WillReturnRows(sqlmock.NewRows([]string{"id", "title", "slug", "summary", "published_at", "comment_count"}))
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "social_links" WHERE is_visible = $1 ORDER BY sort_order ASC, id ASC`)).
		WithArgs(true).
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "platform", "display_name", "url", "is_visible", "sort_order", "created_at", "updated_at",
		}))

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/api/v1/homepage", nil)
	engine.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d: %s", http.StatusOK, recorder.Code, recorder.Body.String())
	}
	var response homepageAPIResponse
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("decode fallback homepage response: %v", err)
	}
	if response.Data.Version != 0 {
		t.Fatalf("expected fallback version 0, got %d", response.Data.Version)
	}
	if len(response.Data.Modules) != 4 {
		t.Fatalf("expected four visible default modules after hiding empty featured, got %+v", response.Data.Modules)
	}
	if response.Data.Modules[0].Type != model.HomepageModuleHero ||
		response.Data.Modules[1].Type != model.HomepageModuleAbout ||
		response.Data.Modules[2].Type != model.HomepageModuleLatestArticles ||
		response.Data.Modules[3].Type != model.HomepageModuleSocialLinks {
		t.Fatalf("unexpected fallback modules: %+v", response.Data.Modules)
	}
}

type homepageAPIResponse struct {
	Data struct {
		Version uint64 `json:"version"`
		Modules []struct {
			Type      string          `json:"type"`
			SortOrder int             `json:"sort_order"`
			Enabled   *bool           `json:"enabled"`
			Config    json.RawMessage `json:"config"`
		} `json:"modules"`
	} `json:"data"`
}

func newPublicConfigTestRouter(t *testing.T) (*gin.Engine, sqlmock.Sqlmock) {
	t.Helper()
	sqlDatabase, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("create SQL mock: %v", err)
	}
	t.Cleanup(func() {
		mock.ExpectClose()
		if err := sqlDatabase.Close(); err != nil {
			t.Errorf("close SQL mock: %v", err)
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unmet SQL expectations: %v", err)
		}
	})

	database, err := gorm.Open(postgres.New(postgres.Config{Conn: sqlDatabase}), &gorm.Config{
		DisableAutomaticPing: true,
	})
	if err != nil {
		t.Fatalf("open GORM test database: %v", err)
	}
	return New(Dependencies{
		PublicConfig: handler.NewPublicConfigHandler(service.NewPublicConfigService(database)),
	}), mock
}

func moduleByType(modules []model.HomepageModule, moduleType string) *model.HomepageModule {
	for index := range modules {
		if modules[index].Type == moduleType {
			return &modules[index]
		}
	}
	panic("homepage fixture module not found: " + moduleType)
}

func expectPublishedHomepageVersion(mock sqlmock.Sqlmock, version uint64, config model.JSONDocument) {
	now := time.Now().UTC()
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "homepage_versions" WHERE status = $1 LIMIT $2`)).
		WithArgs(model.HomepageStatusPublished, 1).
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "status", "version_number", "config_json", "published_at", "created_at", "updated_at",
		}).AddRow(uint64(2), model.HomepageStatusPublished, version, string(config), now, now, now))
}
