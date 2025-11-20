package handler

import (
	"github.com/dev-hyunsang/ticketly-backend/internal/usecase"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type AuthHandler struct {
	authUseCase *usecase.AuthUseCase
}

func NewAuthHandler(authUseCase *usecase.AuthUseCase) *AuthHandler {
	return &AuthHandler{
		authUseCase: authUseCase,
	}
}

// Register handles user registration
func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var req usecase.RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	user, err := h.authUseCase.Register(&req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "user created successfully",
		"user":    user,
	})
}

// Login handles user authentication
func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req usecase.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	response, err := h.authUseCase.Login(&req)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(response)
}

// RefreshToken handles token refresh
func (h *AuthHandler) RefreshToken(c *fiber.Ctx) error {
	type RefreshRequest struct {
		RefreshToken string `json:"refresh_token"`
	}

	var req RefreshRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	accessToken, err := h.authUseCase.RefreshToken(req.RefreshToken)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"access_token": accessToken,
	})
}

// Logout handles user logout
func (h *AuthHandler) Logout(c *fiber.Ctx) error {
	// Get user ID from context (set by auth middleware)
	userID := c.Locals("userID").(uuid.UUID)

	// Get access token from header
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "authorization header required",
		})
	}

	// Remove "Bearer " prefix
	token := authHeader[7:]

	err := h.authUseCase.Logout(userID, token)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "logged out successfully",
	})
}

// Me returns current user information
func (h *AuthHandler) Me(c *fiber.Ctx) error {
	// Get user information from context (set by auth middleware)
	userID := c.Locals("userID").(uuid.UUID)

	// Get full user information from database
	user, err := h.authUseCase.GetUserByID(userID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	// Remove password from response
	user.Password = ""

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"user": user,
	})
}
