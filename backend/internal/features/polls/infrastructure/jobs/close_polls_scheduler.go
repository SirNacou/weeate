package jobs

import (
	"context"
	"log/slog"
	"time"

	"github.com/SirNacou/weeate/backend/internal/common/infrastructure/clock"
	"github.com/SirNacou/weeate/backend/internal/features/polls/domain"
	"github.com/SirNacou/weeate/backend/internal/features/polls/services"
	"gorm.io/gorm"
)

type ClosePollScheduler struct {
	db        *gorm.DB
	trigger   chan bool // A channel to "wake up" the sleeper
	closePoll *services.ClosePollCommandHandler
	clock     clock.Clock
}

func NewClosePollScheduler(closePollHandler *services.ClosePollCommandHandler, db *gorm.DB) *ClosePollScheduler {
	return &ClosePollScheduler{
		closePoll: closePollHandler,
		trigger:   make(chan bool, 1),
		db:        db,
		clock:     clock.NewRealClock(),
	}
}

// Call this when your app starts (in a goroutine)
func (s *ClosePollScheduler) Start(ctx context.Context) {
	for {
		// 1. Find the next closest poll time
		var nextPoll domain.Poll
		var err error
		{
			queryCtx, cancel := context.WithTimeout(ctx, time.Second*5)

			nextPoll, err = gorm.G[domain.Poll](s.db).Where("closed_at IS NULL").
				Order("scheduled_closes_at ASC").
				First(queryCtx)

			cancel()
		}

		var duration time.Duration
		if err != nil {
			duration = 1 * time.Hour
		} else {
			duration = time.Until(nextPoll.ScheduledClosesAt)
			if duration < 0 {
				duration = 0
			}
		}

		slog.InfoContext(ctx, "Sleeping for %v...\n", "duration", duration)

		select {
		case <-s.clock.After(duration):
			workerCtx, cancel := context.WithTimeout(ctx, time.Second*30)
			s.closeDuePolls(workerCtx)
			cancel()
		case <-s.trigger:
			slog.InfoContext(ctx, "New poll created! Recalculating schedule...")
		case <-ctx.Done():
			slog.InfoContext(ctx, "Stopping ClosePollScheduler...")
			return
		}
	}
}

// Call this function whenever a user creates/updates a poll
func (s *ClosePollScheduler) TriggerUpdate() {
	// Non-blocking send
	select {
	case s.trigger <- true:
	default:
	}
}

func (s *ClosePollScheduler) closeDuePolls(ctx context.Context) {
	duePolls, err := gorm.G[domain.Poll](s.db).Where("closed_at IS NULL").
		Where("scheduled_closes_at <= ? AND closed_at IS NULL", time.Now()).
		Find(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to fetch due polls", "error", err)
		return
	}

	for _, poll := range duePolls {
		slog.InfoContext(ctx, "Closing poll", "poll_id", poll.ID)
		closeCmd := services.ClosePollCommand{
			PollID: poll.ID,
		}
		err := s.closePoll.Handle(ctx, closeCmd)
		if err != nil {
			slog.ErrorContext(ctx, "Failed to close poll", "poll_id", poll.ID, "error", err)
		} else {
			slog.InfoContext(ctx, "Successfully closed poll", "poll_id", poll.ID)
		}
	}
}
