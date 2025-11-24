package domain

import (
	"errors"

	"github.com/SirNacou/weeate/backend/internal/common/domain"
	"github.com/gofrs/uuid/v5"
)

type PollOption struct {
	domain.UUID
	domain.Audit
	PollID          uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:idx_poll_option_poll_food"`
	FoodID          uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:idx_poll_option_poll_food"`
	PriceAtCreation int64     `gorm:"not null"`
	Votes           []Vote    `gorm:"foreignKey:PollOptionID;constraint:OnDelete:CASCADE;"`
}

func NewPollOption(pollID, foodID uuid.UUID, priceAtCreation int64) *PollOption {
	return &PollOption{
		UUID:            domain.NewUUID(),
		PollID:          pollID,
		FoodID:          foodID,
		PriceAtCreation: priceAtCreation,
		Votes:           []Vote{},
	}
}

var (
	ErrPollOptionAlreadyExists = errors.New("poll option already exists for this food in the poll")
	ErrInvalidPollOption       = errors.New("invalid poll option")
)
