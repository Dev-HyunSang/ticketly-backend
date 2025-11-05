package handler

import (
	"github.com/dev-hyunsang/ticketly-backend/internal/usecase"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type UserHandler struct {
	userUseCase usecase.UserUseCase
}

func NewUserHandler(userUseCase usecase.UserUseCase) *UserHandler {
	return &UserHandler{
		userUseCase: userUseCase,
	}
}

// GetProfile retrieves the current user's profile
func (h *UserHandler) GetProfile(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uuid.UUID)

	user, err := h.userUseCase.GetUserByID(userID)
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

// GetUserByID retrieves a user by ID
func (h *UserHandler) GetUserByID(c *fiber.Ctx) error {
	userID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID",
		})
	}

	user, err := h.userUseCase.GetUserByID(userID)
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
