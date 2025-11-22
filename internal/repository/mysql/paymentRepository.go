package mysql

import (
	"context"
	"fmt"

	"github.com/dev-hyunsang/ticketly-backend/internal/domain"
	"github.com/dev-hyunsang/ticketly-backend/lib/ent"
	"github.com/dev-hyunsang/ticketly-backend/lib/ent/payment"
	"github.com/google/uuid"
)

type PaymentRepository struct {
	client *ent.Client
}

func NewPaymentRepository(client *ent.Client) *PaymentRepository {
	return &PaymentRepository{
		client: client,
	}
}

func (r *PaymentRepository) Create(p *domain.Payment) (*domain.Payment, error) {
	ctx := context.Background()

	builder := r.client.Payment.
		Create().
		SetID(p.ID).
		SetEventID(p.EventID).
		SetEventTitle(p.EventTitle).
		SetTicketQuantity(p.TicketQuantity).
		SetTotalPrice(p.TotalPrice).
		SetCurrency(p.Currency).
		SetBuyerName(p.BuyerName).
		SetBuyerEmail(p.BuyerEmail).
		SetBuyerPhone(p.BuyerPhone).
		SetStatus(payment.Status(p.Status))

	if p.UserID != nil {
		builder.SetUserID(*p.UserID)
	}

	if p.PaymentKey != "" {
		builder.SetPaymentKey(p.PaymentKey)
	}

	if p.OrderID != "" {
		builder.SetOrderID(p.OrderID)
	}

	createdPayment, err := builder.Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create payment: %w", err)
	}

	return r.mapToDomain(createdPayment), nil
}

func (r *PaymentRepository) GetByID(paymentID uuid.UUID) (*domain.Payment, error) {
	ctx := context.Background()

	p, err := r.client.Payment.
		Query().
		Where(payment.ID(paymentID)).
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("failed to get payment: %w", err)
	}

	return r.mapToDomain(p), nil
}

func (r *PaymentRepository) GetByOrderID(orderID string) (*domain.Payment, error) {
	ctx := context.Background()

	p, err := r.client.Payment.
		Query().
		Where(payment.OrderID(orderID)).
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("failed to get payment by order ID: %w", err)
	}

	return r.mapToDomain(p), nil
}

func (r *PaymentRepository) GetByUserID(userID uuid.UUID) ([]*domain.Payment, error) {
	ctx := context.Background()

	payments, err := r.client.Payment.
		Query().
		Where(payment.UserID(userID)).
		Order(ent.Desc(payment.FieldCreatedAt)).
		All(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get payments by user ID: %w", err)
	}

	result := make([]*domain.Payment, len(payments))
	for i, p := range payments {
		result[i] = r.mapToDomain(p)
	}

	return result, nil
}

func (r *PaymentRepository) GetByEventID(eventID uuid.UUID) ([]*domain.Payment, error) {
	ctx := context.Background()

	payments, err := r.client.Payment.
		Query().
		Where(payment.EventID(eventID)).
		Order(ent.Desc(payment.FieldCreatedAt)).
		All(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get payments by event ID: %w", err)
	}

	result := make([]*domain.Payment, len(payments))
	for i, p := range payments {
		result[i] = r.mapToDomain(p)
	}

	return result, nil
}

func (r *PaymentRepository) UpdateStatus(paymentID uuid.UUID, status string, paymentKey string) error {
	ctx := context.Background()

	builder := r.client.Payment.
		UpdateOneID(paymentID).
		SetStatus(payment.Status(status))

	if paymentKey != "" {
		builder.SetPaymentKey(paymentKey)
	}

	err := builder.Exec(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return domain.ErrNotFound
		}
		return fmt.Errorf("failed to update payment status: %w", err)
	}

	return nil
}

func (r *PaymentRepository) mapToDomain(p *ent.Payment) *domain.Payment {
	var userID *uuid.UUID
	if p.UserID != uuid.Nil {
		userID = &p.UserID
	}

	return &domain.Payment{
		ID:             p.ID,
		EventID:        p.EventID,
		UserID:         userID,
		EventTitle:     p.EventTitle,
		TicketQuantity: p.TicketQuantity,
		TotalPrice:     p.TotalPrice,
		Currency:       p.Currency,
		BuyerName:      p.BuyerName,
		BuyerEmail:     p.BuyerEmail,
		BuyerPhone:     p.BuyerPhone,
		PaymentKey:     p.PaymentKey,
		OrderID:        p.OrderID,
		Status:         string(p.Status),
		CreatedAt:      p.CreatedAt,
		UpdatedAt:      p.UpdatedAt,
	}
}
