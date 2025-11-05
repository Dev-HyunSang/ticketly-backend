package handler

import (
	"github.com/dev-hyunsang/ticketly-backend/internal/usecase"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type OrganizationHandler struct {
	orgUseCase usecase.OrganizationUseCase
}

func NewOrganizationHandler(orgUseCase usecase.OrganizationUseCase) *OrganizationHandler {
	return &OrganizationHandler{
		orgUseCase: orgUseCase,
	}
}

type CreateOrganizationRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	LogoURL     string `json:"logo_url"`
}

type UpdateOrganizationRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	LogoURL     string `json:"logo_url"`
}

type AddMemberRequest struct {
	UserID uuid.UUID `json:"user_id"`
	Role   string    `json:"role"` // "admin" or "member"
}

type UpdateMemberRoleRequest struct {
	Role string `json:"role"` // "admin" or "member"
}

// CreateOrganization creates a new organization
func (h *OrganizationHandler) CreateOrganization(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uuid.UUID)

	var req CreateOrganizationRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	org, err := h.orgUseCase.CreateOrganization(req.Name, req.Description, req.LogoURL, userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message":      "Organization created successfully",
		"organization": org,
	})
}

// GetOrganization retrieves an organization by ID
func (h *OrganizationHandler) GetOrganization(c *fiber.Ctx) error {
	orgID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid organization ID",
		})
	}

	org, err := h.orgUseCase.GetOrganization(orgID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Organization not found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"organization": org,
	})
}

// GetMyOrganizations retrieves all organizations the current user belongs to
func (h *OrganizationHandler) GetMyOrganizations(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uuid.UUID)

	orgs, err := h.orgUseCase.GetUserOrganizations(userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"organizations": orgs,
	})
}

// UpdateOrganization updates an organization (admin only)
func (h *OrganizationHandler) UpdateOrganization(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uuid.UUID)

	orgID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid organization ID",
		})
	}

	var req UpdateOrganizationRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	err = h.orgUseCase.UpdateOrganization(orgID, req.Name, req.Description, req.LogoURL, userID)
	if err != nil {
		status := fiber.StatusInternalServerError
		if err.Error() == "permission denied: admin role required" {
			status = fiber.StatusForbidden
		}
		return c.Status(status).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Organization updated successfully",
	})
}

// DeleteOrganization deletes an organization (owner only)
func (h *OrganizationHandler) DeleteOrganization(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uuid.UUID)

	orgID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid organization ID",
		})
	}

	err = h.orgUseCase.DeleteOrganization(orgID, userID)
	if err != nil {
		status := fiber.StatusInternalServerError
		if err.Error() == "only organization owner can delete the organization" {
			status = fiber.StatusForbidden
		}
		return c.Status(status).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Organization deleted successfully",
	})
}

// GetMembers retrieves all members of an organization
func (h *OrganizationHandler) GetMembers(c *fiber.Ctx) error {
	orgID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid organization ID",
		})
	}

	members, err := h.orgUseCase.GetMembers(orgID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"members": members,
	})
}

// AddMember adds a member to an organization (admin only)
func (h *OrganizationHandler) AddMember(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uuid.UUID)

	orgID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid organization ID",
		})
	}

	var req AddMemberRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	err = h.orgUseCase.AddMember(orgID, req.UserID, userID, req.Role)
	if err != nil {
		status := fiber.StatusInternalServerError
		if err.Error() == "permission denied: admin role required" {
			status = fiber.StatusForbidden
		}
		return c.Status(status).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Member added successfully",
	})
}

// RemoveMember removes a member from an organization (admin only)
func (h *OrganizationHandler) RemoveMember(c *fiber.Ctx) error {
	requesterID := c.Locals("userID").(uuid.UUID)

	orgID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid organization ID",
		})
	}

	userID, err := uuid.Parse(c.Params("userId"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID",
		})
	}

	err = h.orgUseCase.RemoveMember(orgID, userID, requesterID)
	if err != nil {
		status := fiber.StatusInternalServerError
		if err.Error() == "permission denied: admin role required" || err.Error() == "cannot remove organization owner" {
			status = fiber.StatusForbidden
		}
		return c.Status(status).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Member removed successfully",
	})
}

// UpdateMemberRole updates a member's role (admin only)
func (h *OrganizationHandler) UpdateMemberRole(c *fiber.Ctx) error {
	requesterID := c.Locals("userID").(uuid.UUID)

	orgID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid organization ID",
		})
	}

	userID, err := uuid.Parse(c.Params("userId"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID",
		})
	}

	var req UpdateMemberRoleRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	err = h.orgUseCase.UpdateMemberRole(orgID, userID, requesterID, req.Role)
	if err != nil {
		status := fiber.StatusInternalServerError
		if err.Error() == "permission denied: admin role required" || err.Error() == "cannot change organization owner's role" {
			status = fiber.StatusForbidden
		}
		return c.Status(status).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Member role updated successfully",
	})
}
