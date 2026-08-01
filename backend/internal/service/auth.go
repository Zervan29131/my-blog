package service

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"personal-blog/backend/internal/model"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidToken       = errors.New("invalid token")
	ErrAdministratorGone  = errors.New("administrator not found")
)

type AuthService struct {
	database  *gorm.DB
	jwtSecret []byte
	expiresIn time.Duration
}

type LoginResult struct {
	Token     string
	ExpiresIn int64
}

type AdminClaims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func NewAuthService(database *gorm.DB, jwtSecret string, expiresHours int) *AuthService {
	return &AuthService{
		database:  database,
		jwtSecret: []byte(jwtSecret),
		expiresIn: time.Duration(expiresHours) * time.Hour,
	}
}

func (service *AuthService) InitializeAdministrator(
	ctx context.Context,
	username string,
	password string,
) (bool, error) {
	var count int64
	if err := service.database.WithContext(ctx).Model(&model.Administrator{}).Count(&count).Error; err != nil {
		return false, fmt.Errorf("count administrators: %w", err)
	}
	if count > 0 {
		return false, nil
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return false, fmt.Errorf("hash administrator password: %w", err)
	}

	administrator := model.Administrator{
		Username:     username,
		PasswordHash: string(passwordHash),
	}
	if err := service.database.WithContext(ctx).Create(&administrator).Error; err != nil {
		return false, fmt.Errorf("create administrator: %w", err)
	}

	return true, nil
}

func (service *AuthService) Login(ctx context.Context, username, password string) (LoginResult, error) {
	var administrator model.Administrator
	result := service.database.WithContext(ctx).
		Where("username = ?", username).
		First(&administrator)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return LoginResult{}, ErrInvalidCredentials
	}
	if result.Error != nil {
		return LoginResult{}, fmt.Errorf("find administrator: %w", result.Error)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(administrator.PasswordHash), []byte(password)); err != nil {
		return LoginResult{}, ErrInvalidCredentials
	}

	token, err := service.issueToken(administrator)
	if err != nil {
		return LoginResult{}, err
	}

	return LoginResult{
		Token:     token,
		ExpiresIn: int64(service.expiresIn.Seconds()),
	}, nil
}

func (service *AuthService) ParseToken(tokenValue string) (*AdminClaims, error) {
	claims := &AdminClaims{}
	token, err := jwt.ParseWithClaims(
		tokenValue,
		claims,
		func(token *jwt.Token) (any, error) {
			if token.Method.Alg() != jwt.SigningMethodHS256.Alg() {
				return nil, ErrInvalidToken
			}
			return service.jwtSecret, nil
		},
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
		jwt.WithExpirationRequired(),
	)
	if err != nil || !token.Valid {
		return nil, ErrInvalidToken
	}

	administratorID, err := strconv.ParseUint(claims.Subject, 10, 64)
	if err != nil || administratorID == 0 || claims.Username == "" {
		return nil, ErrInvalidToken
	}

	return claims, nil
}

func (service *AuthService) CurrentAdministrator(ctx context.Context, id uint64) (model.Administrator, error) {
	var administrator model.Administrator
	result := service.database.WithContext(ctx).First(&administrator, id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return model.Administrator{}, ErrAdministratorGone
	}
	if result.Error != nil {
		return model.Administrator{}, fmt.Errorf("find current administrator: %w", result.Error)
	}
	return administrator, nil
}

func (service *AuthService) issueToken(administrator model.Administrator) (string, error) {
	now := time.Now().UTC()
	claims := AdminClaims{
		Username: administrator.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   strconv.FormatUint(administrator.ID, 10),
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(service.expiresIn)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(service.jwtSecret)
	if err != nil {
		return "", fmt.Errorf("sign JWT: %w", err)
	}
	return signedToken, nil
}
