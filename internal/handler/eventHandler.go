package handler

import (
	"github.com/dev-hyunsang/ticketly-backend/internal/usecase"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type EventHandler struct {
	eventUseCase usecase.EventUseCase
}

func NewEventHandler(eventUseCase usecase.EventUseCase) *EventHandler {
	return &EventHandler{
		eventUseCase: eventUseCase,
	}
}

// CreateEvent creates a new event (admin only)
func (h *EventHandler) CreateEvent(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uuid.UUID)

	orgID, err := uuid.Parse(c.Params("orgId"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid organization ID",
		})
	}

	var req usecase.CreateEventRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	event, err := h.eventUseCase.CreateEvent(orgID, userID, req)
	if err != nil {
		status := fiber.StatusInternalServerError
		if err.Error() == "permission denied: only admins can create events" {
			status = fiber.StatusForbidden
		}
		return c.Status(status).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Event created successfully",
		"event":   event,
	})
}

// GetEvent retrieves an event by ID
func (h *EventHandler) GetEvent(c *fiber.Ctx) error {
	eventID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid event ID",
		})
	}

	event, err := h.eventUseCase.GetEvent(eventID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Event not found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"event": event,
	})
}

// GetOrganizationEvents retrieves all events for an organization
func (h *EventHandler) GetOrganizationEvents(c *fiber.Ctx) error {
	orgID, err := uuid.Parse(c.Params("orgId"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid organization ID",
		})
	}

	events, err := h.eventUseCase.GetOrganizationEvents(orgID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"events": events,
	})
}

// UpdateEvent updates an event (admin only)
func (h *EventHandler) UpdateEvent(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uuid.UUID)

	eventID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid event ID",
		})
	}

	var req usecase.UpdateEventRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	err = h.eventUseCase.UpdateEvent(eventID, userID, req)
	if err != nil {
		status := fiber.StatusInternalServerError
		if err.Error() == "permission denied: only admins can update events" {
			status = fiber.StatusForbidden
		}
		return c.Status(status).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Event updated successfully",
	})
}

// DeleteEvent deletes an event (admin only)
func (h *EventHandler) DeleteEvent(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uuid.UUID)

	eventID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid event ID",
		})
	}

	err = h.eventUseCase.DeleteEvent(eventID, userID)
	if err != nil {
		status := fiber.StatusInternalServerError
		if err.Error() == "permission denied: only admins can delete events" {
			status = fiber.StatusForbidden
		}
		return c.Status(status).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Event deleted successfully",
	})
}

// GetPublicEvents retrieves all public events
func (h *EventHandler) GetPublicEvents(c *fiber.Ctx) error {
	events, err := h.eventUseCase.GetPublicEvents()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"events": events,
	})
}

// GetUpcomingEvents retrieves upcoming public events
func (h *EventHandler) GetUpcomingEvents(c *fiber.Ctx) error {
	events, err := h.eventUseCase.GetUpcomingEvents()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"events": events,
	})
}

// SearchEvents searches events by keyword
func (h *EventHandler) SearchEvents(c *fiber.Ctx) error {
	keyword := c.Query("q")
	if keyword == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Search keyword is required",
		})
	}

	events, err := h.eventUseCase.SearchEvents(keyword)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"events": events,
	})
}
