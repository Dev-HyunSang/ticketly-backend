package handler

import (
	"github.com/dev-hyunsang/ticketly-backend/internal/usecase"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type PaymentHandler struct {
	paymentUseCase usecase.PaymentUseCase
}

func NewPaymentHandler(paymentUseCase usecase.PaymentUseCase) *PaymentHandler {
	return &PaymentHandler{
		paymentUseCase: paymentUseCase,
	}
}

// CreatePayment creates a new payment record
func (h *PaymentHandler) CreatePayment(c *fiber.Ctx) error {
	// Get user ID from context (optional for guest checkout)
	var userID *uuid.UUID
	if id, ok := c.Locals("userID").(uuid.UUID); ok {
		userID = &id
	}

	var req usecase.CreatePaymentRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	payment, err := h.paymentUseCase.CreatePayment(req, userID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Payment created successfully",
		"payment": payment,
	})
}

// GetPayment retrieves a payment by ID
func (h *PaymentHandler) GetPayment(c *fiber.Ctx) error {
	paymentID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid payment ID",
		})
	}

	payment, err := h.paymentUseCase.GetPaymentByID(paymentID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Payment not found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"payment": payment,
	})
}

// GetPaymentByOrderID retrieves a payment by order ID
func (h *PaymentHandler) GetPaymentByOrderID(c *fiber.Ctx) error {
	orderID := c.Params("orderId")
	if orderID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Order ID is required",
		})
	}

	payment, err := h.paymentUseCase.GetPaymentByOrderID(orderID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Payment not found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"payment": payment,
	})
}

// GetMyPayments retrieves all payments for the current user
func (h *PaymentHandler) GetMyPayments(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uuid.UUID)

	payments, err := h.paymentUseCase.GetUserPayments(userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"payments": payments,
	})
}

// GetEventPayments retrieves all payments for an event
func (h *PaymentHandler) GetEventPayments(c *fiber.Ctx) error {
	eventID, err := uuid.Parse(c.Params("eventId"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid event ID",
		})
	}

	payments, err := h.paymentUseCase.GetEventPayments(eventID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"payments": payments,
	})
}

// GetEventAttendees retrieves the attendee list for an event (completed payments only)
func (h *PaymentHandler) GetEventAttendees(c *fiber.Ctx) error {
	eventID, err := uuid.Parse(c.Params("eventId"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid event ID",
		})
	}

	attendees, err := h.paymentUseCase.GetEventAttendees(eventID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"attendees": attendees,
		"count":     len(attendees),
	})
}

// CompletePayment completes a payment after payment gateway confirmation
func (h *PaymentHandler) CompletePayment(c *fiber.Ctx) error {
	type CompleteRequest struct {
		OrderID    string `json:"order_id"`
		PaymentKey string `json:"payment_key"`
	}

	var req CompleteRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if req.OrderID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Order ID is required",
		})
	}

	payment, err := h.paymentUseCase.CompletePayment(req.OrderID, req.PaymentKey)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Payment completed successfully",
		"payment": payment,
	})
}

// CancelPayment cancels a payment and restores tickets
func (h *PaymentHandler) CancelPayment(c *fiber.Ctx) error {
	// Get user ID from context (optional for guest checkout)
	var userID *uuid.UUID
	if id, ok := c.Locals("userID").(uuid.UUID); ok {
		userID = &id
	}

	paymentID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid payment ID",
		})
	}

	payment, err := h.paymentUseCase.CancelPayment(paymentID, userID)
	if err != nil {
		status := fiber.StatusBadRequest
		if err.Error() == "permission denied: you can only cancel your own payments" {
			status = fiber.StatusForbidden
		}
		return c.Status(status).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Payment cancelled successfully",
		"payment": payment,
	})
}
