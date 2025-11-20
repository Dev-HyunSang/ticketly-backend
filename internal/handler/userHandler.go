package handler

import (
	"github.com/dev-hyunsang/ticketly-backend/internal/usecase"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type UserHandler struct {
	userUseCase usecase.UserUseCase
}

type UpdateUserRequest struct {
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	NickName    string `json:"nick_name"`
	Birthday    string `json:"birthday"`
	PhoneNumber string `json:"phone_number"`
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

// UpdateProfile updates the current user's profile
func (h *UserHandler) UpdateProfile(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uuid.UUID)

	var req UpdateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Get current user data
	user, err := h.userUseCase.GetUserByID(userID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	// Update only provided fields
	if req.FirstName != "" {
		user.FirstName = req.FirstName
	}
	if req.LastName != "" {
		user.LastName = req.LastName
	}
	if req.NickName != "" {
		user.NickName = req.NickName
	}
	if req.Birthday != "" {
		user.Birthday = req.Birthday
	}
	if req.PhoneNumber != "" {
		user.PhoneNumber = req.PhoneNumber
	}

	// Update user
	if err := h.userUseCase.Update(user); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update user profile",
		})
	}

	// Remove password from response
	user.Password = ""

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Profile updated successfully",
		"user":    user,
	})
}
