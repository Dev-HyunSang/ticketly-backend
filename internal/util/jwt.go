package util

import (
	"fmt"
	"time"

	"github.com/dev-hyunsang/ticketly-backend/config"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type Claims struct {
	UserID uuid.UUID `json:"user_id"`
	Email  string    `json:"email"`
	jwt.RegisteredClaims
}

type JWTUtil struct {
	accessSecret  string
	refreshSecret string
}

func NewJWTUtil() *JWTUtil {
	accessSecret := config.Getenv("JWT_ACCESS_SECRET")
	refreshSecret := config.Getenv("JWT_REFRESH_SECRET")

	if accessSecret == "" {
		accessSecret = "default-access-secret-change-this-in-production"
	}
	if refreshSecret == "" {
		refreshSecret = "default-refresh-secret-change-this-in-production"
	}

	return &JWTUtil{
		accessSecret:  accessSecret,
		refreshSecret: refreshSecret,
	}
}

// GenerateAccessToken generates JWT access token (15 minutes expiry)
func (j *JWTUtil) GenerateAccessToken(userID uuid.UUID, email string) (string, error) {
	claims := &Claims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(j.accessSecret))
	if err != nil {
		return "", fmt.Errorf("failed to generate access token: %w", err)
	}

	return tokenString, nil
}

// GenerateRefreshToken generates JWT refresh token (7 days expiry)
func (j *JWTUtil) GenerateRefreshToken(userID uuid.UUID, email string) (string, error) {
	claims := &Claims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(j.refreshSecret))
	if err != nil {
		return "", fmt.Errorf("failed to generate refresh token: %w", err)
	}

	return tokenString, nil
}

// ValidateAccessToken validates JWT access token
func (j *JWTUtil) ValidateAccessToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(j.accessSecret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}

// ValidateRefreshToken validates JWT refresh token
func (j *JWTUtil) ValidateRefreshToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(j.refreshSecret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}

// GetAccessTokenExpiry returns access token expiration duration
func (j *JWTUtil) GetAccessTokenExpiry() time.Duration {
	return 15 * time.Minute
}

// GetRefreshTokenExpiry returns refresh token expiration duration
func (j *JWTUtil) GetRefreshTokenExpiry() time.Duration {
	return 7 * 24 * time.Hour
}
