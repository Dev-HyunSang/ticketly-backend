package usecase

import (
	"errors"
	"fmt"
	"time"

	"github.com/dev-hyunsang/ticketly-backend/internal/domain"
	"github.com/dev-hyunsang/ticketly-backend/internal/repository/mysql"
	"github.com/google/uuid"
)

type PaymentUseCase interface {
	CreatePayment(req CreatePaymentRequest, userID *uuid.UUID) (*domain.Payment, error)
	GetPaymentByID(paymentID uuid.UUID) (*domain.Payment, error)
	GetPaymentByOrderID(orderID string) (*domain.Payment, error)
	GetUserPayments(userID uuid.UUID) ([]*domain.Payment, error)
	GetEventPayments(eventID uuid.UUID) ([]*domain.Payment, error)
	UpdatePaymentStatus(paymentID uuid.UUID, status string, paymentKey string) error
	CompletePayment(orderID string, paymentKey string) (*domain.Payment, error)
}

type CreatePaymentRequest struct {
	EventID        uuid.UUID `json:"event_id"`
	EventTitle     string    `json:"event_title"`
	TicketQuantity int       `json:"ticket_quantity"`
	TotalPrice     float64   `json:"total_price"`
	BuyerName      string    `json:"buyer_name"`
	BuyerEmail     string    `json:"buyer_email"`
	BuyerPhone     string    `json:"buyer_phone"`
}

type paymentUseCase struct {
	paymentRepo *mysql.PaymentRepository
	eventRepo   domain.EventRepository
}

func NewPaymentUseCase(paymentRepo *mysql.PaymentRepository, eventRepo domain.EventRepository) PaymentUseCase {
	return &paymentUseCase{
		paymentRepo: paymentRepo,
		eventRepo:   eventRepo,
	}
}

func (uc *paymentUseCase) CreatePayment(req CreatePaymentRequest, userID *uuid.UUID) (*domain.Payment, error) {
	// Validate required fields
	if req.EventID == uuid.Nil {
		return nil, errors.New("event ID is required")
	}
	if req.EventTitle == "" {
		return nil, errors.New("event title is required")
	}
	if req.TicketQuantity <= 0 {
		return nil, errors.New("ticket quantity must be positive")
	}
	if req.TotalPrice < 0 {
		return nil, errors.New("total price must be non-negative")
	}
	if req.BuyerName == "" {
		return nil, errors.New("buyer name is required")
	}
	if req.BuyerEmail == "" {
		return nil, errors.New("buyer email is required")
	}
	if req.BuyerPhone == "" {
		return nil, errors.New("buyer phone is required")
	}

	// Check if event exists
	event, err := uc.eventRepo.GetByID(req.EventID)
	if err != nil {
		return nil, fmt.Errorf("event not found: %w", err)
	}

	// Check if enough tickets are available
	if event.AvailableTickets < req.TicketQuantity {
		return nil, errors.New("not enough tickets available")
	}

	// Generate order ID
	orderID := fmt.Sprintf("ORDER-%s", uuid.New().String()[:8])

	payment := &domain.Payment{
		ID:             uuid.New(),
		EventID:        req.EventID,
		UserID:         userID,
		EventTitle:     req.EventTitle,
		TicketQuantity: req.TicketQuantity,
		TotalPrice:     req.TotalPrice,
		Currency:       "KRW",
		BuyerName:      req.BuyerName,
		BuyerEmail:     req.BuyerEmail,
		BuyerPhone:     req.BuyerPhone,
		OrderID:        orderID,
		Status:         "pending",
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	return uc.paymentRepo.Create(payment)
}

func (uc *paymentUseCase) GetPaymentByID(paymentID uuid.UUID) (*domain.Payment, error) {
	return uc.paymentRepo.GetByID(paymentID)
}

func (uc *paymentUseCase) GetPaymentByOrderID(orderID string) (*domain.Payment, error) {
	return uc.paymentRepo.GetByOrderID(orderID)
}

func (uc *paymentUseCase) GetUserPayments(userID uuid.UUID) ([]*domain.Payment, error) {
	return uc.paymentRepo.GetByUserID(userID)
}

func (uc *paymentUseCase) GetEventPayments(eventID uuid.UUID) ([]*domain.Payment, error) {
	return uc.paymentRepo.GetByEventID(eventID)
}

func (uc *paymentUseCase) UpdatePaymentStatus(paymentID uuid.UUID, status string, paymentKey string) error {
	validStatuses := map[string]bool{
		"pending": true, "completed": true, "failed": true, "cancelled": true, "refunded": true,
	}
	if !validStatuses[status] {
		return errors.New("invalid payment status")
	}

	return uc.paymentRepo.UpdateStatus(paymentID, status, paymentKey)
}

func (uc *paymentUseCase) CompletePayment(orderID string, paymentKey string) (*domain.Payment, error) {
	// Get payment by order ID
	payment, err := uc.paymentRepo.GetByOrderID(orderID)
	if err != nil {
		return nil, fmt.Errorf("payment not found: %w", err)
	}

	// Check if payment is in pending status
	if payment.Status != "pending" {
		return nil, errors.New("payment is not in pending status")
	}

	// Update payment status to completed
	err = uc.paymentRepo.UpdateStatus(payment.ID, "completed", paymentKey)
	if err != nil {
		return nil, fmt.Errorf("failed to update payment status: %w", err)
	}

	// Reserve tickets
	err = uc.eventRepo.UpdateAvailableTickets(payment.EventID, -payment.TicketQuantity)
	if err != nil {
		// Rollback payment status if ticket update fails
		_ = uc.paymentRepo.UpdateStatus(payment.ID, "failed", paymentKey)
		return nil, fmt.Errorf("failed to reserve tickets: %w", err)
	}

	// Get updated payment
	return uc.paymentRepo.GetByID(payment.ID)
}
