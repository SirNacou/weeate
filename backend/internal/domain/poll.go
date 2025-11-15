package domain

import (
	"time"

	"github.com/gofrs/uuid/v5"
	"github.com/samber/lo"
)

type Poll struct {
	Base
	PollDate          time.Time    `gorm:"type:date;not null;index"`
	UserTimeZone      string       `gorm:"type:text;not null"`
	ScheduledClosesAt time.Time    `gorm:"not null"`
	FinalTotalPrice   int64        `gorm:"not null;default:0"`
	Strategy          PollStrategy `gorm:"type:varchar(100);not null;check:strategy IN ('ORDER_MULTIPLE_ITEMS', 'ORDER_CONSENSUS_ITEM')"`
	ClosedAt          time.Time    `gorm:""`
	PollOptions       []PollOption `gorm:"foreignKey:PollID;constraint:OnDelete:CASCADE;"`
}

func NewPoll(pollDate, scheduledClosesAt time.Time, userTimeZone string, strategy PollStrategy) (*Poll, error) {
	poll := &Poll{
		PollDate:          pollDate,
		UserTimeZone:      userTimeZone,
		Strategy:          strategy,
		ScheduledClosesAt: scheduledClosesAt,
		FinalTotalPrice:   0,
	}

	return poll, nil
}

func (p *Poll) Close(finalTotalPrice int64) {
	p.ClosedAt = time.Now()
	p.FinalTotalPrice = finalTotalPrice
}

func (p *Poll) Vote(optionID, userID uuid.UUID) error {
	// Find the poll option the user is voting for.
	newVoteOption, ok := lo.Find(p.PollOptions, func(option PollOption) bool {
		return option.ID == optionID
	})
	if !ok {
		// If the option ID is invalid, do nothing.
		return ErrInvalidPollOption
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

	return nil
}
