package domain

import (
	"time"

	"github.com/google/uuid"
)

type Payment struct {
	ID             uuid.UUID  `json:"id"`
	EventID        uuid.UUID  `json:"event_id"`
	UserID         *uuid.UUID `json:"user_id,omitempty"`
	EventTitle     string     `json:"event_title"`
	TicketQuantity int        `json:"ticket_quantity"`
	TotalPrice     float64    `json:"total_price"`
	Currency       string     `json:"currency"`
	BuyerName      string     `json:"buyer_name"`
	BuyerEmail     string     `json:"buyer_email"`
	BuyerPhone     string     `json:"buyer_phone"`
	PaymentKey     string     `json:"payment_key,omitempty"`
	OrderID        string     `json:"order_id,omitempty"`
	Status         string     `json:"status"` // pending, completed, failed, cancelled, refunded
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

type PaymentWithEvent struct {
	Payment
	EventTitle string `json:"event_title"`
}

// PaymentRepository defines the interface for payment data access
type PaymentRepository interface {
	Create(payment *Payment) (*Payment, error)
	GetByID(paymentID uuid.UUID) (*Payment, error)
	GetByOrderID(orderID string) (*Payment, error)
	GetByUserID(userID uuid.UUID) ([]*Payment, error)
	GetByEventID(eventID uuid.UUID) ([]*Payment, error)
	UpdateStatus(paymentID uuid.UUID, status string, paymentKey string) error
}
