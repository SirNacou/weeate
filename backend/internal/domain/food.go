package domain

import (
	"context"

	"github.com/gofrs/uuid/v5"
	"gorm.io/gorm"
)

type Food struct {
	ID          uuid.UUID `gorm:"primaryKey;autoIncrement"`
	Name        string    `gorm:"type:varchar(100);not null"`
	ImageURL    string    `gorm:"type:text"`
	Description string    `gorm:"type:text"`
	Price       float64   `gorm:"type:decimal(10,2);not null"`
	UserID      uuid.UUID `gorm:"not null;index"`
}

func NewFood(name, imageUrl, description string, price float64, userID uuid.UUID) (*Food, error) {
	foodID, err := uuid.NewV7()
	if err != nil {
		return nil, err
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

func (f *Food) UpdateDetails(name, imageUrl, description string, price float64) {
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
