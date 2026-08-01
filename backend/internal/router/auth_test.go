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
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"personal-blog/backend/internal/handler"
	"personal-blog/backend/internal/middleware"
	"personal-blog/backend/internal/service"
)

const testJWTSecret = "test-jwt-secret-with-at-least-32-characters"

func TestAdministratorLoginAndCurrentAdministrator(t *testing.T) {
	engine, mock := newAuthTestRouter(t)
	passwordHash, err := bcrypt.GenerateFromPassword([]byte("correct-password"), bcrypt.MinCost)
	if err != nil {
		t.Fatalf("hash test password: %v", err)
	}

	expectAdministratorByUsername(mock, string(passwordHash))
	loginRecorder := httptest.NewRecorder()
	loginRequest := httptest.NewRequest(
		http.MethodPost,
		"/api/v1/admin/auth/login",
		bytes.NewBufferString(`{"username":"admin","password":"correct-password"}`),
	)
	loginRequest.Header.Set("Content-Type", "application/json")
	engine.ServeHTTP(loginRecorder, loginRequest)

	if loginRecorder.Code != http.StatusOK {
		t.Fatalf("expected login status %d, got %d: %s", http.StatusOK, loginRecorder.Code, loginRecorder.Body.String())
	}

	var loginResponse struct {
		Data struct {
			Token     string `json:"token"`
			ExpiresIn int64  `json:"expires_in"`
		} `json:"data"`
	}
	if err := json.Unmarshal(loginRecorder.Body.Bytes(), &loginResponse); err != nil {
		t.Fatalf("decode login response: %v", err)
	}
	if loginResponse.Data.Token == "" {
		t.Fatal("expected a JWT token")
	}
	if loginResponse.Data.ExpiresIn != 86400 {
		t.Fatalf("expected expires_in 86400, got %d", loginResponse.Data.ExpiresIn)
	}

	expectAdministratorByID(mock, string(passwordHash))
	meRecorder := httptest.NewRecorder()
	meRequest := httptest.NewRequest(http.MethodGet, "/api/v1/admin/auth/me", nil)
	meRequest.Header.Set("Authorization", "Bearer "+loginResponse.Data.Token)
	engine.ServeHTTP(meRecorder, meRequest)

	if meRecorder.Code != http.StatusOK {
		t.Fatalf("expected current administrator status %d, got %d: %s", http.StatusOK, meRecorder.Code, meRecorder.Body.String())
	}
	if strings.Contains(meRecorder.Body.String(), "password") {
		t.Fatalf("current administrator response exposed a password field: %s", meRecorder.Body.String())
	}

	var meResponse struct {
		Data struct {
			ID       uint64 `json:"id"`
			Username string `json:"username"`
		} `json:"data"`
	}
	if err := json.Unmarshal(meRecorder.Body.Bytes(), &meResponse); err != nil {
		t.Fatalf("decode current administrator response: %v", err)
	}
	if meResponse.Data.ID != 1 || meResponse.Data.Username != "admin" {
		t.Fatalf("unexpected current administrator: %+v", meResponse.Data)
	}
}

func TestAdministratorLoginRejectsWrongPassword(t *testing.T) {
	engine, mock := newAuthTestRouter(t)
	passwordHash, err := bcrypt.GenerateFromPassword([]byte("correct-password"), bcrypt.MinCost)
	if err != nil {
		t.Fatalf("hash test password: %v", err)
	}
	expectAdministratorByUsername(mock, string(passwordHash))

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(
		http.MethodPost,
		"/api/v1/admin/auth/login",
		bytes.NewBufferString(`{"username":"admin","password":"wrong-password"}`),
	)
	request.Header.Set("Content-Type", "application/json")
	engine.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusUnauthorized {
		t.Fatalf("expected status %d, got %d: %s", http.StatusUnauthorized, recorder.Code, recorder.Body.String())
	}
	assertErrorCode(t, recorder.Body.Bytes(), "INVALID_CREDENTIALS")
}

func TestCurrentAdministratorRequiresToken(t *testing.T) {
	engine, _ := newAuthTestRouter(t)
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/api/v1/admin/auth/me", nil)

	engine.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusUnauthorized {
		t.Fatalf("expected status %d, got %d: %s", http.StatusUnauthorized, recorder.Code, recorder.Body.String())
	}
	assertErrorCode(t, recorder.Body.Bytes(), "UNAUTHORIZED")
}

func newAuthTestRouter(t *testing.T) (*gin.Engine, sqlmock.Sqlmock) {
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
	authHandler := handler.NewAuthHandler(authService)
	return New(Dependencies{
		AuthHandler:    authHandler,
		AuthMiddleware: middleware.Authenticate(authService),
	}), mock
}

func expectAdministratorByUsername(mock sqlmock.Sqlmock, passwordHash string) {
	mock.ExpectQuery(`SELECT \* FROM "administrators" WHERE username = \$1 ORDER BY "administrators"\."id" LIMIT \$2`).
		WithArgs("admin", 1).
		WillReturnRows(administratorRows(passwordHash))
}

func expectAdministratorByID(mock sqlmock.Sqlmock, passwordHash string) {
	mock.ExpectQuery(`SELECT \* FROM "administrators" WHERE "administrators"\."id" = \$1 ORDER BY "administrators"\."id" LIMIT \$2`).
		WithArgs(uint64(1), 1).
		WillReturnRows(administratorRows(passwordHash))
}

func administratorRows(passwordHash string) *sqlmock.Rows {
	now := time.Now().UTC()
	return sqlmock.NewRows([]string{"id", "username", "password_hash", "created_at", "updated_at"}).
		AddRow(uint64(1), "admin", passwordHash, now, now)
}

func assertErrorCode(t *testing.T, body []byte, expected string) {
	t.Helper()
	var response struct {
		Error struct {
			Code string `json:"code"`
		} `json:"error"`
	}
	if err := json.Unmarshal(body, &response); err != nil {
		t.Fatalf("decode error response: %v", err)
	}
	if response.Error.Code != expected {
		t.Fatalf("expected error code %q, got %q", expected, response.Error.Code)
	}
}
