package events

import (
	"time"

	"github.com/gofrs/uuid/v5"
)

type PollClosedEvent struct {
	PollID    uuid.UUID
	BuyerID   string
	OrderDate time.Time
	Strategy  string
	ClosedAt  time.Time
	Results   []OptionResult
}

func (e PollClosedEvent) Name() string {
	return "PollClosedEvent"
}

type OptionResult struct {
	FoodID          uuid.UUID
	PriceAtCreation int64
	Votes           []VoteResult
}

type VoteResult struct {
	UserID   string
	Quantity int
}

type PollStrategy string

const (
	OrderMultiple  PollStrategy = "ORDER_MULTIPLE_ITEMS"
	OrderConsensus PollStrategy = "ORDER_CONSENSUS_ITEM"
)
