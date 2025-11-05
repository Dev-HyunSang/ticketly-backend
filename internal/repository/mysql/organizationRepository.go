package mysql

import (
	"context"
	"fmt"

	"github.com/dev-hyunsang/ticketly-backend/internal/domain"
	"github.com/dev-hyunsang/ticketly-backend/lib/ent"
	"github.com/dev-hyunsang/ticketly-backend/lib/ent/organization"
	"github.com/dev-hyunsang/ticketly-backend/lib/ent/organizationmember"
	"github.com/google/uuid"
)

type organizationRepository struct {
	client *ent.Client
}

func NewOrganizationRepository(client *ent.Client) domain.OrganizationRepository {
	return &organizationRepository{
		client: client,
	}
}

// Create creates a new organization and adds the owner as an admin member
func (r *organizationRepository) Create(org *domain.Organization) (*domain.Organization, error) {
	ctx := context.Background()

	// Start transaction
	tx, err := r.client.Tx(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to start transaction: %w", err)
	}

	// Create organization
	createdOrg, err := tx.Organization.
		Create().
		SetID(org.ID).
		SetName(org.Name).
		SetNillableDescription(&org.Description).
		SetNillableLogoURL(&org.LogoURL).
		SetOwnerID(org.OwnerID).
		SetIsActive(org.IsActive).
		Save(ctx)
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to create organization: %w", err)
	}

	// Add owner as admin member
	_, err = tx.OrganizationMember.
		Create().
		SetOrganizationID(createdOrg.ID).
		SetUserID(org.OwnerID).
		SetRole(organizationmember.RoleAdmin).
		Save(ctx)
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to add owner as admin: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return &domain.Organization{
		ID:          createdOrg.ID,
		Name:        createdOrg.Name,
		Description: createdOrg.Description,
		LogoURL:     createdOrg.LogoURL,
		OwnerID:     createdOrg.OwnerID,
		IsActive:    createdOrg.IsActive,
		CreatedAt:   createdOrg.CreatedAt,
		UpdatedAt:   createdOrg.UpdatedAt,
	}, nil
}

// GetByID retrieves an organization by ID
func (r *organizationRepository) GetByID(orgID uuid.UUID) (*domain.Organization, error) {
	ctx := context.Background()

	org, err := r.client.Organization.
		Query().
		Where(organization.ID(orgID)).
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("failed to get organization: %w", err)
	}

	return &domain.Organization{
		ID:          org.ID,
		Name:        org.Name,
		Description: org.Description,
		LogoURL:     org.LogoURL,
		OwnerID:     org.OwnerID,
		IsActive:    org.IsActive,
		CreatedAt:   org.CreatedAt,
		UpdatedAt:   org.UpdatedAt,
	}, nil
}

// GetByOwnerID retrieves all organizations owned by a user
func (r *organizationRepository) GetByOwnerID(ownerID uuid.UUID) ([]*domain.Organization, error) {
	ctx := context.Background()

	orgs, err := r.client.Organization.
		Query().
		Where(organization.OwnerID(ownerID)).
		All(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get organizations by owner: %w", err)
	}

	result := make([]*domain.Organization, len(orgs))
	for i, org := range orgs {
		result[i] = &domain.Organization{
			ID:          org.ID,
			Name:        org.Name,
			Description: org.Description,
			LogoURL:     org.LogoURL,
			OwnerID:     org.OwnerID,
			IsActive:    org.IsActive,
			CreatedAt:   org.CreatedAt,
			UpdatedAt:   org.UpdatedAt,
		}
	}

	return result, nil
}

// Update updates an organization
func (r *organizationRepository) Update(org *domain.Organization) error {
	ctx := context.Background()

	err := r.client.Organization.
		UpdateOneID(org.ID).
		SetName(org.Name).
		SetDescription(org.Description).
		SetLogoURL(org.LogoURL).
		SetIsActive(org.IsActive).
		Exec(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return domain.ErrNotFound
		}
		return fmt.Errorf("failed to update organization: %w", err)
	}

	return nil
}

// Delete deletes an organization
func (r *organizationRepository) Delete(orgID uuid.UUID) error {
	ctx := context.Background()

	err := r.client.Organization.
		DeleteOneID(orgID).
		Exec(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return domain.ErrNotFound
		}
		return fmt.Errorf("failed to delete organization: %w", err)
	}

	return nil
}

// AddMember adds a member to an organization
func (r *organizationRepository) AddMember(member *domain.OrganizationMember) error {
	ctx := context.Background()

	_, err := r.client.OrganizationMember.
		Create().
		SetOrganizationID(member.OrganizationID).
		SetUserID(member.UserID).
		SetRole(organizationmember.Role(member.Role)).
		Save(ctx)
	if err != nil {
		return fmt.Errorf("failed to add member: %w", err)
	}

	return nil
}

// RemoveMember removes a member from an organization
func (r *organizationRepository) RemoveMember(orgID, userID uuid.UUID) error {
	ctx := context.Background()

	_, err := r.client.OrganizationMember.
		Delete().
		Where(
			organizationmember.OrganizationID(orgID),
			organizationmember.UserID(userID),
		).
		Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to remove member: %w", err)
	}

	return nil
}

// GetMember retrieves a specific member
func (r *organizationRepository) GetMember(orgID, userID uuid.UUID) (*domain.OrganizationMember, error) {
	ctx := context.Background()

	member, err := r.client.OrganizationMember.
		Query().
		Where(
			organizationmember.OrganizationID(orgID),
			organizationmember.UserID(userID),
		).
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("failed to get member: %w", err)
	}

	return &domain.OrganizationMember{
		ID:             member.ID,
		OrganizationID: member.OrganizationID,
		UserID:         member.UserID,
		Role:           string(member.Role),
		JoinedAt:       member.JoinedAt,
		CreatedAt:      member.CreatedAt,
		UpdatedAt:      member.UpdatedAt,
	}, nil
}

// GetMembersByOrgID retrieves all members of an organization
func (r *organizationRepository) GetMembersByOrgID(orgID uuid.UUID) ([]*domain.OrganizationMember, error) {
	ctx := context.Background()

	members, err := r.client.OrganizationMember.
		Query().
		Where(organizationmember.OrganizationID(orgID)).
		All(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get members: %w", err)
	}

	result := make([]*domain.OrganizationMember, len(members))
	for i, member := range members {
		result[i] = &domain.OrganizationMember{
			ID:             member.ID,
			OrganizationID: member.OrganizationID,
			UserID:         member.UserID,
			Role:           string(member.Role),
			JoinedAt:       member.JoinedAt,
			CreatedAt:      member.CreatedAt,
			UpdatedAt:      member.UpdatedAt,
		}
	}

	return result, nil
}

// UpdateMemberRole updates a member's role
func (r *organizationRepository) UpdateMemberRole(orgID, userID uuid.UUID, role string) error {
	ctx := context.Background()

	_, err := r.client.OrganizationMember.
		Update().
		Where(
			organizationmember.OrganizationID(orgID),
			organizationmember.UserID(userID),
		).
		SetRole(organizationmember.Role(role)).
		Save(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return domain.ErrNotFound
		}
		return fmt.Errorf("failed to update member role: %w", err)
	}

	return nil
}

// IsUserAdmin checks if a user is an admin of an organization
func (r *organizationRepository) IsUserAdmin(orgID, userID uuid.UUID) (bool, error) {
	ctx := context.Background()

	count, err := r.client.OrganizationMember.
		Query().
		Where(
			organizationmember.OrganizationID(orgID),
			organizationmember.UserID(userID),
			organizationmember.RoleEQ(organizationmember.RoleAdmin),
		).
		Count(ctx)
	if err != nil {
		return false, fmt.Errorf("failed to check admin status: %w", err)
	}

	return count > 0, nil
}

// IsUserMember checks if a user is a member of an organization
func (r *organizationRepository) IsUserMember(orgID, userID uuid.UUID) (bool, error) {
	ctx := context.Background()

	count, err := r.client.OrganizationMember.
		Query().
		Where(
			organizationmember.OrganizationID(orgID),
			organizationmember.UserID(userID),
		).
		Count(ctx)
	if err != nil {
		return false, fmt.Errorf("failed to check member status: %w", err)
	}

	return count > 0, nil
}

// GetUserOrganizations retrieves all organizations a user belongs to with their role
func (r *organizationRepository) GetUserOrganizations(userID uuid.UUID) ([]*domain.OrganizationWithRole, error) {
	ctx := context.Background()

	members, err := r.client.OrganizationMember.
		Query().
		Where(organizationmember.UserID(userID)).
		WithOrganization().
		All(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get user organizations: %w", err)
	}

	result := make([]*domain.OrganizationWithRole, len(members))
	for i, member := range members {
		org := member.Edges.Organization
		result[i] = &domain.OrganizationWithRole{
			Organization: domain.Organization{
				ID:          org.ID,
				Name:        org.Name,
				Description: org.Description,
				LogoURL:     org.LogoURL,
				OwnerID:     org.OwnerID,
				IsActive:    org.IsActive,
				CreatedAt:   org.CreatedAt,
				UpdatedAt:   org.UpdatedAt,
			},
			UserRole: string(member.Role),
		}
	}

	return result, nil
}
