package domain

import "errors"

type PollOption struct {
	Base
	PollID          string `gorm:"type:uuid;not null;uniqueIndex:idx_poll_option_poll"`
	FoodID          string `gorm:"type:uuid;not null;uniqueIndex:idx_poll_option_food"`
	PriceAtCreation int64  `gorm:"not null"`
	Votes           []Vote `gorm:"foreignKey:PollOptionID;constraint:OnDelete:CASCADE;"`
}

var (
	ErrPollOptionAlreadyExists = errors.New("poll option already exists for this food in the poll")
	ErrInvalidPollOption       = errors.New("invalid poll option")
)
