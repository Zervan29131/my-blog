package router

import (
	"bytes"
	"encoding/json"
	"errors"
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
	"personal-blog/backend/internal/middleware"
	"personal-blog/backend/internal/model"
	"personal-blog/backend/internal/service"
)

func TestHomepageAdminEndpointsRequireAuthentication(t *testing.T) {
	engine, _ := newHomepageAdminTestRouter(t)
	tests := []struct {
		method string
		path   string
	}{
		{http.MethodGet, "/api/v1/admin/homepage/draft"},
		{http.MethodPut, "/api/v1/admin/homepage/draft"},
		{http.MethodGet, "/api/v1/admin/homepage/published"},
		{http.MethodPost, "/api/v1/admin/homepage/publish"},
		{http.MethodPost, "/api/v1/admin/homepage/reset-draft"},
		{http.MethodGet, "/api/v1/admin/homepage/preview"},
	}
	for _, test := range tests {
		t.Run(test.method+" "+test.path, func(t *testing.T) {
			recorder := httptest.NewRecorder()
			request := httptest.NewRequest(test.method, test.path, nil)
			engine.ServeHTTP(recorder, request)
			if recorder.Code != http.StatusUnauthorized {
				t.Fatalf("expected status %d, got %d: %s", http.StatusUnauthorized, recorder.Code, recorder.Body.String())
			}
			assertErrorCode(t, recorder.Body.Bytes(), "UNAUTHORIZED")
		})
	}
}

func TestHomepageAdminReadsDraftAndPublishedConfigs(t *testing.T) {
	engine, mock := newHomepageAdminTestRouter(t)
	draftJSON := mustEncodeHomepageConfig(t, heroOnlyHomepageConfig("Draft title"))
	publishedJSON := mustEncodeHomepageConfig(t, heroOnlyHomepageConfig("Published title"))
	expectHomepageVersionQuery(mock, model.HomepageStatusDraft, 4, draftJSON)
	expectHomepageVersionQuery(mock, model.HomepageStatusPublished, 3, publishedJSON)

	draft := performAdminHomepageRequest(t, engine, http.MethodGet, "/api/v1/admin/homepage/draft", nil)
	published := performAdminHomepageRequest(t, engine, http.MethodGet, "/api/v1/admin/homepage/published", nil)

	if draft.Code != http.StatusOK || !strings.Contains(draft.Body.String(), "Draft title") ||
		!strings.Contains(draft.Body.String(), `"status":"draft"`) {
		t.Fatalf("unexpected draft response: %d %s", draft.Code, draft.Body.String())
	}
	if published.Code != http.StatusOK || !strings.Contains(published.Body.String(), "Published title") ||
		!strings.Contains(published.Body.String(), `"status":"published"`) {
		t.Fatalf("unexpected published response: %d %s", published.Code, published.Body.String())
	}
	if draft.Header().Get("Cache-Control") != "no-store" || published.Header().Get("Cache-Control") != "no-store" {
		t.Fatal("administrator config responses must disable caching")
	}
}

func TestSavingHomepageDraftDoesNotChangePublicHomepage(t *testing.T) {
	engine, mock := newHomepageAdminTestRouter(t)
	draftConfig := heroOnlyHomepageConfig("Unsaved public title")
	publishedJSON := mustEncodeHomepageConfig(t, heroOnlyHomepageConfig("Current public title"))
	now := time.Now().UTC()

	mock.ExpectBegin()
	expectLockedHomepageVersionQuery(mock, model.HomepageStatusDraft, 1, 1, publishedJSON, nil, now)
	mock.ExpectExec(`UPDATE "homepage_versions" SET .* WHERE id = \$[0-9]+`).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	body, err := json.Marshal(draftConfig)
	if err != nil {
		t.Fatalf("encode save draft request: %v", err)
	}
	saved := performAdminHomepageRequest(
		t, engine, http.MethodPut, "/api/v1/admin/homepage/draft", bytes.NewReader(body),
	)
	if saved.Code != http.StatusOK || !strings.Contains(saved.Body.String(), "Unsaved public title") {
		t.Fatalf("unexpected save response: %d %s", saved.Code, saved.Body.String())
	}

	expectHomepageVersionQuery(mock, model.HomepageStatusPublished, 1, publishedJSON)
	publicRecorder := httptest.NewRecorder()
	publicRequest := httptest.NewRequest(http.MethodGet, "/api/v1/homepage", nil)
	engine.ServeHTTP(publicRecorder, publicRequest)
	if publicRecorder.Code != http.StatusOK || !strings.Contains(publicRecorder.Body.String(), "Current public title") {
		t.Fatalf("draft affected public homepage: %d %s", publicRecorder.Code, publicRecorder.Body.String())
	}
	if strings.Contains(publicRecorder.Body.String(), "Unsaved public title") {
		t.Fatalf("public homepage exposed draft config: %s", publicRecorder.Body.String())
	}
}

func TestPublishingHomepageDraftIsTransactionalAndImmediatelyPublic(t *testing.T) {
	engine, mock := newHomepageAdminTestRouter(t)
	draftJSON := mustEncodeHomepageConfig(t, heroOnlyHomepageConfig("New published title"))
	publishedJSON := mustEncodeHomepageConfig(t, heroOnlyHomepageConfig("Old published title"))
	now := time.Now().UTC()

	mock.ExpectBegin()
	expectLockedHomepageVersionQuery(mock, model.HomepageStatusDraft, 1, 1, draftJSON, nil, now)
	expectLockedHomepageVersionQuery(mock, model.HomepageStatusPublished, 2, 1, publishedJSON, &now, now)
	mock.ExpectExec(`UPDATE "homepage_versions" SET .* WHERE id = \$[0-9]+`).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectExec(`UPDATE "homepage_versions" SET .* WHERE id = \$[0-9]+`).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	published := performAdminHomepageRequest(t, engine, http.MethodPost, "/api/v1/admin/homepage/publish", nil)
	if published.Code != http.StatusOK {
		t.Fatalf("unexpected publish response: %d %s", published.Code, published.Body.String())
	}
	var publishResponse struct {
		Data service.HomepagePublishResult `json:"data"`
	}
	if err := json.Unmarshal(published.Body.Bytes(), &publishResponse); err != nil {
		t.Fatalf("decode publish response: %v", err)
	}
	if publishResponse.Data.Version != 2 || publishResponse.Data.PublishedAt.IsZero() {
		t.Fatalf("unexpected publish result: %+v", publishResponse.Data)
	}

	expectHomepageVersionQuery(mock, model.HomepageStatusPublished, 2, draftJSON)
	publicRecorder := httptest.NewRecorder()
	publicRequest := httptest.NewRequest(http.MethodGet, "/api/v1/homepage", nil)
	engine.ServeHTTP(publicRecorder, publicRequest)
	if publicRecorder.Code != http.StatusOK || !strings.Contains(publicRecorder.Body.String(), "New published title") {
		t.Fatalf("published config was not immediately public: %d %s", publicRecorder.Code, publicRecorder.Body.String())
	}
}

func TestPublishingHomepageDraftRollsBackOnFailure(t *testing.T) {
	engine, mock := newHomepageAdminTestRouter(t)
	draftJSON := mustEncodeHomepageConfig(t, heroOnlyHomepageConfig("Draft title"))
	publishedJSON := mustEncodeHomepageConfig(t, heroOnlyHomepageConfig("Published title"))
	now := time.Now().UTC()

	mock.ExpectBegin()
	expectLockedHomepageVersionQuery(mock, model.HomepageStatusDraft, 1, 1, draftJSON, nil, now)
	expectLockedHomepageVersionQuery(mock, model.HomepageStatusPublished, 2, 1, publishedJSON, &now, now)
	mock.ExpectExec(`UPDATE "homepage_versions" SET .* WHERE id = \$[0-9]+`).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectExec(`UPDATE "homepage_versions" SET .* WHERE id = \$[0-9]+`).
		WillReturnError(errors.New("database write failed"))
	mock.ExpectRollback()

	recorder := performAdminHomepageRequest(t, engine, http.MethodPost, "/api/v1/admin/homepage/publish", nil)
	if recorder.Code != http.StatusInternalServerError {
		t.Fatalf("expected status %d, got %d: %s", http.StatusInternalServerError, recorder.Code, recorder.Body.String())
	}
	assertErrorCode(t, recorder.Body.Bytes(), "INTERNAL_ERROR")
}

func TestResetAndPreviewHomepageDraft(t *testing.T) {
	engine, mock := newHomepageAdminTestRouter(t)
	draftJSON := mustEncodeHomepageConfig(t, heroOnlyHomepageConfig("Discarded draft"))
	publishedJSON := mustEncodeHomepageConfig(t, heroOnlyHomepageConfig("Restored published title"))
	now := time.Now().UTC()

	mock.ExpectBegin()
	expectLockedHomepageVersionQuery(mock, model.HomepageStatusDraft, 1, 2, draftJSON, nil, now)
	expectLockedHomepageVersionQuery(mock, model.HomepageStatusPublished, 2, 2, publishedJSON, &now, now)
	mock.ExpectExec(`UPDATE "homepage_versions" SET .* WHERE id = \$[0-9]+`).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	reset := performAdminHomepageRequest(t, engine, http.MethodPost, "/api/v1/admin/homepage/reset-draft", nil)
	if reset.Code != http.StatusOK || !strings.Contains(reset.Body.String(), "Restored published title") ||
		strings.Contains(reset.Body.String(), "Discarded draft") {
		t.Fatalf("unexpected reset response: %d %s", reset.Code, reset.Body.String())
	}

	expectHomepageVersionQuery(mock, model.HomepageStatusDraft, 2, publishedJSON)
	preview := performAdminHomepageRequest(t, engine, http.MethodGet, "/api/v1/admin/homepage/preview", nil)
	if preview.Code != http.StatusOK || !strings.Contains(preview.Body.String(), "Restored published title") {
		t.Fatalf("unexpected preview response: %d %s", preview.Code, preview.Body.String())
	}
	if preview.Header().Get("Cache-Control") != "no-store" {
		t.Fatal("administrator homepage preview must disable caching")
	}
}

func TestSaveHomepageDraftRejectsInvalidConfig(t *testing.T) {
	engine, _ := newHomepageAdminTestRouter(t)
	invalidBodies := []string{
		`{"modules":[]}`,
		`{"modules":[{"type":"unknown","enabled":true,"sort_order":1,"config":{}}]}`,
		`{"modules":[],"unexpected":true}`,
	}
	for _, body := range invalidBodies {
		recorder := performAdminHomepageRequest(
			t, engine, http.MethodPut, "/api/v1/admin/homepage/draft", strings.NewReader(body),
		)
		if recorder.Code != http.StatusBadRequest {
			t.Fatalf("expected status %d, got %d: %s", http.StatusBadRequest, recorder.Code, recorder.Body.String())
		}
		assertErrorCode(t, recorder.Body.Bytes(), "VALIDATION_ERROR")
	}
}

func newHomepageAdminTestRouter(t *testing.T) (*gin.Engine, sqlmock.Sqlmock) {
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
		AuthHandler:    handler.NewAuthHandler(authService),
		AuthMiddleware: middleware.Authenticate(authService),
		PublicConfig:   handler.NewPublicConfigHandler(service.NewPublicConfigService(database)),
		HomepageAdmin:  handler.NewHomepageAdminHandler(service.NewHomepageAdminService(database)),
	}), mock
}

func performAdminHomepageRequest(
	t *testing.T,
	engine *gin.Engine,
	method string,
	path string,
	body interface{ Read([]byte) (int, error) },
) *httptest.ResponseRecorder {
	t.Helper()
	request := httptest.NewRequest(method, path, body)
	request.Header.Set("Authorization", administratorBearerToken(t))
	if body != nil {
		request.Header.Set("Content-Type", "application/json")
	}
	recorder := httptest.NewRecorder()
	engine.ServeHTTP(recorder, request)
	return recorder
}

func heroOnlyHomepageConfig(title string) model.HomepageConfig {
	config := model.DefaultHomepageConfig()
	for index := range config.Modules {
		config.Modules[index].Enabled = config.Modules[index].Type == model.HomepageModuleHero
		if config.Modules[index].Type == model.HomepageModuleHero {
			config.Modules[index].Hero.Title = title
		}
	}
	return config
}

func mustEncodeHomepageConfig(t *testing.T, config model.HomepageConfig) model.JSONDocument {
	t.Helper()
	document, err := model.EncodeHomepageConfig(config)
	if err != nil {
		t.Fatalf("encode homepage test config: %v", err)
	}
	return document
}

func expectHomepageVersionQuery(
	mock sqlmock.Sqlmock,
	status string,
	version uint64,
	config model.JSONDocument,
) {
	now := time.Now().UTC()
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "homepage_versions" WHERE status = $1 LIMIT $2`)).
		WithArgs(status, 1).
		WillReturnRows(homepageVersionRows().AddRow(
			uint64(1), status, version, string(config), nullablePublishedAt(status, now), now, now,
		))
}

func expectLockedHomepageVersionQuery(
	mock sqlmock.Sqlmock,
	status string,
	id uint64,
	version uint64,
	config model.JSONDocument,
	publishedAt *time.Time,
	updatedAt time.Time,
) {
	mock.ExpectQuery(`SELECT \* FROM "homepage_versions" WHERE status = \$1 LIMIT \$2 FOR UPDATE`).
		WithArgs(status, 1).
		WillReturnRows(homepageVersionRows().AddRow(
			id, status, version, string(config), publishedAt, updatedAt, updatedAt,
		))
}

func homepageVersionRows() *sqlmock.Rows {
	return sqlmock.NewRows([]string{
		"id", "status", "version_number", "config_json", "published_at", "created_at", "updated_at",
	})
}

func nullablePublishedAt(status string, value time.Time) any {
	if status == model.HomepageStatusPublished {
		return value
	}
	return nil
}
