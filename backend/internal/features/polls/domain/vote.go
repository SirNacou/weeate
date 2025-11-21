package domain

import (
	"github.com/SirNacou/weeate/backend/internal/common/domain"
	"github.com/gofrs/uuid/v5"
)

type Vote struct {
	domain.Base
	PollID       uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:idx_poll_user_vote"`
	UserID       string    `gorm:"not null;uniqueIndex:idx_poll_user_vote"`
	PollOptionID uuid.UUID `gorm:"type:uuid;not null;index"`
}
