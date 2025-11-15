package domain

import "github.com/gofrs/uuid/v5"

type Vote struct {
	Base
	PollID       uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:idx_poll_user_vote"`
	UserID       uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:idx_poll_user_vote"`
	PollOptionID uuid.UUID `gorm:"type:uuid;not null;index"`
}
