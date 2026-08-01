package router

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"personal-blog/backend/internal/handler"
	"personal-blog/backend/internal/middleware"
	"personal-blog/backend/internal/model"
	"personal-blog/backend/internal/service"
)

func TestDashboardReturnsStatistics(t *testing.T) {
	engine, mock := newDashboardTestRouter(t)
	mock.ExpectQuery(`SELECT\s+\(SELECT COUNT\(\*\) FROM articles\) AS article_total,\s+\(SELECT COUNT\(\*\) FROM articles WHERE status = \$1\) AS article_published,\s+\(SELECT COUNT\(\*\) FROM articles WHERE status = \$2\) AS article_draft,\s+\(SELECT COUNT\(\*\) FROM comments WHERE status = \$3\) AS comment_pending,\s+\(SELECT COUNT\(\*\) FROM comments WHERE status = \$4\) AS comment_approved`).
		WithArgs(
			model.ArticleStatusPublished,
			model.ArticleStatusDraft,
			model.CommentStatusPending,
			model.CommentStatusApproved,
		).
		WillReturnRows(sqlmock.NewRows([]string{
			"article_total", "article_published", "article_draft", "comment_pending", "comment_approved",
		}).AddRow(8, 5, 3, 4, 9))

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/api/v1/admin/dashboard", nil)
	request.Header.Set("Authorization", administratorBearerToken(t))
	engine.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d: %s", http.StatusOK, recorder.Code, recorder.Body.String())
	}
	var response struct {
		Data service.DashboardStats `json:"data"`
	}
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("decode dashboard response: %v", err)
	}
	if response.Data.ArticleTotal != 8 || response.Data.ArticlePublished != 5 ||
		response.Data.ArticleDraft != 3 || response.Data.CommentPending != 4 ||
		response.Data.CommentApproved != 9 {
		t.Fatalf("unexpected dashboard statistics: %+v", response.Data)
	}
}

func TestDashboardRequiresAdministratorToken(t *testing.T) {
	engine, _ := newDashboardTestRouter(t)
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/api/v1/admin/dashboard", nil)
	engine.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusUnauthorized {
		t.Fatalf("expected status %d, got %d: %s", http.StatusUnauthorized, recorder.Code, recorder.Body.String())
	}
	assertErrorCode(t, recorder.Body.Bytes(), "UNAUTHORIZED")
}

func newDashboardTestRouter(t *testing.T) (*gin.Engine, sqlmock.Sqlmock) {
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
		AuthHandler:      handler.NewAuthHandler(authService),
		DashboardHandler: handler.NewDashboardHandler(service.NewDashboardService(database)),
		AuthMiddleware:   middleware.Authenticate(authService),
	}), mock
}
