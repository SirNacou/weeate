package auth

import (
	"errors"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string `gorm:"uniqueIndex;size:30;not null"`
	Password string `gorm:"not null;size:100"`
	FullName string `gorm:"not null;size:100"`
}

type UserRepository interface {
	WithTx(tx *gorm.DB) UserRepository
	CreateUser(user *User) error
	GetUserByUsername(username string) (*User, error)
	DeleteUser(id uint) error
}

var (
	ErrInvalidUsername       = errors.New("Invalid username")
	ErrInvalidUsernameLength = errors.New("Username length must be between 3 and 30 characters")
	ErrInvalidPassword       = errors.New("Invalid password")
	ErrInvalidPasswordLength = errors.New("Password length must be at least 8 characters")
	ErrInvalidCredentials    = errors.New("Invalid username or password")

	ErrUserNotFound         = errors.New("User not found")
	ErrUsernameAlreadyExist = errors.New("Username already exists")
)
