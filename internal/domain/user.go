package domain

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID          uuid.UUID `json:"id"`
	FirstName   string    `json:"first_name"`
	LastName    string    `json:"last_name"`
	NickName    string    `json:"nick_name"`
	Birthday    string    `json:"birthday"`
	Email       string    `json:"email"`
	Password    string    `json:"password"`
	PhoneNumber string    `json:"phone_number"`
	IsValid     bool      `json:"is_valid"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type UserRepository interface {
	Save(user *User) (*User, error)
	GetUserByID(userID uuid.UUID) (*User, error)
	GetUserByEmail(userEmail string) (*User, error)
	Update(user *User) error
	DeleteUserByID(userID uuid.UUID) error
}

type userUseCase interface {
	Save(user *User) (*User, error)
	GetUserByID(userID uuid.UUID) (*User, error)
	GetUserByEmail(userEmail string) (*User, error)
	Update(user *User) error
	DeleteUserByID(userID uuid.UUID) error
}
