package middleware

import (
	"strings"

	"github.com/dev-hyunsang/ticketly-backend/internal/usecase"
	"github.com/gofiber/fiber/v2"
)

type AuthMiddleware struct {
	authUseCase *usecase.AuthUseCase
}

func NewAuthMiddleware(authUseCase *usecase.AuthUseCase) *AuthMiddleware {
	return &AuthMiddleware{
		authUseCase: authUseCase,
	}
}

// Authenticate validates JWT token and sets user info in context
func (m *AuthMiddleware) Authenticate(c *fiber.Ctx) error {
	// Get Authorization header
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "authorization header required",
		})
	}

	// Check if it starts with "Bearer "
	if !strings.HasPrefix(authHeader, "Bearer ") {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "invalid authorization header format",
		})
	}

	// Extract token
	token := authHeader[7:]

	// Validate token
	claims, err := m.authUseCase.ValidateAccessToken(token)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "invalid or expired token",
		})
	}

	// Set user info in context
	c.Locals("userID", claims.UserID)
	c.Locals("email", claims.Email)

	return c.Next()
}
