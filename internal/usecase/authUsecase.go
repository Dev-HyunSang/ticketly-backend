package usecase

import (
	"fmt"
	"time"

	"github.com/dev-hyunsang/ticketly-backend/internal/domain"
	"github.com/dev-hyunsang/ticketly-backend/internal/repository/mysql"
	"github.com/dev-hyunsang/ticketly-backend/internal/repository/redis"
	"github.com/dev-hyunsang/ticketly-backend/internal/util"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthUseCase struct {
	userRepo  domain.UserRepository
	tokenRepo *redis.TokenRepository
	jwtUtil   *util.JWTUtil
}

type LoginResponse struct {
	AccessToken  string       `json:"access_token"`
	RefreshToken string       `json:"refresh_token"`
	User         *domain.User `json:"user"`
}

type RegisterRequest struct {
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	NickName    string `json:"nick_name"`
	Birthday    string `json:"birthday"`
	Email       string `json:"email"`
	Password    string `json:"password"`
	PhoneNumber string `json:"phone_number"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func NewAuthUseCase(userRepo *mysql.UserRepository, tokenRepo *redis.TokenRepository, jwtUtil *util.JWTUtil) *AuthUseCase {
	return &AuthUseCase{
		userRepo:  userRepo,
		tokenRepo: tokenRepo,
		jwtUtil:   jwtUtil,
	}
}

// Register creates a new user account
func (uc *AuthUseCase) Register(req *RegisterRequest) (*domain.User, error) {
	// Validate input
	if req.Email == "" || req.Password == "" {
		return nil, domain.ErrInvalidInput
	}

	// Check if user already exists
	existingUser, err := uc.userRepo.GetUserByEmail(req.Email)
	if err == nil && existingUser != nil {
		return nil, fmt.Errorf("user with email %s already exists", req.Email)
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Create user
	user := &domain.User{
		ID:          uuid.New(),
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		NickName:    req.NickName,
		Birthday:    req.Birthday,
		Email:       req.Email,
		Password:    string(hashedPassword),
		PhoneNumber: req.PhoneNumber,
		IsValid:     true,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	createdUser, err := uc.userRepo.Save(user)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Don't return password in response
	createdUser.Password = ""

	return createdUser, nil
}

// Login authenticates user and returns JWT tokens
func (uc *AuthUseCase) Login(req *LoginRequest) (*LoginResponse, error) {
	// Validate input
	if req.Email == "" || req.Password == "" {
		return nil, domain.ErrInvalidInput
	}

	// Get user by email
	user, err := uc.userRepo.GetUserByEmail(req.Email)
	if err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	// Check if user is valid
	if !user.IsValid {
		return nil, fmt.Errorf("user account is not valid")
	}

	// Compare password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	// Generate access token
	accessToken, err := uc.jwtUtil.GenerateAccessToken(user.ID, user.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	// Generate refresh token
	refreshToken, err := uc.jwtUtil.GenerateRefreshToken(user.ID, user.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	// Save refresh token to Redis
	err = uc.tokenRepo.SaveRefreshToken(user.ID, refreshToken, uc.jwtUtil.GetRefreshTokenExpiry())
	if err != nil {
		return nil, fmt.Errorf("failed to save refresh token: %w", err)
	}

	// Don't return password in response
	user.Password = ""

	return &LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         user,
	}, nil
}

// RefreshToken generates new access token using refresh token
func (uc *AuthUseCase) RefreshToken(refreshToken string) (string, error) {
	// Validate refresh token
	claims, err := uc.jwtUtil.ValidateRefreshToken(refreshToken)
	if err != nil {
		return "", fmt.Errorf("invalid refresh token: %w", err)
	}

	// Check if refresh token exists in Redis
	storedToken, err := uc.tokenRepo.GetRefreshToken(claims.UserID)
	if err != nil {
		return "", fmt.Errorf("refresh token not found or expired")
	}

	// Compare tokens
	if storedToken != refreshToken {
		return "", fmt.Errorf("invalid refresh token")
	}

	// Generate new access token
	accessToken, err := uc.jwtUtil.GenerateAccessToken(claims.UserID, claims.Email)
	if err != nil {
		return "", fmt.Errorf("failed to generate access token: %w", err)
	}

	return accessToken, nil
}

// Logout invalidates user's tokens
func (uc *AuthUseCase) Logout(userID uuid.UUID, accessToken string) error {
	// Delete refresh token from Redis
	err := uc.tokenRepo.DeleteRefreshToken(userID)
	if err != nil {
		return fmt.Errorf("failed to delete refresh token: %w", err)
	}

	// Blacklist current access token
	err = uc.tokenRepo.BlacklistToken(accessToken, uc.jwtUtil.GetAccessTokenExpiry())
	if err != nil {
		return fmt.Errorf("failed to blacklist token: %w", err)
	}

	return nil
}

// ValidateAccessToken validates access token and checks if it's blacklisted
func (uc *AuthUseCase) ValidateAccessToken(token string) (*util.Claims, error) {
	// Check if token is blacklisted
	isBlacklisted, err := uc.tokenRepo.IsTokenBlacklisted(token)
	if err != nil {
		return nil, fmt.Errorf("failed to check token blacklist: %w", err)
	}

	if isBlacklisted {
		return nil, fmt.Errorf("token is blacklisted")
	}

	// Validate token
	claims, err := uc.jwtUtil.ValidateAccessToken(token)
	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	return claims, nil
}
