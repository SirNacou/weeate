package services

import (
	"context"

	"github.com/SirNacou/weeate/backend/internal/features/polls/domain"
	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"github.com/gofrs/uuid/v5"
	"gorm.io/gorm"
)

type ClosePollCommand struct {
	PollID uuid.UUID `json:"poll_id" format:"uuid" doc:"The ID of the poll to be closed"`
}

type ClosePollCommandHandler struct {
	db       *gorm.DB
	eventBus *cqrs.EventBus
}

func NewClosePollCommandHandler(db *gorm.DB, eventBus *cqrs.EventBus) *ClosePollCommandHandler {
	return &ClosePollCommandHandler{
		db:       db,
		eventBus: eventBus,
	}
}

func (c *ClosePollCommandHandler) Handle(ctx context.Context, req ClosePollCommand) error {
	if req.PollID.IsNil() {
		return domain.ErrInvalidPollID
	}

	poll, err := gorm.G[domain.Poll](c.db).
		Where("id = ?", req.PollID).
		First(ctx)
	if err != nil {
		return err
	}

	if poll.ClosedAt != nil {
		return domain.ErrPollAlreadyClosed
	}

	poll.Close()

	return c.db.Transaction(func(tx *gorm.DB) error {
		if _, err := gorm.G[domain.Poll](tx).Updates(ctx, poll); err != nil {
			return err
		}

		events := poll.PullEvents()
		for _, event := range events {
			if err := c.eventBus.Publish(ctx, event); err != nil {
				return err
			}
		}

		return nil
	})
}
