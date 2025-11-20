package handler

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/dev-hyunsang/ticketly-backend/internal/usecase"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type EventHandler struct {
	eventUseCase usecase.EventUseCase
}

type BuyEventTicket struct {
	EventID uuid.UUID `json:event_id`
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

// GetPublicEvent retrieves a public event by ID (no authentication required)
func (h *EventHandler) GetPublicEvent(c *fiber.Ctx) error {
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

	// Check if event is public
	if !event.IsPublic {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "This event is not publicly accessible",
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

	log.Println(events)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"events": events,
	})
}

// GetPopularEvents retrieves popular upcoming events
func (h *EventHandler) GetPopularEvents(c *fiber.Ctx) error {
	events, err := h.eventUseCase.GetPopularEvents()
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

func ConfirmPayment(paymentKey, orderId string, amount int) (string, error) {
	url := "https://api.tosspayments.com/v1/payments/confirm"

	// Prepare payload
	payload := map[string]interface{}{
		"paymentKey": paymentKey,
		"orderId":    orderId,
		"amount":     amount,
	}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("failed to marshal payload: %v", err)
	}

	// Authorization (replace this with your actual secret key)
	secretKey := "test_sk_GePWvyJnrKazDJ56ZLRa3gLzNN7E:"
	auth := "Basic " + base64.StdEncoding.EncodeToString([]byte(secretKey))

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return "", fmt.Errorf("failed to create HTTP request: %v", err)
	}

	req.Header.Set("Authorization", auth)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("HTTP request failed: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("payment confirm failed: %s", string(body))
	}

	return string(body), nil
}

func (h *EventHandler) BuyEvents(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uuid.UUID)

	req := new(BuyEventTicket)

	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "올바르지 않은 형식으로 요청하셨습니다. 확인 후 다시 시도해 주세요.",
		})
	}

	result, err := h.eventUseCase.GetEvent(userID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "해당 이벤트를 찾을 수 없습니다. 확인 후 다시 시도해 주세요.",
		})
	}

	log.Println(result)

	// ConfirmPayment()
	return nil
}
