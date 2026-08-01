package router

import (
	"bytes"
	"net/http"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"personal-blog/backend/internal/handler"
	"personal-blog/backend/internal/middleware"
	"personal-blog/backend/internal/service"
)

func TestSiteSettingsRequireAuthentication(t *testing.T) {
	engine, _ := newSiteSettingsAdminTestRouter(t)
	for _, method := range []string{http.MethodGet, http.MethodPut} {
		recorder := performSiteContentRequest(t, engine, method, "/api/v1/admin/site/settings", nil, false)
		if recorder.Code != http.StatusUnauthorized {
			t.Fatalf("expected status %d, got %d: %s", http.StatusUnauthorized, recorder.Code, recorder.Body.String())
		}
		assertErrorCode(t, recorder.Body.Bytes(), "UNAUTHORIZED")
	}
}

func TestAdministratorReadsAndUpdatesSiteSettings(t *testing.T) {
	engine, mock := newSiteSettingsAdminTestRouter(t)
	now := time.Now().UTC()
	expectSiteSettingQuery(mock, "Old Blog", now, false)

	loaded := performSiteContentRequest(t, engine, http.MethodGet, "/api/v1/admin/site/settings", nil, true)
	if loaded.Code != http.StatusOK || !bytes.Contains(loaded.Body.Bytes(), []byte(`"site_name":"Old Blog"`)) {
		t.Fatalf("unexpected site setting response: %d %s", loaded.Code, loaded.Body.String())
	}
	if loaded.Header().Get("Cache-Control") != "no-store" {
		t.Fatal("administrator site setting response must disable caching")
	}

	mock.ExpectBegin()
	expectSiteSettingQuery(mock, "Old Blog", now, true)
	mock.ExpectExec(`UPDATE "site_settings" SET .* WHERE id = \$[0-9]+`).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()
	updated := performSiteContentRequest(t, engine, http.MethodPut, "/api/v1/admin/site/settings", []byte(
		`{"site_name":"New Blog","site_short_name":"New","site_description":"A new description",`+
			`"title_suffix":"New Blog","logo_url":"https://cdn.example.com/logo.png",`+
			`"favicon_url":null,"default_share_image_url":null,"copyright_name":"New Blog",`+
			`"start_year":2024,"additional_text":"All rights reserved.","filing_number":null,`+
			`"filing_url":null,"show_technology":true,"technology_text":"Built with Vue 3 and Go"}`,
	), true)
	if updated.Code != http.StatusOK || !bytes.Contains(updated.Body.Bytes(), []byte(`"site_name":"New Blog"`)) {
		t.Fatalf("unexpected site setting update response: %d %s", updated.Code, updated.Body.String())
	}

	expectSiteSettingQuery(mock, "New Blog", now, false)
	mock.ExpectQuery(`SELECT \* FROM "navigation_items" WHERE is_visible = \$1 ORDER BY sort_order ASC, id ASC`).
		WithArgs(true).
		WillReturnRows(navigationRows())
	mock.ExpectQuery(`SELECT \* FROM "social_links" WHERE is_visible = \$1 ORDER BY sort_order ASC, id ASC`).
		WithArgs(true).
		WillReturnRows(socialLinkRows())
	public := performSiteContentRequest(t, engine, http.MethodGet, "/api/v1/site/config", nil, false)
	if public.Code != http.StatusOK || !bytes.Contains(public.Body.Bytes(), []byte(`"name":"New Blog"`)) {
		t.Fatalf("site setting change was not public: %d %s", public.Code, public.Body.String())
	}
}

func TestSiteSettingsRejectInvalidValues(t *testing.T) {
	engine, _ := newSiteSettingsAdminTestRouter(t)
	bodies := [][]byte{
		[]byte(`{"site_name":"` + string(bytes.Repeat([]byte("x"), 51)) + `","site_description":"Description","copyright_name":"Blog","show_technology":true}`),
		[]byte(`{"site_name":"Blog","site_description":"Description","logo_url":"javascript:alert(1)","copyright_name":"Blog","show_technology":true}`),
		[]byte(`{"site_name":"Blog","site_description":"Description","copyright_name":"Blog","show_technology":true,"unknown":true}`),
	}
	for _, body := range bodies {
		recorder := performSiteContentRequest(
			t, engine, http.MethodPut, "/api/v1/admin/site/settings", body, true,
		)
		if recorder.Code != http.StatusBadRequest {
			t.Fatalf("expected status %d, got %d: %s", http.StatusBadRequest, recorder.Code, recorder.Body.String())
		}
		assertErrorCode(t, recorder.Body.Bytes(), "VALIDATION_ERROR")
	}
}

func newSiteSettingsAdminTestRouter(t *testing.T) (*gin.Engine, sqlmock.Sqlmock) {
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
	authService := service.NewAuthService(database, testJWTSecret, 24)
	return New(Dependencies{
		AuthHandler:       handler.NewAuthHandler(authService),
		AuthMiddleware:    middleware.Authenticate(authService),
		PublicConfig:      handler.NewPublicConfigHandler(service.NewPublicConfigService(database)),
		SiteSettingsAdmin: handler.NewSiteSettingsAdminHandler(service.NewSiteSettingsAdminService(database)),
	}), mock
}

func expectSiteSettingQuery(mock sqlmock.Sqlmock, siteName string, now time.Time, locked bool) {
	query := `SELECT \* FROM "site_settings" WHERE id = \$1 LIMIT \$2`
	if locked {
		query += ` FOR UPDATE`
	}
	mock.ExpectQuery(query).
		WithArgs(uint64(1), 1).
		WillReturnRows(siteSettingRows().AddRow(
			uint64(1), siteName, "Short", "Description", siteName, nil, nil, nil,
			siteName, 2024, nil, nil, nil, true, "Built with Vue 3 and Go", now, now,
		))
}

func siteSettingRows() *sqlmock.Rows {
	return sqlmock.NewRows([]string{
		"id", "site_name", "site_short_name", "site_description", "title_suffix", "logo_url",
		"favicon_url", "default_share_image_url", "copyright_name", "start_year", "additional_text",
		"filing_number", "filing_url", "show_technology", "technology_text", "created_at", "updated_at",
	})
}
