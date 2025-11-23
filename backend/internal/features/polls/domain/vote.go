package domain

import (
	"github.com/SirNacou/weeate/backend/internal/common/domain"
	"github.com/gofrs/uuid/v5"
)

type Vote struct {
	domain.Audit
	PollID       uuid.UUID `gorm:"type:uuid;not null;primaryKey"`
	UserID       string    `gorm:"not null;primaryKey"`
	PollOptionID uuid.UUID `gorm:"type:uuid;not null"`
}

func NewVote(pollID uuid.UUID, userID string, pollOptionID uuid.UUID) Vote {
	return Vote{
		PollID:       pollID,
		UserID:       userID,
		PollOptionID: pollOptionID,
	}
}

func (v *Vote) ChangeOption(newOptionID uuid.UUID) {
	v.PollOptionID = newOptionID
}
