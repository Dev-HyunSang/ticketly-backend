package mysql

import (
	"context"
	"fmt"

	"github.com/dev-hyunsang/ticketly-backend/internal/domain"
	"github.com/dev-hyunsang/ticketly-backend/lib/ent"
	"github.com/dev-hyunsang/ticketly-backend/lib/ent/user"
	"github.com/google/uuid"
)

type UserRepository struct {
	client *ent.Client
}

func NewUserRepository(client *ent.Client) *UserRepository {
	return &UserRepository{
		client: client,
	}
}

func (r *UserRepository) Save(userDomain *domain.User) (*domain.User, error) {
	ctx := context.Background()

	user, err := r.client.User.
		Create().
		SetID(userDomain.ID).
		SetFirstName(userDomain.FirstName).
		SetLastName(userDomain.LastName).
		SetNickName(userDomain.NickName).
		SetBirthday(userDomain.Birthday).
		SetEmail(userDomain.Email).
		SetPassword(userDomain.Password).
		SetPhoneNumber(userDomain.PhoneNumber).
		SetIsValid(userDomain.IsValid).
		Save(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return &domain.User{
		ID:          user.ID,
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		NickName:    user.NickName,
		Birthday:    user.Birthday,
		Email:       user.Email,
		Password:    user.Password,
		PhoneNumber: user.PhoneNumber,
		IsValid:     user.IsValid,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
	}, nil
}

func (r *UserRepository) GetUserByID(userID uuid.UUID) (*domain.User, error) {
	ctx := context.Background()

	user, err := r.client.User.
		Query().
		Where(user.ID(userID)).
		Only(ctx)

	if err != nil {
		if ent.IsNotFound(err) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("failed to get user by id: %w", err)
	}

	return &domain.User{
		ID:          user.ID,
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		NickName:    user.NickName,
		Birthday:    user.Birthday,
		Email:       user.Email,
		Password:    user.Password,
		PhoneNumber: user.PhoneNumber,
		IsValid:     user.IsValid,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
	}, nil
}

func (r *UserRepository) GetUserByEmail(userEmail string) (*domain.User, error) {
	ctx := context.Background()

	user, err := r.client.User.
		Query().
		Where(user.Email(userEmail)).
		Only(ctx)

	if err != nil {
		if ent.IsNotFound(err) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}

	return &domain.User{
		ID:          user.ID,
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		NickName:    user.NickName,
		Birthday:    user.Birthday,
		Email:       user.Email,
		Password:    user.Password,
		PhoneNumber: user.PhoneNumber,
		IsValid:     user.IsValid,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
	}, nil
}

func (r *UserRepository) Update(userDomain *domain.User) error {
	ctx := context.Background()

	_, err := r.client.User.
		UpdateOneID(userDomain.ID).
		SetFirstName(userDomain.FirstName).
		SetLastName(userDomain.LastName).
		SetNickName(userDomain.NickName).
		SetBirthday(userDomain.Birthday).
		SetPhoneNumber(userDomain.PhoneNumber).
		SetIsValid(userDomain.IsValid).
		Save(ctx)

	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	return nil
}

func (r *UserRepository) DeleteUserByID(userID uuid.UUID) error {
	ctx := context.Background()

	err := r.client.User.
		DeleteOneID(userID).
		Exec(ctx)

	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	return nil
}
