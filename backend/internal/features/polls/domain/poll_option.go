package domain

import (
	"errors"

	"github.com/SirNacou/weeate/backend/internal/common/domain"
	"github.com/gofrs/uuid/v5"
)

type PollOption struct {
	domain.Base
	PollID          uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:idx_poll_option_poll"`
	FoodID          uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:idx_poll_option_food"`
	PriceAtCreation int64     `gorm:"not null"`
	Votes           []Vote    `gorm:"foreignKey:PollOptionID;constraint:OnDelete:CASCADE;"`
}

var (
	ErrPollOptionAlreadyExists = errors.New("poll option already exists for this food in the poll")
	ErrInvalidPollOption       = errors.New("invalid poll option")
)
