package services

import (
	"context"
	"database/sql"
	"errors"

	"github.com/SirNacou/weeate/backend/internal/common/infrastructure/bus"
	"github.com/SirNacou/weeate/backend/internal/features/polls/domain"
	"github.com/gofrs/uuid/v5"
	"gorm.io/gorm"
)

type ClosePollCommand struct {
	PollID uuid.UUID `json:"poll_id" format:"uuid" doc:"The ID of the poll to be closed"`
}

type ClosePollCommandHandler struct {
	db  *gorm.DB
	bus *bus.Bus
}

func NewClosePollCommandHandler(db *gorm.DB, bus *bus.Bus) *ClosePollCommandHandler {
	return &ClosePollCommandHandler{
		db:  db,
		bus: bus,
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
		db, ok := tx.Statement.ConnPool.(*sql.Tx)
		if !ok {
			return errors.New("Failed to get sql transaction")
		}

		if _, err := gorm.G[domain.Poll](tx).Updates(ctx, poll); err != nil {
			return err
		}

		events := poll.PullEvents()

		publisher, err := c.bus.NewSqlPublisher(db)
		if err != nil {
			return err
		}

		for _, event := range events {
			if err := publisher.Publish(ctx, event); err != nil {
				return err
			}
		}

		return nil
	})
}
