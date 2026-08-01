package service

import (
	"context"
	"database/sql/driver"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestInitializeAdministratorStoresBcryptHash(t *testing.T) {
	database, mock := newServiceTestDatabase(t)
	mock.ExpectQuery(`SELECT count\(\*\) FROM "administrators"`).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))
	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO "administrators" \("username","password_hash","created_at","updated_at"\) VALUES \(\$1,\$2,\$3,\$4\) RETURNING "id"`).
		WithArgs("admin", bcryptHashMatcher{password: "correct-password"}, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(uint64(1)))
	mock.ExpectCommit()

	authService := NewAuthService(database, testServiceJWTSecret, 24)
	created, err := authService.InitializeAdministrator(context.Background(), "admin", "correct-password")
	if err != nil {
		t.Fatalf("initialize administrator: %v", err)
	}
	if !created {
		t.Fatal("expected administrator to be created")
	}
}

func TestInitializeAdministratorDoesNotCreateSecondAccount(t *testing.T) {
	database, mock := newServiceTestDatabase(t)
	mock.ExpectQuery(`SELECT count\(\*\) FROM "administrators"`).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

	authService := NewAuthService(database, testServiceJWTSecret, 24)
	created, err := authService.InitializeAdministrator(context.Background(), "another-admin", "password")
	if err != nil {
		t.Fatalf("initialize administrator: %v", err)
	}
	if created {
		t.Fatal("expected existing administrator to be preserved")
	}
}

func TestParseTokenRejectsExpiredJWT(t *testing.T) {
	authService := NewAuthService(nil, testServiceJWTSecret, 24)
	claims := AdminClaims{
		Username: "admin",
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   "1",
			IssuedAt:  jwt.NewNumericDate(time.Now().Add(-2 * time.Hour)),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(-time.Hour)),
		},
	}
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(testServiceJWTSecret))
	if err != nil {
		t.Fatalf("sign expired test token: %v", err)
	}

	if _, err := authService.ParseToken(token); !errors.Is(err, ErrInvalidToken) {
		t.Fatalf("expected expired token error, got %v", err)
	}
}

const testServiceJWTSecret = "test-service-jwt-secret-with-32-characters"

type bcryptHashMatcher struct {
	password string
}

func (matcher bcryptHashMatcher) Match(value driver.Value) bool {
	hash, ok := value.(string)
	if !ok || hash == matcher.password {
		return false
	}
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(matcher.password)) == nil
}

func newServiceTestDatabase(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
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
	return database, mock
}
