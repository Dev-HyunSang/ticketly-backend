package domain

import (
	"time"

	"github.com/google/uuid"
)

type Event struct {
	ID               uuid.UUID `json:"id"`
	OrganizationID   uuid.UUID `json:"organization_id"`
	Title            string    `json:"title"`
	Description      string    `json:"description,omitempty"`
	Location         string    `json:"location,omitempty"`
	Venue            string    `json:"venue,omitempty"`
	StartTime        time.Time `json:"start_time"`
	EndTime          time.Time `json:"end_time"`
	TotalTickets     int       `json:"total_tickets"`
	AvailableTickets int       `json:"available_tickets"`
	ParticipantCount int       `json:"participant_count"` // Real-time count based on completed payments
	TicketPrice      float64   `json:"ticket_price"`
	Currency         string    `json:"currency"`
	ThumbnailURL     string    `json:"thumbnail_url,omitempty"`
	Status           string    `json:"status"` // draft, published, ongoing, completed, cancelled
	IsPublic         bool      `json:"is_public"`
	CreatedBy        uuid.UUID `json:"created_by"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

type EventWithOrganization struct {
	Event
	OrganizationName string `json:"organization_name"`
}

// EventRepository defines the interface for event data access
type EventRepository interface {
	// Event CRUD
	Create(event *Event) (*Event, error)
	GetByID(eventID uuid.UUID) (*Event, error)
	GetByOrganizationID(orgID uuid.UUID) ([]*Event, error)
	Update(event *Event) error
	Delete(eventID uuid.UUID) error

	// Queries
	GetPublicEvents() ([]*EventWithOrganization, error)
	GetEventsByStatus(status string) ([]*EventWithOrganization, error)
	GetUpcomingEvents() ([]*EventWithOrganization, error)
	GetPopularEvents(threshold float64) ([]*EventWithOrganization, error)
	SearchEvents(keyword string) ([]*EventWithOrganization, error)

	// Ticket management
	UpdateAvailableTickets(eventID uuid.UUID, tickets int) error
	UpdateParticipantCount(eventID uuid.UUID, count int) error
}
