package router

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"personal-blog/backend/internal/handler"
	"personal-blog/backend/internal/middleware"
	"personal-blog/backend/internal/model"
	"personal-blog/backend/internal/service"
)

func TestCreateDraftArticle(t *testing.T) {
	engine, mock := newArticleTestRouter(t)
	mock.ExpectQuery(`SELECT count\(\*\) FROM "articles" WHERE slug = \$1`).
		WithArgs("hello-world").
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))
	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO "articles" .* RETURNING "id"`).
		WithArgs(
			"Hello World",
			"hello-world",
			"Summary",
			"# Draft content",
			model.ArticleStatusDraft,
			nil,
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
		).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(uint64(1)))
	mock.ExpectCommit()

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(
		http.MethodPost,
		"/api/v1/admin/articles",
		bytes.NewBufferString(`{"title":"Hello World","summary":"Summary","content":"# Draft content"}`),
	)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", administratorBearerToken(t))
	engine.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d: %s", http.StatusCreated, recorder.Code, recorder.Body.String())
	}
	var response struct {
		Data struct {
			ID          uint64  `json:"id"`
			Slug        string  `json:"slug"`
			Status      string  `json:"status"`
			PublishedAt *string `json:"published_at"`
		} `json:"data"`
	}
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("decode create response: %v", err)
	}
	if response.Data.ID != 1 || response.Data.Slug != "hello-world" || response.Data.Status != model.ArticleStatusDraft {
		t.Fatalf("unexpected draft article: %+v", response.Data)
	}
	if response.Data.PublishedAt != nil {
		t.Fatalf("expected draft published_at to be null, got %q", *response.Data.PublishedAt)
	}
}

func TestDraftArticleIsNotPublic(t *testing.T) {
	engine, mock := newArticleTestRouter(t)
	mock.ExpectQuery(`SELECT \* FROM "articles" WHERE slug = \$1 AND status = \$2 ORDER BY "articles"\."id" LIMIT \$3`).
		WithArgs("draft-post", model.ArticleStatusPublished, 1).
		WillReturnRows(articleRows())

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/api/v1/articles/draft-post", nil)
	engine.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d: %s", http.StatusNotFound, recorder.Code, recorder.Body.String())
	}
	assertErrorCode(t, recorder.Body.Bytes(), "ARTICLE_NOT_FOUND")
}

func TestPublishDraftMakesArticlePublic(t *testing.T) {
	engine, mock := newArticleTestRouter(t)
	createdAt := time.Now().UTC().Add(-time.Hour)
	mock.ExpectQuery(`SELECT \* FROM "articles" WHERE "articles"\."id" = \$1 ORDER BY "articles"\."id" LIMIT \$2`).
		WithArgs(uint64(1), 1).
		WillReturnRows(articleRows().AddRow(
			uint64(1), "Draft", "draft-post", "Summary", "# Content",
			model.ArticleStatusDraft, nil, createdAt, createdAt,
		))
	mock.ExpectQuery(`SELECT count\(\*\) FROM "articles" WHERE slug = \$1 AND id <> \$2`).
		WithArgs("draft-post", uint64(1)).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))
	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE "articles" SET .* WHERE id = \$[0-9]+`).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	updateRecorder := httptest.NewRecorder()
	updateRequest := httptest.NewRequest(
		http.MethodPut,
		"/api/v1/admin/articles/1",
		bytes.NewBufferString(`{"title":"Published","slug":"draft-post","summary":"Summary","content":"# Published content","status":"published"}`),
	)
	updateRequest.Header.Set("Content-Type", "application/json")
	updateRequest.Header.Set("Authorization", administratorBearerToken(t))
	engine.ServeHTTP(updateRecorder, updateRequest)

	if updateRecorder.Code != http.StatusOK {
		t.Fatalf("expected update status %d, got %d: %s", http.StatusOK, updateRecorder.Code, updateRecorder.Body.String())
	}
	var updateResponse struct {
		Data struct {
			PublishedAt *time.Time `json:"published_at"`
			Status      string     `json:"status"`
		} `json:"data"`
	}
	if err := json.Unmarshal(updateRecorder.Body.Bytes(), &updateResponse); err != nil {
		t.Fatalf("decode update response: %v", err)
	}
	if updateResponse.Data.PublishedAt == nil || updateResponse.Data.Status != model.ArticleStatusPublished {
		t.Fatalf("article was not published: %+v", updateResponse.Data)
	}

	updatedAt := time.Now().UTC()
	mock.ExpectQuery(`SELECT \* FROM "articles" WHERE slug = \$1 AND status = \$2 ORDER BY "articles"\."id" LIMIT \$3`).
		WithArgs("draft-post", model.ArticleStatusPublished, 1).
		WillReturnRows(articleRows().AddRow(
			uint64(1), "Published", "draft-post", "Summary", "# Published content",
			model.ArticleStatusPublished, *updateResponse.Data.PublishedAt, createdAt, updatedAt,
		))

	publicRecorder := httptest.NewRecorder()
	publicRequest := httptest.NewRequest(http.MethodGet, "/api/v1/articles/draft-post", nil)
	engine.ServeHTTP(publicRecorder, publicRequest)
	if publicRecorder.Code != http.StatusOK {
		t.Fatalf("expected public status %d, got %d: %s", http.StatusOK, publicRecorder.Code, publicRecorder.Body.String())
	}
	if !strings.Contains(publicRecorder.Body.String(), "# Published content") {
		t.Fatalf("public response did not contain article content: %s", publicRecorder.Body.String())
	}
}

func TestCreateArticleRejectsDuplicateSlug(t *testing.T) {
	engine, mock := newArticleTestRouter(t)
	mock.ExpectQuery(`SELECT count\(\*\) FROM "articles" WHERE slug = \$1`).
		WithArgs("duplicate-slug").
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(
		http.MethodPost,
		"/api/v1/admin/articles",
		bytes.NewBufferString(`{"title":"Duplicate","slug":"duplicate-slug","summary":"","content":"Content","status":"draft"}`),
	)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", administratorBearerToken(t))
	engine.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusConflict {
		t.Fatalf("expected status %d, got %d: %s", http.StatusConflict, recorder.Code, recorder.Body.String())
	}
	assertErrorCode(t, recorder.Body.Bytes(), "SLUG_CONFLICT")
}

func TestCreateArticleRejectsInvalidSlug(t *testing.T) {
	engine, _ := newArticleTestRouter(t)
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(
		http.MethodPost,
		"/api/v1/admin/articles",
		bytes.NewBufferString(`{"title":"Invalid Slug","slug":"Invalid_Slug","content":"Content","status":"draft"}`),
	)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", administratorBearerToken(t))
	engine.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d: %s", http.StatusBadRequest, recorder.Code, recorder.Body.String())
	}
	assertErrorCode(t, recorder.Body.Bytes(), "VALIDATION_ERROR")
}

func TestAdministratorArticleRoutesRequireToken(t *testing.T) {
	engine, _ := newArticleTestRouter(t)
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/api/v1/admin/articles", nil)
	engine.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusUnauthorized {
		t.Fatalf("expected status %d, got %d: %s", http.StatusUnauthorized, recorder.Code, recorder.Body.String())
	}
	assertErrorCode(t, recorder.Body.Bytes(), "UNAUTHORIZED")
}

func TestPublishedArticleListIsPaginatedWithoutContent(t *testing.T) {
	engine, mock := newArticleTestRouter(t)
	publishedAt := time.Now().UTC()
	mock.ExpectQuery(`SELECT count\(\*\) FROM "articles" WHERE status = \$1`).
		WithArgs(model.ArticleStatusPublished).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(3))
	mock.ExpectQuery(`SELECT articles\.id, articles\.title, articles\.slug, articles\.summary, articles\.published_at, \(SELECT COUNT\(\*\) FROM comments WHERE comments\.article_id = articles\.id AND comments\.status = \$1\) AS comment_count FROM "articles" WHERE status = \$2 ORDER BY published_at DESC, id DESC LIMIT \$3`).
		WithArgs(model.CommentStatusApproved, model.ArticleStatusPublished, 2).
		WillReturnRows(sqlmock.NewRows([]string{"id", "title", "slug", "summary", "published_at", "comment_count"}).
			AddRow(uint64(2), "Second", "second", "Second summary", publishedAt, int64(2)).
			AddRow(uint64(1), "First", "first", "First summary", publishedAt.Add(-time.Minute), int64(0)))

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/api/v1/articles?page=1&page_size=2", nil)
	engine.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d: %s", http.StatusOK, recorder.Code, recorder.Body.String())
	}
	if strings.Contains(recorder.Body.String(), "content") {
		t.Fatalf("article list exposed full content: %s", recorder.Body.String())
	}
	var response struct {
		Data struct {
			Items []struct {
				CommentCount int64 `json:"comment_count"`
			} `json:"items"`
			Page       int   `json:"page"`
			PageSize   int   `json:"page_size"`
			Total      int64 `json:"total"`
			TotalPages int   `json:"total_pages"`
		} `json:"data"`
	}
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("decode list response: %v", err)
	}
	if len(response.Data.Items) != 2 || response.Data.Page != 1 || response.Data.PageSize != 2 ||
		response.Data.Total != 3 || response.Data.TotalPages != 2 {
		t.Fatalf("unexpected pagination response: %+v", response.Data)
	}
	if response.Data.Items[0].CommentCount != 2 || response.Data.Items[1].CommentCount != 0 {
		t.Fatalf("unexpected approved comment counts: %+v", response.Data.Items)
	}
}

func TestDeleteArticle(t *testing.T) {
	engine, mock := newArticleTestRouter(t)
	mock.ExpectBegin()
	mock.ExpectExec(`DELETE FROM "articles" WHERE "articles"\."id" = \$1`).
		WithArgs(uint64(1)).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodDelete, "/api/v1/admin/articles/1", nil)
	request.Header.Set("Authorization", administratorBearerToken(t))
	engine.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusNoContent {
		t.Fatalf("expected status %d, got %d: %s", http.StatusNoContent, recorder.Code, recorder.Body.String())
	}
}

func newArticleTestRouter(t *testing.T) (*gin.Engine, sqlmock.Sqlmock) {
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
	authHandler := handler.NewAuthHandler(authService)
	articleHandler := handler.NewArticleHandler(service.NewArticleService(database))
	return New(Dependencies{
		AuthHandler:    authHandler,
		ArticleHandler: articleHandler,
		AuthMiddleware: middleware.Authenticate(authService),
	}), mock
}

func administratorBearerToken(t *testing.T) string {
	t.Helper()
	now := time.Now().UTC()
	claims := service.AdminClaims{
		Username: "admin",
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   "1",
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Hour)),
		},
	}
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(testJWTSecret))
	if err != nil {
		t.Fatalf("sign administrator test token: %v", err)
	}
	return "Bearer " + token
}

func articleRows() *sqlmock.Rows {
	return sqlmock.NewRows([]string{
		"id", "title", "slug", "summary", "content", "status", "published_at", "created_at", "updated_at",
	})
}
