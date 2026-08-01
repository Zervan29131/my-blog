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
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"personal-blog/backend/internal/handler"
	"personal-blog/backend/internal/middleware"
	"personal-blog/backend/internal/model"
	"personal-blog/backend/internal/service"
)

func TestSubmitCommentAlwaysCreatesPendingWithoutExposingEmail(t *testing.T) {
	engine, mock := newCommentTestRouter(t, false)
	expectPublishedArticleForComments(mock, "published-post")
	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO "comments" \("article_id","nickname","email","content","status","created_at","updated_at"\) VALUES \(\$1,\$2,\$3,\$4,\$5,\$6,\$7\) RETURNING "id"`).
		WithArgs(
			uint64(1), "访客甲", "visitor@example.com", "这是一条待审核评论。",
			model.CommentStatusPending, sqlmock.AnyArg(), sqlmock.AnyArg(),
		).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(uint64(9)))
	mock.ExpectCommit()

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(
		http.MethodPost,
		"/api/v1/articles/published-post/comments",
		bytes.NewBufferString(`{"nickname":"访客甲","email":"visitor@example.com","content":"这是一条待审核评论。"}`),
	)
	request.Header.Set("Content-Type", "application/json")
	engine.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d: %s", http.StatusCreated, recorder.Code, recorder.Body.String())
	}
	if strings.Contains(recorder.Body.String(), "email") || strings.Contains(recorder.Body.String(), "visitor@example.com") {
		t.Fatalf("public submit response exposed email: %s", recorder.Body.String())
	}
	var response struct {
		Data struct {
			ID     uint64 `json:"id"`
			Status string `json:"status"`
		} `json:"data"`
		Message string `json:"message"`
	}
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("decode submit response: %v", err)
	}
	if response.Data.ID != 9 || response.Data.Status != model.CommentStatusPending {
		t.Fatalf("unexpected submitted comment: %+v", response.Data)
	}
	if response.Message != "评论已提交，审核通过后将会显示。" {
		t.Fatalf("unexpected submit message: %q", response.Message)
	}
}

func TestPendingCommentIsNotPublic(t *testing.T) {
	engine, mock := newCommentTestRouter(t, false)
	expectPublishedArticleForComments(mock, "published-post")
	mock.ExpectQuery(`SELECT count\(\*\) FROM "comments" WHERE article_id = \$1 AND status = \$2`).
		WithArgs(uint64(1), model.CommentStatusApproved).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))
	mock.ExpectQuery(`SELECT "id","nickname","content","created_at" FROM "comments" WHERE article_id = \$1 AND status = \$2 ORDER BY created_at ASC, id ASC LIMIT \$3`).
		WithArgs(uint64(1), model.CommentStatusApproved, 20).
		WillReturnRows(publicCommentRows())

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/api/v1/articles/published-post/comments", nil)
	engine.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d: %s", http.StatusOK, recorder.Code, recorder.Body.String())
	}
	var response struct {
		Data struct {
			Items []json.RawMessage `json:"items"`
			Total int64             `json:"total"`
		} `json:"data"`
	}
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("decode public comments: %v", err)
	}
	if len(response.Data.Items) != 0 || response.Data.Total != 0 {
		t.Fatalf("pending comment became public: %+v", response.Data)
	}
}

func TestApprovedCommentIsPublicWithoutEmail(t *testing.T) {
	engine, mock := newCommentTestRouter(t, false)
	createdAt := time.Now().UTC()
	expectPublishedArticleForComments(mock, "published-post")
	mock.ExpectQuery(`SELECT count\(\*\) FROM "comments" WHERE article_id = \$1 AND status = \$2`).
		WithArgs(uint64(1), model.CommentStatusApproved).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))
	mock.ExpectQuery(`SELECT "id","nickname","content","created_at" FROM "comments" WHERE article_id = \$1 AND status = \$2 ORDER BY created_at ASC, id ASC LIMIT \$3`).
		WithArgs(uint64(1), model.CommentStatusApproved, 20).
		WillReturnRows(publicCommentRows().AddRow(uint64(3), "访客乙", "审核通过的评论", createdAt))

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/api/v1/articles/published-post/comments", nil)
	engine.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d: %s", http.StatusOK, recorder.Code, recorder.Body.String())
	}
	if strings.Contains(recorder.Body.String(), "email") || strings.Contains(recorder.Body.String(), "@") {
		t.Fatalf("public comments response exposed email data: %s", recorder.Body.String())
	}
	var response struct {
		Data struct {
			Items []struct {
				ID       uint64 `json:"id"`
				Nickname string `json:"nickname"`
			} `json:"items"`
			Total int64 `json:"total"`
		} `json:"data"`
	}
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("decode public comments: %v", err)
	}
	if response.Data.Total != 1 || len(response.Data.Items) != 1 || response.Data.Items[0].ID != 3 {
		t.Fatalf("approved comment was not public: %+v", response.Data)
	}
}

func TestAdministratorCanApproveAndDeleteComment(t *testing.T) {
	engine, mock := newCommentTestRouter(t, false)
	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE "comments" SET .* WHERE id = \$[0-9]+`).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	updateRecorder := httptest.NewRecorder()
	updateRequest := httptest.NewRequest(
		http.MethodPut,
		"/api/v1/admin/comments/3/status",
		bytes.NewBufferString(`{"status":"approved"}`),
	)
	updateRequest.Header.Set("Content-Type", "application/json")
	updateRequest.Header.Set("Authorization", administratorBearerToken(t))
	engine.ServeHTTP(updateRecorder, updateRequest)
	if updateRecorder.Code != http.StatusOK {
		t.Fatalf("expected status update %d, got %d: %s", http.StatusOK, updateRecorder.Code, updateRecorder.Body.String())
	}

	mock.ExpectBegin()
	mock.ExpectExec(`DELETE FROM "comments" WHERE "comments"\."id" = \$1`).
		WithArgs(uint64(3)).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	deleteRecorder := httptest.NewRecorder()
	deleteRequest := httptest.NewRequest(http.MethodDelete, "/api/v1/admin/comments/3", nil)
	deleteRequest.Header.Set("Authorization", administratorBearerToken(t))
	engine.ServeHTTP(deleteRecorder, deleteRequest)
	if deleteRecorder.Code != http.StatusNoContent {
		t.Fatalf("expected delete status %d, got %d: %s", http.StatusNoContent, deleteRecorder.Code, deleteRecorder.Body.String())
	}
}

func TestAdministratorCanFilterCommentsAndSeeArticleAssociation(t *testing.T) {
	engine, mock := newCommentTestRouter(t, false)
	now := time.Now().UTC()
	mock.ExpectQuery(`SELECT count\(\*\) FROM "comments" WHERE comments\.status = \$1`).
		WithArgs(model.CommentStatusPending).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))
	mock.ExpectQuery(`SELECT comments\.id, comments\.article_id, articles\.title AS article_title, articles\.slug AS article_slug, comments\.nickname, comments\.email, comments\.content, comments\.status, comments\.created_at, comments\.updated_at FROM "comments" JOIN articles ON articles\.id = comments\.article_id WHERE comments\.status = \$1 ORDER BY comments\.created_at DESC, comments\.id DESC LIMIT \$2`).
		WithArgs(model.CommentStatusPending, 20).
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "article_id", "article_title", "article_slug", "nickname", "email",
			"content", "status", "created_at", "updated_at",
		}).AddRow(
			uint64(4), uint64(1), "Published", "published-post", "访客丁", "guest@example.com",
			"等待管理员审核", model.CommentStatusPending, now, now,
		))

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/api/v1/admin/comments?status=pending", nil)
	request.Header.Set("Authorization", administratorBearerToken(t))
	engine.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d: %s", http.StatusOK, recorder.Code, recorder.Body.String())
	}
	var response struct {
		Data struct {
			Items []struct {
				ID      uint64 `json:"id"`
				Email   string `json:"email"`
				Status  string `json:"status"`
				Article struct {
					ID   uint64 `json:"id"`
					Slug string `json:"slug"`
				} `json:"article"`
			} `json:"items"`
		} `json:"data"`
	}
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("decode administrator comments: %v", err)
	}
	if len(response.Data.Items) != 1 || response.Data.Items[0].ID != 4 ||
		response.Data.Items[0].Status != model.CommentStatusPending ||
		response.Data.Items[0].Email != "guest@example.com" ||
		response.Data.Items[0].Article.ID != 1 ||
		response.Data.Items[0].Article.Slug != "published-post" {
		t.Fatalf("unexpected administrator comment response: %+v", response.Data.Items)
	}
}

func TestCommentSubmissionRateLimitIsFivePerMinutePerIP(t *testing.T) {
	engine, _ := newCommentTestRouter(t, true)
	for attempt := 1; attempt <= 6; attempt++ {
		recorder := httptest.NewRecorder()
		request := httptest.NewRequest(
			http.MethodPost,
			"/api/v1/articles/published-post/comments",
			bytes.NewBufferString(`{}`),
		)
		request.Header.Set("Content-Type", "application/json")
		request.RemoteAddr = "198.51.100.10:41000"
		engine.ServeHTTP(recorder, request)

		if attempt <= 5 && recorder.Code != http.StatusBadRequest {
			t.Fatalf("attempt %d should reach validation, got %d: %s", attempt, recorder.Code, recorder.Body.String())
		}
		if attempt == 6 {
			if recorder.Code != http.StatusTooManyRequests {
				t.Fatalf("sixth attempt should be rate limited, got %d: %s", recorder.Code, recorder.Body.String())
			}
			assertErrorCode(t, recorder.Body.Bytes(), "RATE_LIMITED")
		}
	}
}

func TestCommentSubmitRejectsClientProvidedStatus(t *testing.T) {
	engine, _ := newCommentTestRouter(t, false)
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(
		http.MethodPost,
		"/api/v1/articles/published-post/comments",
		bytes.NewBufferString(`{"nickname":"访客丙","content":"试图绕过审核","status":"approved"}`),
	)
	request.Header.Set("Content-Type", "application/json")
	engine.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d: %s", http.StatusBadRequest, recorder.Code, recorder.Body.String())
	}
	assertErrorCode(t, recorder.Body.Bytes(), "VALIDATION_ERROR")
}

func newCommentTestRouter(t *testing.T, withRateLimit bool) (*gin.Engine, sqlmock.Sqlmock) {
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
	commentHandler := handler.NewCommentHandler(service.NewCommentService(database))
	dependencies := Dependencies{
		AuthHandler:    handler.NewAuthHandler(authService),
		ArticleHandler: handler.NewArticleHandler(service.NewArticleService(database)),
		CommentHandler: commentHandler,
		AuthMiddleware: middleware.Authenticate(authService),
	}
	if withRateLimit {
		dependencies.CommentRateLimit = middleware.NewCommentRateLimiter(5, time.Minute).Handler()
	}
	return New(dependencies), mock
}

func expectPublishedArticleForComments(mock sqlmock.Sqlmock, slug string) {
	mock.ExpectQuery(`SELECT "id","title","slug" FROM "articles" WHERE slug = \$1 AND status = \$2 ORDER BY "articles"\."id" LIMIT \$3`).
		WithArgs(slug, model.ArticleStatusPublished, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "title", "slug"}).
			AddRow(uint64(1), "Published", slug))
}

func publicCommentRows() *sqlmock.Rows {
	return sqlmock.NewRows([]string{"id", "nickname", "content", "created_at"})
}
