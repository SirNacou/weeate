package domain

import (
	"time"

	"github.com/SirNacou/weeate/backend/internal/common/domain"
	"github.com/SirNacou/weeate/backend/internal/common/events"
	"github.com/gofrs/uuid/v5"
	"github.com/samber/lo"
)

type Poll struct {
	domain.Base
	BuyerID           string         `gorm:"not null;uniqueIndex:idx_poll_order_date_buyer"`
	OrderDate         time.Time      `gorm:"type:date;not null;uniqueIndex:idx_poll_order_date_buyer"`
	ScheduledClosesAt time.Time      `gorm:"not null"`
	Strategy          PollStrategy   `gorm:"type:varchar(100);not null;check:strategy IN ('ORDER_MULTIPLE_ITEMS', 'ORDER_CONSENSUS_ITEM')"`
	ClosedAt          *time.Time     `gorm:"nullable"`
	PollOptions       []PollOption   `gorm:"foreignKey:PollID;constraint:OnDelete:CASCADE;"`
	events            []domain.Event `gorm:"-"`
}

func NewPoll(orderDate, scheduledClosesAt time.Time, strategy PollStrategy, pollOptions []PollOption, buyerID string) (*Poll, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}

	pollOptions = lo.Map(pollOptions, func(option PollOption, _ int) PollOption {
		option.PollID = id
		return option
	})

	poll := &Poll{
		Base:              domain.NewBase(),
		BuyerID:           buyerID,
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
	p.events = append(p.events, events.PollClosedEvent{
		PollID:    p.ID,
		BuyerID:   p.BuyerID,
		OrderDate: p.OrderDate,
		Strategy:  string(p.Strategy),
		ClosedAt:  now,
		Results: lo.Map(p.PollOptions, func(option PollOption, _ int) events.OptionResult {
			return events.OptionResult{
				FoodID:          option.FoodID,
				PriceAtCreation: option.PriceAtCreation,
				Votes: lo.Map(option.Votes, func(vote Vote, _ int) events.VoteResult {
					return events.VoteResult{
						UserID:   vote.UserID,
						Quantity: 1, // Each vote represents a quantity of 1
					}
				}),
			}
		}),
	})
}

func (p *Poll) PullEvents() []domain.Event {
	events := p.events
	p.events = nil
	return events
}

type CastVoteResult struct {
	VoteToSave   *Vote
	VoteToDelete *Vote
	OldOptionID  uuid.UUID
}

func (p *Poll) CastVote(userID string, optionID uuid.UUID, existingVote *Vote) (result CastVoteResult, err error) {
	// Find the poll option the user is voting for.
	_, ok := lo.Find(p.PollOptions, func(option PollOption) bool {
		return option.ID == optionID
	})
	if !ok {
		// If the option ID is invalid, do nothing.
		return CastVoteResult{}, ErrInvalidPollOption
	}

	if existingVote != nil {
		if existingVote.PollOptionID == optionID {
			// User is voting for the same option -> Toggle off (remove vote)
			return CastVoteResult{
				VoteToDelete: existingVote,
				OldOptionID:  existingVote.PollOptionID,
			}, nil
		}
		// User is changing vote -> Update vote
		oldID := existingVote.PollOptionID
		existingVote.ChangeOption(optionID)
		return CastVoteResult{
			VoteToSave:  existingVote,
			OldOptionID: oldID,
		}, nil
	}

	// User is voting for the first time -> Create vote
	newVote := NewVote(p.ID, userID, optionID)
	return CastVoteResult{
		VoteToSave: &newVote,
	}, nil
}
