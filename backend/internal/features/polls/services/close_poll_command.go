package services

import (
	"context"
	"database/sql"
	"errors"

	"github.com/SirNacou/weeate/backend/internal/common/infrastructure/bus"
	"github.com/SirNacou/weeate/backend/internal/common/infrastructure/centrifugo"
	"github.com/SirNacou/weeate/backend/internal/features/polls/domain"
	"github.com/gofrs/uuid/v5"
	"gorm.io/gorm"
)

type ClosePollCommand struct {
	PollID uuid.UUID `json:"poll_id" format:"uuid" doc:"The ID of the poll to be closed"`
}

type ClosePollCommandHandler struct {
	db         *gorm.DB
	bus        *bus.Bus
	centrifugo centrifugo.WebsocketClient
}

func NewClosePollCommandHandler(db *gorm.DB, bus *bus.Bus, centrifugo centrifugo.WebsocketClient) *ClosePollCommandHandler {
	return &ClosePollCommandHandler{
		db:         db,
		bus:        bus,
		centrifugo: centrifugo,
	}
}

func (c *ClosePollCommandHandler) Handle(ctx context.Context, req ClosePollCommand) error {
	if req.PollID.IsNil() {
		return domain.ErrInvalidPollID
	}

	poll, err := gorm.G[domain.Poll](c.db).
		Preload("PollOptions.Votes", nil).
		Where("id = ?", req.PollID).
		First(ctx)
	if err != nil {
		return err
	}

	if poll.ClosedAt != nil {
		return domain.ErrPollAlreadyClosed
	}

	closedPoll, event := poll.Close()

	return c.db.Transaction(func(tx *gorm.DB) error {
		db, ok := tx.Statement.ConnPool.(*sql.Tx)
		if !ok {
			return errors.New("failed to get sql transaction")
		}

		if _, err := gorm.G[domain.Poll](tx).Updates(ctx, *closedPoll); err != nil {
			return err
		}

		publisher, err := c.bus.NewSqlPublisher(db)
		if err != nil {
			return err
		}

		if err := publisher.Publish(ctx, event); err != nil {
			return err
		}

		return nil
	})
}
