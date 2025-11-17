package close_poll

import (
	"context"

	"github.com/SirNacou/weeate/backend/internal/features/polls/domain"
	"gorm.io/gorm"
)

type ClosePollCommand struct {
	db *gorm.DB
}

func NewClosePollCommand(db *gorm.DB) *ClosePollCommand {
	return &ClosePollCommand{
		db: db,
	}
}

func (c *ClosePollCommand) Handle(ctx context.Context, req ClosePollRequest) error {
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
	if _, err := gorm.G[domain.Poll](c.db).Updates(ctx, poll); err != nil {
		return err
	}

	return nil
}
