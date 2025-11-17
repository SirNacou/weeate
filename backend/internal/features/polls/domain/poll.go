package domain

import (
	"time"

	"github.com/SirNacou/weeate/backend/internal/common/domain"
	"github.com/gofrs/uuid/v5"
	"github.com/samber/lo"
)

type Poll struct {
	domain.Base
	OrderDate         time.Time    `gorm:"type:date;not null;index"`
	ScheduledClosesAt time.Time    `gorm:"not null"`
	Strategy          PollStrategy `gorm:"type:varchar(100);not null;check:strategy IN ('ORDER_MULTIPLE_ITEMS', 'ORDER_CONSENSUS_ITEM')"`
	ClosedAt          *time.Time   `gorm:"nullable"`
	PollOptions       []PollOption `gorm:"foreignKey:PollID;constraint:OnDelete:CASCADE;"`
}

func NewPoll(orderDate, scheduledClosesAt time.Time, strategy PollStrategy, pollOptions []PollOption) (*Poll, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}

	pollOptions = lo.Map(pollOptions, func(option PollOption, _ int) PollOption {
		option.PollID = id
		return option
	})

	poll := &Poll{
		Base: domain.Base{
			ID: id,
		},
		OrderDate:         orderDate,
		Strategy:          strategy,
		ScheduledClosesAt: scheduledClosesAt,
		PollOptions:       pollOptions,
	}

	return poll, nil
}

func (p *Poll) Close() {
	now := time.Now().UTC()
	p.ClosedAt = &now
}

func (p *Poll) CastVote(optionID uuid.UUID, userID string) (*Vote, error) {
	// Find the poll option the user is voting for.
	newVoteOption, ok := lo.Find(p.PollOptions, func(option PollOption) bool {
		return option.ID == optionID
	})
	if !ok {
		// If the option ID is invalid, do nothing.
		return nil, ErrInvalidPollOption
	}

	// Remove any existing vote from the user across all options.
	for i := range p.PollOptions {
		// Filter out the user's previous vote, if any.
		p.PollOptions[i].Votes = lo.Filter(p.PollOptions[i].Votes, func(vote Vote, _ int) bool {
			return vote.UserID != userID
		})
	}

	// Add the new vote to the correct option.
	newVote := Vote{
		PollID:       p.ID,
		UserID:       userID,
		PollOptionID: newVoteOption.ID,
	}
	for i := range p.PollOptions {
		if p.PollOptions[i].ID == newVoteOption.ID {
			p.PollOptions[i].Votes = append(p.PollOptions[i].Votes, newVote)
			break
		}
	}

	return &newVote, nil
}
