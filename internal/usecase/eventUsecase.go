package usecase

import (
	"errors"
	"time"

	"github.com/dev-hyunsang/ticketly-backend/internal/domain"
	"github.com/google/uuid"
)

type EventUseCase interface {
	// Event management
	CreateEvent(orgID, userID uuid.UUID, req CreateEventRequest) (*domain.Event, error)
	GetEvent(eventID uuid.UUID) (*domain.Event, error)
	GetOrganizationEvents(orgID uuid.UUID) ([]*domain.Event, error)
	UpdateEvent(eventID, userID uuid.UUID, req UpdateEventRequest) error
	DeleteEvent(eventID, userID uuid.UUID) error

	// Public queries
	GetPublicEvents() ([]*domain.EventWithOrganization, error)
	GetUpcomingEvents() ([]*domain.EventWithOrganization, error)
	SearchEvents(keyword string) ([]*domain.EventWithOrganization, error)

	// Ticket management
	ReserveTickets(eventID uuid.UUID, quantity int) error
	ReleaseTickets(eventID uuid.UUID, quantity int) error
}

type CreateEventRequest struct {
	Title            string    `json:"title"`
	Description      string    `json:"description"`
	Location         string    `json:"location"`
	Venue            string    `json:"venue"`
	StartTime        time.Time `json:"start_time"`
	EndTime          time.Time `json:"end_time"`
	TotalTickets     int       `json:"total_tickets"`
	TicketPrice      float64   `json:"ticket_price"`
	Currency         string    `json:"currency"`
	ThumbnailURL     string    `json:"thumbnail_url"`
	IsPublic         bool      `json:"is_public"`
}

type UpdateEventRequest struct {
	Title            string    `json:"title"`
	Description      string    `json:"description"`
	Location         string    `json:"location"`
	Venue            string    `json:"venue"`
	StartTime        time.Time `json:"start_time"`
	EndTime          time.Time `json:"end_time"`
	TotalTickets     int       `json:"total_tickets"`
	TicketPrice      float64   `json:"ticket_price"`
	Currency         string    `json:"currency"`
	ThumbnailURL     string    `json:"thumbnail_url"`
	Status           string    `json:"status"`
	IsPublic         bool      `json:"is_public"`
}

type eventUseCase struct {
	eventRepo domain.EventRepository
	orgRepo   domain.OrganizationRepository
}

func NewEventUseCase(eventRepo domain.EventRepository, orgRepo domain.OrganizationRepository) EventUseCase {
	return &eventUseCase{
		eventRepo: eventRepo,
		orgRepo:   orgRepo,
	}
}

// CreateEvent creates a new event (admin only)
func (uc *eventUseCase) CreateEvent(orgID, userID uuid.UUID, req CreateEventRequest) (*domain.Event, error) {
	// Check if user is admin of the organization
	isAdmin, err := uc.orgRepo.IsUserAdmin(orgID, userID)
	if err != nil {
		return nil, err
	}
	if !isAdmin {
		return nil, errors.New("permission denied: only admins can create events")
	}

	// Validate required fields
	if req.Title == "" {
		return nil, errors.New("event title is required")
	}
	if req.StartTime.IsZero() {
		return nil, errors.New("start time is required")
	}
	if req.EndTime.IsZero() {
		return nil, errors.New("end time is required")
	}
	if req.StartTime.After(req.EndTime) {
		return nil, errors.New("start time must be before end time")
	}
	if req.TotalTickets < 0 {
		return nil, errors.New("total tickets must be non-negative")
	}
	if req.TicketPrice < 0 {
		return nil, errors.New("ticket price must be non-negative")
	}

	// Set default currency
	if req.Currency == "" {
		req.Currency = "KRW"
	}

	event := &domain.Event{
		ID:               uuid.New(),
		OrganizationID:   orgID,
		Title:            req.Title,
		Description:      req.Description,
		Location:         req.Location,
		Venue:            req.Venue,
		StartTime:        req.StartTime,
		EndTime:          req.EndTime,
		TotalTickets:     req.TotalTickets,
		AvailableTickets: req.TotalTickets,
		TicketPrice:      req.TicketPrice,
		Currency:         req.Currency,
		ThumbnailURL:     req.ThumbnailURL,
		Status:           "draft",
		IsPublic:         req.IsPublic,
		CreatedBy:        userID,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	return uc.eventRepo.Create(event)
}

// GetEvent retrieves an event by ID
func (uc *eventUseCase) GetEvent(eventID uuid.UUID) (*domain.Event, error) {
	return uc.eventRepo.GetByID(eventID)
}

// GetOrganizationEvents retrieves all events for an organization
func (uc *eventUseCase) GetOrganizationEvents(orgID uuid.UUID) ([]*domain.Event, error) {
	return uc.eventRepo.GetByOrganizationID(orgID)
}

// UpdateEvent updates an event (admin only)
func (uc *eventUseCase) UpdateEvent(eventID, userID uuid.UUID, req UpdateEventRequest) error {
	// Get event
	event, err := uc.eventRepo.GetByID(eventID)
	if err != nil {
		return err
	}

	// Check if user is admin of the organization
	isAdmin, err := uc.orgRepo.IsUserAdmin(event.OrganizationID, userID)
	if err != nil {
		return err
	}
	if !isAdmin {
		return errors.New("permission denied: only admins can update events")
	}

	// Validate updates
	if req.StartTime.After(req.EndTime) {
		return errors.New("start time must be before end time")
	}
	if req.TotalTickets < 0 {
		return errors.New("total tickets must be non-negative")
	}
	if req.TicketPrice < 0 {
		return errors.New("ticket price must be non-negative")
	}

	// Validate status
	validStatuses := map[string]bool{
		"draft": true, "published": true, "ongoing": true, "completed": true, "cancelled": true,
	}
	if req.Status != "" && !validStatuses[req.Status] {
		return errors.New("invalid status")
	}

	// Update fields
	if req.Title != "" {
		event.Title = req.Title
	}
	event.Description = req.Description
	event.Location = req.Location
	event.Venue = req.Venue
	if !req.StartTime.IsZero() {
		event.StartTime = req.StartTime
	}
	if !req.EndTime.IsZero() {
		event.EndTime = req.EndTime
	}
	if req.TotalTickets > 0 {
		// Adjust available tickets proportionally
		diff := req.TotalTickets - event.TotalTickets
		event.AvailableTickets += diff
		if event.AvailableTickets < 0 {
			event.AvailableTickets = 0
		}
		event.TotalTickets = req.TotalTickets
	}
	event.TicketPrice = req.TicketPrice
	if req.Currency != "" {
		event.Currency = req.Currency
	}
	event.ThumbnailURL = req.ThumbnailURL
	if req.Status != "" {
		event.Status = req.Status
	}
	event.IsPublic = req.IsPublic
	event.UpdatedAt = time.Now()

	return uc.eventRepo.Update(event)
}

// DeleteEvent deletes an event (admin only)
func (uc *eventUseCase) DeleteEvent(eventID, userID uuid.UUID) error {
	// Get event
	event, err := uc.eventRepo.GetByID(eventID)
	if err != nil {
		return err
	}

	// Check if user is admin of the organization
	isAdmin, err := uc.orgRepo.IsUserAdmin(event.OrganizationID, userID)
	if err != nil {
		return err
	}
	if !isAdmin {
		return errors.New("permission denied: only admins can delete events")
	}

	return uc.eventRepo.Delete(eventID)
}

// GetPublicEvents retrieves all public events
func (uc *eventUseCase) GetPublicEvents() ([]*domain.EventWithOrganization, error) {
	return uc.eventRepo.GetPublicEvents()
}

// GetUpcomingEvents retrieves upcoming public events
func (uc *eventUseCase) GetUpcomingEvents() ([]*domain.EventWithOrganization, error) {
	return uc.eventRepo.GetUpcomingEvents()
}

// SearchEvents searches events by keyword
func (uc *eventUseCase) SearchEvents(keyword string) ([]*domain.EventWithOrganization, error) {
	if keyword == "" {
		return nil, errors.New("search keyword is required")
	}
	return uc.eventRepo.SearchEvents(keyword)
}

// ReserveTickets reserves tickets for an event (decreases available tickets)
func (uc *eventUseCase) ReserveTickets(eventID uuid.UUID, quantity int) error {
	if quantity <= 0 {
		return errors.New("quantity must be positive")
	}

	event, err := uc.eventRepo.GetByID(eventID)
	if err != nil {
		return err
	}

	if event.AvailableTickets < quantity {
		return errors.New("not enough tickets available")
	}

	newAvailable := event.AvailableTickets - quantity
	return uc.eventRepo.UpdateAvailableTickets(eventID, newAvailable)
}

// ReleaseTickets releases tickets for an event (increases available tickets)
func (uc *eventUseCase) ReleaseTickets(eventID uuid.UUID, quantity int) error {
	if quantity <= 0 {
		return errors.New("quantity must be positive")
	}

	event, err := uc.eventRepo.GetByID(eventID)
	if err != nil {
		return err
	}

	newAvailable := event.AvailableTickets + quantity
	if newAvailable > event.TotalTickets {
		newAvailable = event.TotalTickets
	}

	return uc.eventRepo.UpdateAvailableTickets(eventID, newAvailable)
}
