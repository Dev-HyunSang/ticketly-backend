package usecase

import (
	"github.com/dev-hyunsang/ticketly-backend/internal/domain"
	"github.com/dev-hyunsang/ticketly-backend/internal/repository/mysql"
	"github.com/google/uuid"
)

type userUseCase struct {
	userRepo domain.UserRepository
}

func NewUserUseCase(userRepo *mysql.UserRepository) *userUseCase {
	return &userUseCase{
		userRepo: userRepo,
	}
}

func (uc *userUseCase) Save(user *domain.User) (*domain.User, error) {
	if user.Email == "" || user.Password == "" || user.PhoneNumber == "" {
		return nil, domain.ErrInvalidInput
	}

	if len(user.Email) == 0 || len(user.Password) == 0 || len(user.PhoneNumber) == 0 {
		return nil, domain.ErrInvalidInput
	}

	return uc.userRepo.Save(user)
}

func (uc *userUseCase) GetUserByID(userID uuid.UUID) (*domain.User, error) {
	return uc.userRepo.GetUserByID(userID)
}

func (uc *userUseCase) GetUserByEmail(userEmail string) (*domain.User, error) {
	return uc.userRepo.GetUserByEmail(userEmail)
}

func (uc *userUseCase) Update(user *domain.User) error {
	return uc.userRepo.Update(user)
}

func (uc *userUseCase) DeleteUserByID(userID uuid.UUID) error {
	return uc.userRepo.DeleteUserByID(userID)
}
