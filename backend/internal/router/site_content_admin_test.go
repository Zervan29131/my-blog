package router

import (
	"bytes"
	"net/http"
	"net/http/httptest"
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

func TestSiteContentAdminRoutesRequireAuthentication(t *testing.T) {
	engine, _ := newSiteContentAdminTestRouter(t)
	routes := []struct {
		method string
		path   string
	}{
		{http.MethodGet, "/api/v1/admin/navigation"},
		{http.MethodPost, "/api/v1/admin/navigation"},
		{http.MethodPut, "/api/v1/admin/navigation/order"},
		{http.MethodPut, "/api/v1/admin/navigation/1"},
		{http.MethodDelete, "/api/v1/admin/navigation/1"},
		{http.MethodGet, "/api/v1/admin/social-links"},
		{http.MethodPost, "/api/v1/admin/social-links"},
		{http.MethodPut, "/api/v1/admin/social-links/order"},
		{http.MethodPut, "/api/v1/admin/social-links/1"},
		{http.MethodDelete, "/api/v1/admin/social-links/1"},
		{http.MethodGet, "/api/v1/admin/featured-articles"},
		{http.MethodPost, "/api/v1/admin/featured-articles"},
		{http.MethodPut, "/api/v1/admin/featured-articles/order"},
		{http.MethodPut, "/api/v1/admin/featured-articles/1"},
		{http.MethodDelete, "/api/v1/admin/featured-articles/1"},
	}
	for _, route := range routes {
		t.Run(route.method+" "+route.path, func(t *testing.T) {
			recorder := performSiteContentRequest(t, engine, route.method, route.path, nil, false)
			if recorder.Code != http.StatusUnauthorized {
				t.Fatalf("expected status %d, got %d: %s", http.StatusUnauthorized, recorder.Code, recorder.Body.String())
			}
			assertErrorCode(t, recorder.Body.Bytes(), "UNAUTHORIZED")
		})
	}
}

func TestNavigationCRUDAndReorder(t *testing.T) {
	engine, mock := newSiteContentAdminTestRouter(t)
	now := time.Now().UTC()
	mock.ExpectQuery(`SELECT \* FROM "navigation_items" ORDER BY sort_order ASC, id ASC`).
		WillReturnRows(navigationRows().AddRow(
			uint64(1), "Home", "/", model.LinkTypeInternal, false, true, 10, now, now,
		))
	list := performSiteContentRequest(t, engine, http.MethodGet, "/api/v1/admin/navigation", nil, true)
	if list.Code != http.StatusOK || !bytes.Contains(list.Body.Bytes(), []byte(`"name":"Home"`)) {
		t.Fatalf("unexpected navigation list: %d %s", list.Code, list.Body.String())
	}

	mock.ExpectBegin()
	mock.ExpectQuery(`SELECT count\(\*\) FROM "navigation_items"`).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))
	mock.ExpectQuery(`INSERT INTO "navigation_items" .* RETURNING "id"`).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(uint64(2)))
	mock.ExpectCommit()
	created := performSiteContentRequest(t, engine, http.MethodPost, "/api/v1/admin/navigation", []byte(
		`{"name":"Docs","url":"https://example.com/docs","link_type":"external","open_in_new_tab":true,"is_visible":true,"sort_order":20}`,
	), true)
	if created.Code != http.StatusCreated || !bytes.Contains(created.Body.Bytes(), []byte(`"id":2`)) {
		t.Fatalf("unexpected navigation create response: %d %s", created.Code, created.Body.String())
	}

	mock.ExpectQuery(`SELECT \* FROM "navigation_items" WHERE id = \$1 LIMIT \$2`).
		WithArgs(uint64(2), 1).
		WillReturnRows(navigationRows().AddRow(
			uint64(2), "Docs", "https://example.com/docs", model.LinkTypeExternal, true, true, 20, now, now,
		))
	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE "navigation_items" SET .* WHERE id = \$[0-9]+`).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()
	updated := performSiteContentRequest(t, engine, http.MethodPut, "/api/v1/admin/navigation/2", []byte(
		`{"name":"Documentation","url":"/archive","link_type":"internal","open_in_new_tab":false,"is_visible":false,"sort_order":30}`,
	), true)
	if updated.Code != http.StatusOK || !bytes.Contains(updated.Body.Bytes(), []byte(`"is_visible":false`)) {
		t.Fatalf("unexpected navigation update response: %d %s", updated.Code, updated.Body.String())
	}

	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE "navigation_items" SET .* WHERE id = \$[0-9]+`).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectExec(`UPDATE "navigation_items" SET .* WHERE id = \$[0-9]+`).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()
	reordered := performSiteContentRequest(t, engine, http.MethodPut, "/api/v1/admin/navigation/order", []byte(
		`{"items":[{"id":2,"sort_order":10},{"id":1,"sort_order":20}]}`,
	), true)
	if reordered.Code != http.StatusOK || !bytes.Contains(reordered.Body.Bytes(), []byte(`"updated":2`)) {
		t.Fatalf("unexpected navigation reorder response: %d %s", reordered.Code, reordered.Body.String())
	}

	mock.ExpectBegin()
	mock.ExpectExec(`DELETE FROM "navigation_items" WHERE "navigation_items"\."id" = \$1`).
		WithArgs(uint64(2)).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()
	deleted := performSiteContentRequest(t, engine, http.MethodDelete, "/api/v1/admin/navigation/2", nil, true)
	if deleted.Code != http.StatusNoContent {
		t.Fatalf("unexpected navigation delete response: %d %s", deleted.Code, deleted.Body.String())
	}
}

func TestNavigationRejectsDangerousURLsAndUnknownFields(t *testing.T) {
	engine, _ := newSiteContentAdminTestRouter(t)
	bodies := []string{
		`{"name":"Bad","url":"javascript:alert(1)","link_type":"external","open_in_new_tab":true,"is_visible":true,"sort_order":10}`,
		`{"name":"Bad","url":"//evil.example","link_type":"internal","open_in_new_tab":false,"is_visible":true,"sort_order":10}`,
		`{"name":"Bad","url":"/","link_type":"internal","open_in_new_tab":false,"is_visible":true,"sort_order":10,"html":"<script>"}`,
	}
	for _, body := range bodies {
		recorder := performSiteContentRequest(
			t, engine, http.MethodPost, "/api/v1/admin/navigation", []byte(body), true,
		)
		if recorder.Code != http.StatusBadRequest {
			t.Fatalf("expected status %d, got %d: %s", http.StatusBadRequest, recorder.Code, recorder.Body.String())
		}
		assertErrorCode(t, recorder.Body.Bytes(), "VALIDATION_ERROR")
	}
}

func TestNavigationReorderRollsBackAndResourceLimitsAreEnforced(t *testing.T) {
	engine, mock := newSiteContentAdminTestRouter(t)
	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE "navigation_items" SET .* WHERE id = \$[0-9]+`).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectExec(`UPDATE "navigation_items" SET .* WHERE id = \$[0-9]+`).
		WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectRollback()
	reorder := performSiteContentRequest(t, engine, http.MethodPut, "/api/v1/admin/navigation/order", []byte(
		`{"items":[{"id":1,"sort_order":10},{"id":999,"sort_order":20}]}`,
	), true)
	if reorder.Code != http.StatusNotFound {
		t.Fatalf("expected transactional reorder failure, got %d: %s", reorder.Code, reorder.Body.String())
	}

	mock.ExpectBegin()
	mock.ExpectQuery(`SELECT count\(\*\) FROM "navigation_items"`).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(10))
	mock.ExpectRollback()
	limit := performSiteContentRequest(t, engine, http.MethodPost, "/api/v1/admin/navigation", []byte(
		`{"name":"Extra","url":"/extra","link_type":"internal","open_in_new_tab":false,"is_visible":true,"sort_order":100}`,
	), true)
	if limit.Code != http.StatusConflict {
		t.Fatalf("expected navigation limit conflict, got %d: %s", limit.Code, limit.Body.String())
	}
	assertErrorCode(t, limit.Body.Bytes(), "RESOURCE_LIMIT_EXCEEDED")

	mock.ExpectBegin()
	mock.ExpectQuery(`SELECT count\(\*\) FROM "social_links"`).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(15))
	mock.ExpectRollback()
	socialLimit := performSiteContentRequest(t, engine, http.MethodPost, "/api/v1/admin/social-links", []byte(
		`{"platform":"custom","display_name":"Extra","url":"https://example.com","is_visible":true,"sort_order":100}`,
	), true)
	if socialLimit.Code != http.StatusConflict {
		t.Fatalf("expected social link limit conflict, got %d: %s", socialLimit.Code, socialLimit.Body.String())
	}
}

func TestSocialLinkCRUDAndURLValidation(t *testing.T) {
	engine, mock := newSiteContentAdminTestRouter(t)
	now := time.Now().UTC()
	mock.ExpectQuery(`SELECT \* FROM "social_links" ORDER BY sort_order ASC, id ASC`).
		WillReturnRows(socialLinkRows().AddRow(
			uint64(1), model.SocialPlatformGitHub, "GitHub", "https://github.com/example", true, 10, now, now,
		))
	list := performSiteContentRequest(t, engine, http.MethodGet, "/api/v1/admin/social-links", nil, true)
	if list.Code != http.StatusOK || !bytes.Contains(list.Body.Bytes(), []byte(`"platform":"github"`)) {
		t.Fatalf("unexpected social link list: %d %s", list.Code, list.Body.String())
	}

	mock.ExpectBegin()
	mock.ExpectQuery(`SELECT count\(\*\) FROM "social_links"`).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))
	mock.ExpectQuery(`INSERT INTO "social_links" .* RETURNING "id"`).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(uint64(2)))
	mock.ExpectCommit()
	created := performSiteContentRequest(t, engine, http.MethodPost, "/api/v1/admin/social-links", []byte(
		`{"platform":"email","display_name":"Email","url":"mailto:hello@example.com","is_visible":true,"sort_order":20}`,
	), true)
	if created.Code != http.StatusCreated || !bytes.Contains(created.Body.Bytes(), []byte(`"id":2`)) {
		t.Fatalf("unexpected social link create response: %d %s", created.Code, created.Body.String())
	}

	mock.ExpectQuery(`SELECT \* FROM "social_links" WHERE id = \$1 LIMIT \$2`).
		WithArgs(uint64(2), 1).
		WillReturnRows(socialLinkRows().AddRow(
			uint64(2), model.SocialPlatformEmail, "Email", "mailto:hello@example.com", true, 20, now, now,
		))
	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE "social_links" SET .* WHERE id = \$[0-9]+`).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()
	updated := performSiteContentRequest(t, engine, http.MethodPut, "/api/v1/admin/social-links/2", []byte(
		`{"platform":"custom","display_name":"Profile","url":"https://example.com/profile","is_visible":false,"sort_order":30}`,
	), true)
	if updated.Code != http.StatusOK || !bytes.Contains(updated.Body.Bytes(), []byte(`"is_visible":false`)) {
		t.Fatalf("unexpected social link update response: %d %s", updated.Code, updated.Body.String())
	}

	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE "social_links" SET .* WHERE id = \$[0-9]+`).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()
	reordered := performSiteContentRequest(t, engine, http.MethodPut, "/api/v1/admin/social-links/order", []byte(
		`{"items":[{"id":2,"sort_order":5}]}`,
	), true)
	if reordered.Code != http.StatusOK {
		t.Fatalf("unexpected social link reorder response: %d %s", reordered.Code, reordered.Body.String())
	}

	mock.ExpectBegin()
	mock.ExpectExec(`DELETE FROM "social_links" WHERE "social_links"\."id" = \$1`).
		WithArgs(uint64(2)).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()
	deleted := performSiteContentRequest(t, engine, http.MethodDelete, "/api/v1/admin/social-links/2", nil, true)
	if deleted.Code != http.StatusNoContent {
		t.Fatalf("unexpected social link delete response: %d %s", deleted.Code, deleted.Body.String())
	}

	unsafe := performSiteContentRequest(t, engine, http.MethodPost, "/api/v1/admin/social-links", []byte(
		`{"platform":"custom","display_name":"Bad","url":"javascript:alert(1)","is_visible":true,"sort_order":1}`,
	), true)
	if unsafe.Code != http.StatusBadRequest {
		t.Fatalf("expected dangerous URL rejection, got %d: %s", unsafe.Code, unsafe.Body.String())
	}
}

func TestFeaturedArticleManagement(t *testing.T) {
	engine, mock := newSiteContentAdminTestRouter(t)
	now := time.Now().UTC()
	mock.ExpectQuery(`SELECT featured_articles\..* FROM "featured_articles" JOIN articles ON articles\.id = featured_articles\.article_id ORDER BY featured_articles\.sort_order ASC, featured_articles\.id ASC`).
		WillReturnRows(sqlmock.NewRows([]string{
			"article_id", "title", "slug", "summary", "status", "published_at", "sort_order", "is_visible", "created_at", "updated_at",
		}).AddRow(uint64(7), "Published", "published", "Summary", model.ArticleStatusPublished, now, 10, true, now, now))
	list := performSiteContentRequest(t, engine, http.MethodGet, "/api/v1/admin/featured-articles", nil, true)
	if list.Code != http.StatusOK || !bytes.Contains(list.Body.Bytes(), []byte(`"article_id":7`)) {
		t.Fatalf("unexpected featured article list: %d %s", list.Code, list.Body.String())
	}

	mock.ExpectBegin()
	mock.ExpectQuery(`SELECT \* FROM "articles" WHERE id = \$1 LIMIT \$2`).
		WithArgs(uint64(7), 1).
		WillReturnRows(articleRows().AddRow(
			uint64(7), "Published", "published", "Summary", "Content", model.ArticleStatusPublished, now, now, now,
		))
	mock.ExpectQuery(`SELECT count\(\*\) FROM "featured_articles" WHERE article_id = \$1`).
		WithArgs(uint64(7)).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))
	mock.ExpectQuery(`SELECT count\(\*\) FROM "featured_articles"`).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))
	mock.ExpectQuery(`INSERT INTO "featured_articles" .* RETURNING "id"`).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(uint64(1)))
	mock.ExpectCommit()
	created := performSiteContentRequest(t, engine, http.MethodPost, "/api/v1/admin/featured-articles", []byte(
		`{"article_id":7,"sort_order":10,"is_visible":true}`,
	), true)
	if created.Code != http.StatusCreated || !bytes.Contains(created.Body.Bytes(), []byte(`"article_id":7`)) {
		t.Fatalf("unexpected featured article create response: %d %s", created.Code, created.Body.String())
	}

	mock.ExpectQuery(`SELECT \* FROM "featured_articles" WHERE article_id = \$1 LIMIT \$2`).
		WithArgs(uint64(7), 1).
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "article_id", "sort_order", "is_visible", "created_at", "updated_at",
		}).AddRow(uint64(1), uint64(7), 10, true, now, now))
	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE "featured_articles" SET .* WHERE article_id = \$[0-9]+`).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()
	visibility := performSiteContentRequest(t, engine, http.MethodPut, "/api/v1/admin/featured-articles/7", []byte(
		`{"is_visible":false}`,
	), true)
	if visibility.Code != http.StatusOK || !bytes.Contains(visibility.Body.Bytes(), []byte(`"is_visible":false`)) {
		t.Fatalf("unexpected featured visibility response: %d %s", visibility.Code, visibility.Body.String())
	}

	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE "featured_articles" SET .* WHERE article_id = \$[0-9]+`).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()
	reordered := performSiteContentRequest(t, engine, http.MethodPut, "/api/v1/admin/featured-articles/order", []byte(
		`{"items":[{"article_id":7,"sort_order":20}]}`,
	), true)
	if reordered.Code != http.StatusOK {
		t.Fatalf("unexpected featured reorder response: %d %s", reordered.Code, reordered.Body.String())
	}

	mock.ExpectBegin()
	mock.ExpectExec(`DELETE FROM "featured_articles" WHERE article_id = \$1`).
		WithArgs(uint64(7)).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()
	deleted := performSiteContentRequest(t, engine, http.MethodDelete, "/api/v1/admin/featured-articles/7", nil, true)
	if deleted.Code != http.StatusNoContent {
		t.Fatalf("unexpected featured delete response: %d %s", deleted.Code, deleted.Body.String())
	}
}

func TestFeaturedArticleRejectsDraftAndDuplicateArticles(t *testing.T) {
	engine, mock := newSiteContentAdminTestRouter(t)
	now := time.Now().UTC()
	mock.ExpectBegin()
	mock.ExpectQuery(`SELECT \* FROM "articles" WHERE id = \$1 LIMIT \$2`).
		WithArgs(uint64(8), 1).
		WillReturnRows(articleRows().AddRow(
			uint64(8), "Draft", "draft", "Summary", "Content", model.ArticleStatusDraft, nil, now, now,
		))
	mock.ExpectRollback()
	draft := performSiteContentRequest(t, engine, http.MethodPost, "/api/v1/admin/featured-articles", []byte(
		`{"article_id":8,"sort_order":10,"is_visible":true}`,
	), true)
	if draft.Code != http.StatusBadRequest {
		t.Fatalf("expected draft rejection, got %d: %s", draft.Code, draft.Body.String())
	}
	assertErrorCode(t, draft.Body.Bytes(), "ARTICLE_NOT_PUBLISHED")

	mock.ExpectBegin()
	mock.ExpectQuery(`SELECT \* FROM "articles" WHERE id = \$1 LIMIT \$2`).
		WithArgs(uint64(7), 1).
		WillReturnRows(articleRows().AddRow(
			uint64(7), "Published", "published", "Summary", "Content", model.ArticleStatusPublished, now, now, now,
		))
	mock.ExpectQuery(`SELECT count\(\*\) FROM "featured_articles" WHERE article_id = \$1`).
		WithArgs(uint64(7)).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))
	mock.ExpectRollback()
	duplicate := performSiteContentRequest(t, engine, http.MethodPost, "/api/v1/admin/featured-articles", []byte(
		`{"article_id":7,"sort_order":10,"is_visible":true}`,
	), true)
	if duplicate.Code != http.StatusConflict {
		t.Fatalf("expected duplicate conflict, got %d: %s", duplicate.Code, duplicate.Body.String())
	}
	assertErrorCode(t, duplicate.Body.Bytes(), "RESOURCE_CONFLICT")
}

func TestAdministratorCanFilterPublishedArticleCandidates(t *testing.T) {
	engine, mock := newSiteContentAdminTestRouter(t)
	now := time.Now().UTC()
	mock.ExpectQuery(`SELECT count\(\*\) FROM "articles" WHERE status = \$1`).
		WithArgs(model.ArticleStatusPublished).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))
	mock.ExpectQuery(`SELECT "id","title","slug","summary","status","published_at","created_at","updated_at" FROM "articles" WHERE status = \$1 ORDER BY created_at DESC, id DESC LIMIT \$2`).
		WithArgs(model.ArticleStatusPublished, 10).
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "title", "slug", "summary", "status", "published_at", "created_at", "updated_at",
		}).AddRow(uint64(7), "Published", "published", "Summary", model.ArticleStatusPublished, now, now, now))

	recorder := performSiteContentRequest(
		t, engine, http.MethodGet, "/api/v1/admin/articles?status=published", nil, true,
	)
	if recorder.Code != http.StatusOK || !bytes.Contains(recorder.Body.Bytes(), []byte(`"status":"published"`)) {
		t.Fatalf("unexpected published candidate list: %d %s", recorder.Code, recorder.Body.String())
	}

	invalid := performSiteContentRequest(
		t, engine, http.MethodGet, "/api/v1/admin/articles?status=deleted", nil, true,
	)
	if invalid.Code != http.StatusBadRequest {
		t.Fatalf("expected invalid status rejection, got %d: %s", invalid.Code, invalid.Body.String())
	}
}

func newSiteContentAdminTestRouter(t *testing.T) (*gin.Engine, sqlmock.Sqlmock) {
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
		TranslateError:       true,
	})
	if err != nil {
		t.Fatalf("open GORM test database: %v", err)
	}
	authService := service.NewAuthService(database, testJWTSecret, 24)
	return New(Dependencies{
		AuthHandler:    handler.NewAuthHandler(authService),
		AuthMiddleware: middleware.Authenticate(authService),
		ArticleHandler: handler.NewArticleHandler(service.NewArticleService(database)),
		SiteContentAdmin: handler.NewSiteContentAdminHandler(
			service.NewNavigationAdminService(database),
			service.NewSocialLinkAdminService(database),
			service.NewFeaturedArticleAdminService(database),
		),
	}), mock
}

func performSiteContentRequest(
	t *testing.T,
	engine *gin.Engine,
	method string,
	path string,
	body []byte,
	authenticated bool,
) *httptest.ResponseRecorder {
	t.Helper()
	var request *http.Request
	if body == nil {
		request = httptest.NewRequest(method, path, nil)
	} else {
		request = httptest.NewRequest(method, path, bytes.NewReader(body))
		request.Header.Set("Content-Type", "application/json")
	}
	if authenticated {
		request.Header.Set("Authorization", administratorBearerToken(t))
	}
	recorder := httptest.NewRecorder()
	engine.ServeHTTP(recorder, request)
	return recorder
}

func navigationRows() *sqlmock.Rows {
	return sqlmock.NewRows([]string{
		"id", "name", "url", "link_type", "open_in_new_tab", "is_visible", "sort_order", "created_at", "updated_at",
	})
}

func socialLinkRows() *sqlmock.Rows {
	return sqlmock.NewRows([]string{
		"id", "platform", "display_name", "url", "is_visible", "sort_order", "created_at", "updated_at",
	})
}
