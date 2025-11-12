package domain

import (
	"context"
	"errors"

	"github.com/gofrs/uuid/v5"
	"gorm.io/gorm"
)

type Food struct {
	ID          uuid.UUID `gorm:"primaryKey;autoIncrement"`
	Name        string    `gorm:"type:varchar(100);not null"`
	Description string    `gorm:"type:text"`
	Price       int64     `gorm:"not null"`
	UserID      uuid.UUID `gorm:"not null;index"`
	ImageFileID string    `gorm:"type:varchar(255)"`
	ImageURL    string    `gorm:"type:text"`
}

func NewFood(name, image_file_id, imageUrl, description string, price int64, userID uuid.UUID) (*Food, error) {
	foodID, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}

	if name == "" {
		return nil, ErrInvalidName
	}

	if price < 0 {
		return nil, ErrInvalidPrice
	}

	return &Food{
		ID:          foodID,
		Name:        name,
		ImageURL:    imageUrl,
		Description: description,
		Price:       0.0,
		UserID:      userID,
	}, nil
}

func (f *Food) UpdateDetails(name, image_file_id, imageUrl, description string, price int64) {
	f.Name = name
	f.ImageURL = imageUrl
	f.Description = description
	f.Price = price
}

type FoodRepository interface {
	WithTx(tx *gorm.DB) FoodRepository
	FindByID(ctx context.Context, id uuid.UUID) (Food, error)
	FindAllByID(ctx context.Context, ids ...uuid.UUID) ([]Food, error)
	FindAllByUserID(ctx context.Context, userID uuid.UUID) ([]Food, error)
	Create(ctx context.Context, food *Food) error
	Update(ctx context.Context, food *Food) error
	Delete(ctx context.Context, id uuid.UUID) error
}

var (
	ErrInvalidPrice = errors.New("invalid price: must be non-negative")
	ErrInvalidName  = errors.New("invalid name: cannot be empty")
)
