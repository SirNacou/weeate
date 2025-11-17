package domain

import (
	"errors"
	"fmt"
	"time"
)

var (
	ErrInvalidPollID           = errors.New("invalid poll ID")
	ErrPollAlreadyClosed       = errors.New("poll is already closed")
	ErrOrderDateInPast         = errors.New("order date cannot be in the past")
	ErrScheduledCloseAtTooSoon = func(minTime time.Time) error {
		return fmt.Errorf("scheduled time must be at least 1 hour in the future (after %v)", minTime.Format(time.Kitchen))
	}
	ErrUserAlreadyVoted = errors.New("user has already voted in this poll")
)
