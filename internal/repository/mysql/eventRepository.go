package mysql

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/dev-hyunsang/ticketly-backend/internal/domain"
	"github.com/dev-hyunsang/ticketly-backend/lib/ent"
	"github.com/dev-hyunsang/ticketly-backend/lib/ent/event"
	"github.com/google/uuid"
)

type eventRepository struct {
	client *ent.Client
}

func NewEventRepository(client *ent.Client) domain.EventRepository {
	return &eventRepository{
		client: client,
	}
}

// Create creates a new event
func (r *eventRepository) Create(evt *domain.Event) (*domain.Event, error) {
	ctx := context.Background()

	createdEvent, err := r.client.Event.
		Create().
		SetID(evt.ID).
		SetOrganizationID(evt.OrganizationID).
		SetTitle(evt.Title).
		SetNillableDescription(&evt.Description).
		SetNillableLocation(&evt.Location).
		SetNillableVenue(&evt.Venue).
		SetStartTime(evt.StartTime).
		SetEndTime(evt.EndTime).
		SetTotalTickets(evt.TotalTickets).
		SetAvailableTickets(evt.AvailableTickets).
		SetTicketPrice(evt.TicketPrice).
		SetCurrency(evt.Currency).
		SetNillableThumbnailURL(&evt.ThumbnailURL).
		SetStatus(event.Status(evt.Status)).
		SetIsPublic(evt.IsPublic).
		SetCreatedBy(evt.CreatedBy).
		Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create event: %w", err)
	}

	return &domain.Event{
		ID:               createdEvent.ID,
		OrganizationID:   createdEvent.OrganizationID,
		Title:            createdEvent.Title,
		Description:      createdEvent.Description,
		Location:         createdEvent.Location,
		Venue:            createdEvent.Venue,
		StartTime:        createdEvent.StartTime,
		EndTime:          createdEvent.EndTime,
		TotalTickets:     createdEvent.TotalTickets,
		AvailableTickets: createdEvent.AvailableTickets,
		TicketPrice:      createdEvent.TicketPrice,
		Currency:         createdEvent.Currency,
		ThumbnailURL:     createdEvent.ThumbnailURL,
		Status:           string(createdEvent.Status),
		IsPublic:         createdEvent.IsPublic,
		CreatedBy:        createdEvent.CreatedBy,
		CreatedAt:        createdEvent.CreatedAt,
		UpdatedAt:        createdEvent.UpdatedAt,
	}, nil
}

// GetByID retrieves an event by ID
func (r *eventRepository) GetByID(eventID uuid.UUID) (*domain.Event, error) {
	ctx := context.Background()

	evt, err := r.client.Event.
		Query().
		Where(event.ID(eventID)).
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("failed to get event: %w", err)
	}

	return &domain.Event{
		ID:               evt.ID,
		OrganizationID:   evt.OrganizationID,
		Title:            evt.Title,
		Description:      evt.Description,
		Location:         evt.Location,
		Venue:            evt.Venue,
		StartTime:        evt.StartTime,
		EndTime:          evt.EndTime,
		TotalTickets:     evt.TotalTickets,
		AvailableTickets: evt.AvailableTickets,
		TicketPrice:      evt.TicketPrice,
		Currency:         evt.Currency,
		ThumbnailURL:     evt.ThumbnailURL,
		Status:           string(evt.Status),
		IsPublic:         evt.IsPublic,
		CreatedBy:        evt.CreatedBy,
		CreatedAt:        evt.CreatedAt,
		UpdatedAt:        evt.UpdatedAt,
	}, nil
}

// GetByOrganizationID retrieves all events for an organization
func (r *eventRepository) GetByOrganizationID(orgID uuid.UUID) ([]*domain.Event, error) {
	ctx := context.Background()

	events, err := r.client.Event.
		Query().
		Where(event.OrganizationID(orgID)).
		Order(ent.Desc(event.FieldStartTime)).
		All(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get events by organization: %w", err)
	}

	result := make([]*domain.Event, len(events))
	for i, evt := range events {
		result[i] = &domain.Event{
			ID:               evt.ID,
			OrganizationID:   evt.OrganizationID,
			Title:            evt.Title,
			Description:      evt.Description,
			Location:         evt.Location,
			Venue:            evt.Venue,
			StartTime:        evt.StartTime,
			EndTime:          evt.EndTime,
			TotalTickets:     evt.TotalTickets,
			AvailableTickets: evt.AvailableTickets,
			TicketPrice:      evt.TicketPrice,
			Currency:         evt.Currency,
			ThumbnailURL:     evt.ThumbnailURL,
			Status:           string(evt.Status),
			IsPublic:         evt.IsPublic,
			CreatedBy:        evt.CreatedBy,
			CreatedAt:        evt.CreatedAt,
			UpdatedAt:        evt.UpdatedAt,
		}
	}

	return result, nil
}

// Update updates an event
func (r *eventRepository) Update(evt *domain.Event) error {
	ctx := context.Background()

	err := r.client.Event.
		UpdateOneID(evt.ID).
		SetTitle(evt.Title).
		SetDescription(evt.Description).
		SetLocation(evt.Location).
		SetVenue(evt.Venue).
		SetStartTime(evt.StartTime).
		SetEndTime(evt.EndTime).
		SetTotalTickets(evt.TotalTickets).
		SetAvailableTickets(evt.AvailableTickets).
		SetTicketPrice(evt.TicketPrice).
		SetCurrency(evt.Currency).
		SetThumbnailURL(evt.ThumbnailURL).
		SetStatus(event.Status(evt.Status)).
		SetIsPublic(evt.IsPublic).
		Exec(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return domain.ErrNotFound
		}
		return fmt.Errorf("failed to update event: %w", err)
	}

	return nil
}

// Delete deletes an event
func (r *eventRepository) Delete(eventID uuid.UUID) error {
	ctx := context.Background()

	err := r.client.Event.
		DeleteOneID(eventID).
		Exec(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return domain.ErrNotFound
		}
		return fmt.Errorf("failed to delete event: %w", err)
	}

	return nil
}

// GetPublicEvents retrieves all public events
func (r *eventRepository) GetPublicEvents() ([]*domain.EventWithOrganization, error) {
	ctx := context.Background()

	events, err := r.client.Event.
		Query().
		Where(event.IsPublic(true)).
		WithOrganization().
		Order(ent.Desc(event.FieldStartTime)).
		All(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get public events: %w", err)
	}

	return r.mapEventsWithOrganization(events), nil
}

// GetEventsByStatus retrieves events by status
func (r *eventRepository) GetEventsByStatus(status string) ([]*domain.EventWithOrganization, error) {
	ctx := context.Background()

	events, err := r.client.Event.
		Query().
		Where(event.StatusEQ(event.Status(status))).
		WithOrganization().
		Order(ent.Desc(event.FieldStartTime)).
		All(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get events by status: %w", err)
	}

	return r.mapEventsWithOrganization(events), nil
}

// GetUpcomingEvents retrieves upcoming events
func (r *eventRepository) GetUpcomingEvents() ([]*domain.EventWithOrganization, error) {
	ctx := context.Background()
	now := time.Now()

	events, err := r.client.Event.
		Query().
		Where(
			event.IsPublic(true),
			event.StartTimeGT(now),
			event.StatusIn(event.StatusPublished, event.StatusOngoing),
		).
		WithOrganization().
		Order(ent.Asc(event.FieldStartTime)).
		All(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get upcoming events: %w", err)
	}

	log.Println(events)

	return r.mapEventsWithOrganization(events), nil
}

// GetPopularEvents retrieves popular upcoming events based on ticket sales
// threshold is the minimum percentage of tickets sold (0.0 to 1.0)
// For example, 0.7 means events with 70% or more tickets sold
func (r *eventRepository) GetPopularEvents(threshold float64) ([]*domain.EventWithOrganization, error) {
	ctx := context.Background()
	now := time.Now()

	// Get all upcoming public events
	events, err := r.client.Event.
		Query().
		Where(
			event.IsPublic(true),
			event.StartTimeGT(now),
			event.StatusIn(event.StatusPublished, event.StatusOngoing),
			event.TotalTicketsGT(0), // Only events with tickets
		).
		WithOrganization().
		All(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get popular events: %w", err)
	}

	// Filter events based on ticket sales percentage
	popularEvents := make([]*ent.Event, 0)
	for _, evt := range events {
		soldTickets := evt.TotalTickets - evt.AvailableTickets
		salesPercentage := float64(soldTickets) / float64(evt.TotalTickets)

		if salesPercentage >= threshold {
			popularEvents = append(popularEvents, evt)
		}
	}

	// Sort by sales percentage (highest first), then by start time
	// We'll sort by sold tickets count as a proxy since we can't add calculated fields
	// You can implement custom sorting if needed

	return r.mapEventsWithOrganization(popularEvents), nil
}

// SearchEvents searches events by keyword in title or description
func (r *eventRepository) SearchEvents(keyword string) ([]*domain.EventWithOrganization, error) {
	ctx := context.Background()

	events, err := r.client.Event.
		Query().
		Where(
			event.Or(
				event.TitleContains(keyword),
				event.DescriptionContains(keyword),
			),
			event.IsPublic(true),
		).
		WithOrganization().
		Order(ent.Desc(event.FieldStartTime)).
		All(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to search events: %w", err)
	}

	return r.mapEventsWithOrganization(events), nil
}

// UpdateAvailableTickets updates the available tickets count
func (r *eventRepository) UpdateAvailableTickets(eventID uuid.UUID, tickets int) error {
	ctx := context.Background()

	err := r.client.Event.
		UpdateOneID(eventID).
		SetAvailableTickets(tickets).
		Exec(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return domain.ErrNotFound
		}
		return fmt.Errorf("failed to update available tickets: %w", err)
	}

	return nil
}

// Helper function to map events with organization
func (r *eventRepository) mapEventsWithOrganization(events []*ent.Event) []*domain.EventWithOrganization {
	result := make([]*domain.EventWithOrganization, len(events))
	for i, evt := range events {
		orgName := ""
		if evt.Edges.Organization != nil {
			orgName = evt.Edges.Organization.Name
		}

		result[i] = &domain.EventWithOrganization{
			Event: domain.Event{
				ID:               evt.ID,
				OrganizationID:   evt.OrganizationID,
				Title:            evt.Title,
				Description:      evt.Description,
				Location:         evt.Location,
				Venue:            evt.Venue,
				StartTime:        evt.StartTime,
				EndTime:          evt.EndTime,
				TotalTickets:     evt.TotalTickets,
				AvailableTickets: evt.AvailableTickets,
				TicketPrice:      evt.TicketPrice,
				Currency:         evt.Currency,
				ThumbnailURL:     evt.ThumbnailURL,
				Status:           string(evt.Status),
				IsPublic:         evt.IsPublic,
				CreatedBy:        evt.CreatedBy,
				CreatedAt:        evt.CreatedAt,
				UpdatedAt:        evt.UpdatedAt,
			},
			OrganizationName: orgName,
		}
	}
	return result
}
