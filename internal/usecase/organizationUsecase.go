package usecase

import (
	"errors"
	"time"

	"github.com/dev-hyunsang/ticketly-backend/internal/domain"
	"github.com/google/uuid"
)

type OrganizationUseCase interface {
	// Organization management
	CreateOrganization(name, description, logoURL string, ownerID uuid.UUID) (*domain.Organization, error)
	GetOrganization(orgID uuid.UUID) (*domain.Organization, error)
	GetUserOrganizations(userID uuid.UUID) ([]*domain.OrganizationWithRole, error)
	UpdateOrganization(orgID uuid.UUID, name, description, logoURL string, userID uuid.UUID) error
	DeleteOrganization(orgID uuid.UUID, userID uuid.UUID) error

	// Member management
	AddMember(orgID, userID, requesterID uuid.UUID, role string) error
	RemoveMember(orgID, userID, requesterID uuid.UUID) error
	GetMembers(orgID uuid.UUID) ([]*domain.OrganizationMember, error)
	UpdateMemberRole(orgID, userID, requesterID uuid.UUID, newRole string) error

	// Permission checks
	CheckAdminPermission(orgID, userID uuid.UUID) error
	CheckMemberPermission(orgID, userID uuid.UUID) error
}

type organizationUseCase struct {
	orgRepo domain.OrganizationRepository
}

func NewOrganizationUseCase(orgRepo domain.OrganizationRepository) OrganizationUseCase {
	return &organizationUseCase{
		orgRepo: orgRepo,
	}
}

// CreateOrganization creates a new organization
func (uc *organizationUseCase) CreateOrganization(name, description, logoURL string, ownerID uuid.UUID) (*domain.Organization, error) {
	if name == "" {
		return nil, errors.New("organization name is required")
	}

	org := &domain.Organization{
		ID:          uuid.New(),
		Name:        name,
		Description: description,
		LogoURL:     logoURL,
		OwnerID:     ownerID,
		IsActive:    true,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	return uc.orgRepo.Create(org)
}

// GetOrganization retrieves an organization by ID
func (uc *organizationUseCase) GetOrganization(orgID uuid.UUID) (*domain.Organization, error) {
	return uc.orgRepo.GetByID(orgID)
}

// GetUserOrganizations retrieves all organizations a user belongs to
func (uc *organizationUseCase) GetUserOrganizations(userID uuid.UUID) ([]*domain.OrganizationWithRole, error) {
	return uc.orgRepo.GetUserOrganizations(userID)
}

// UpdateOrganization updates an organization (admin only)
func (uc *organizationUseCase) UpdateOrganization(orgID uuid.UUID, name, description, logoURL string, userID uuid.UUID) error {
	// Check if user is admin
	if err := uc.CheckAdminPermission(orgID, userID); err != nil {
		return err
	}

	org, err := uc.orgRepo.GetByID(orgID)
	if err != nil {
		return err
	}

	if name != "" {
		org.Name = name
	}
	org.Description = description
	org.LogoURL = logoURL
	org.UpdatedAt = time.Now()

	return uc.orgRepo.Update(org)
}

// DeleteOrganization deletes an organization (owner only)
func (uc *organizationUseCase) DeleteOrganization(orgID uuid.UUID, userID uuid.UUID) error {
	org, err := uc.orgRepo.GetByID(orgID)
	if err != nil {
		return err
	}

	// Only owner can delete organization
	if org.OwnerID != userID {
		return errors.New("only organization owner can delete the organization")
	}

	return uc.orgRepo.Delete(orgID)
}

// AddMember adds a member to an organization (admin only)
func (uc *organizationUseCase) AddMember(orgID, userID, requesterID uuid.UUID, role string) error {
	// Check if requester is admin
	if err := uc.CheckAdminPermission(orgID, requesterID); err != nil {
		return err
	}

	// Validate role
	if role != "admin" && role != "member" {
		return errors.New("invalid role: must be 'admin' or 'member'")
	}

	// Check if user is already a member
	isMember, err := uc.orgRepo.IsUserMember(orgID, userID)
	if err != nil {
		return err
	}
	if isMember {
		return errors.New("user is already a member of this organization")
	}

	member := &domain.OrganizationMember{
		ID:             uuid.New(),
		OrganizationID: orgID,
		UserID:         userID,
		Role:           role,
		JoinedAt:       time.Now(),
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	return uc.orgRepo.AddMember(member)
}

// RemoveMember removes a member from an organization (admin only)
func (uc *organizationUseCase) RemoveMember(orgID, userID, requesterID uuid.UUID) error {
	// Check if requester is admin
	if err := uc.CheckAdminPermission(orgID, requesterID); err != nil {
		return err
	}

	// Get organization to check owner
	org, err := uc.orgRepo.GetByID(orgID)
	if err != nil {
		return err
	}

	// Cannot remove owner
	if org.OwnerID == userID {
		return errors.New("cannot remove organization owner")
	}

	return uc.orgRepo.RemoveMember(orgID, userID)
}

// GetMembers retrieves all members of an organization (member only)
func (uc *organizationUseCase) GetMembers(orgID uuid.UUID) ([]*domain.OrganizationMember, error) {
	return uc.orgRepo.GetMembersByOrgID(orgID)
}

// UpdateMemberRole updates a member's role (admin only)
func (uc *organizationUseCase) UpdateMemberRole(orgID, userID, requesterID uuid.UUID, newRole string) error {
	// Check if requester is admin
	if err := uc.CheckAdminPermission(orgID, requesterID); err != nil {
		return err
	}

	// Validate role
	if newRole != "admin" && newRole != "member" {
		return errors.New("invalid role: must be 'admin' or 'member'")
	}

	// Get organization to check owner
	org, err := uc.orgRepo.GetByID(orgID)
	if err != nil {
		return err
	}

	// Cannot change owner's role
	if org.OwnerID == userID {
		return errors.New("cannot change organization owner's role")
	}

	return uc.orgRepo.UpdateMemberRole(orgID, userID, newRole)
}

// CheckAdminPermission checks if a user is an admin of an organization
func (uc *organizationUseCase) CheckAdminPermission(orgID, userID uuid.UUID) error {
	isAdmin, err := uc.orgRepo.IsUserAdmin(orgID, userID)
	if err != nil {
		return err
	}

	if !isAdmin {
		return errors.New("permission denied: admin role required")
	}

	return nil
}

// CheckMemberPermission checks if a user is a member of an organization
func (uc *organizationUseCase) CheckMemberPermission(orgID, userID uuid.UUID) error {
	isMember, err := uc.orgRepo.IsUserMember(orgID, userID)
	if err != nil {
		return err
	}

	if !isMember {
		return errors.New("permission denied: organization membership required")
	}

	return nil
}
