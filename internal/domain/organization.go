package domain

import (
	"time"

	"github.com/google/uuid"
)

type Organization struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	LogoURL     string    `json:"logo_url,omitempty"`
	Category    string    `json:"category,omitempty"`
	OwnerID     uuid.UUID `json:"owner_id"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type OrganizationMember struct {
	ID             uuid.UUID `json:"id"`
	OrganizationID uuid.UUID `json:"organization_id"`
	UserID         uuid.UUID `json:"user_id"`
	Role           string    `json:"role"` // "admin" or "member"
	JoinedAt       time.Time `json:"joined_at"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type OrganizationWithRole struct {
	Organization
	UserRole string `json:"user_role"`
}

// OrganizationRepository defines the interface for organization data access
type OrganizationRepository interface {
	// Organization CRUD
	Create(org *Organization) (*Organization, error)
	GetByID(orgID uuid.UUID) (*Organization, error)
	GetByOwnerID(ownerID uuid.UUID) ([]*Organization, error)
	Update(org *Organization) error
	Delete(orgID uuid.UUID) error

	// Member management
	AddMember(member *OrganizationMember) error
	RemoveMember(orgID, userID uuid.UUID) error
	GetMember(orgID, userID uuid.UUID) (*OrganizationMember, error)
	GetMembersByOrgID(orgID uuid.UUID) ([]*OrganizationMember, error)
	UpdateMemberRole(orgID, userID uuid.UUID, role string) error

	// Role checking
	IsUserAdmin(orgID, userID uuid.UUID) (bool, error)
	IsUserMember(orgID, userID uuid.UUID) (bool, error)
	GetUserOrganizations(userID uuid.UUID) ([]*OrganizationWithRole, error)
}
