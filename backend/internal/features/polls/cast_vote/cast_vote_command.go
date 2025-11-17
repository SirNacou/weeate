package cast_vote

import (
	"context"

	"github.com/SirNacou/weeate/backend/internal/common/service"
	"github.com/SirNacou/weeate/backend/internal/features/auth"
	"github.com/SirNacou/weeate/backend/internal/features/polls/domain"
	"github.com/samber/lo"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type CastVoteCommandHandler struct {
	db *gorm.DB
}

func NewCastVoteCommandHandler(db *gorm.DB) *CastVoteCommandHandler {
	return &CastVoteCommandHandler{
		db: db,
	}
}

func (h *CastVoteCommandHandler) Handle(ctx context.Context, req CastVoteRequest) error {
	user := ctx.Value(service.ContextKeyUser).(auth.User)

	poll, err := gorm.G[domain.Poll](h.db).
		Where("id = ?", req.PollID).
		First(ctx)
	if err != nil {
		return err
	}

	if poll.ClosedAt != nil {
		return domain.ErrPollAlreadyClosed
	}

	pollOption, err := gorm.G[domain.PollOption](h.db).
		Where("id = ? AND poll_id = ?", req.PollOptionID, req.PollID).
		First(ctx)
	if err != nil {
		return err
	}

	if lo.ContainsBy(pollOption.Votes, func(vote domain.Vote) bool {
		return vote.UserID == user.ID
	}) {
		return domain.ErrUserAlreadyVoted
	}

	vote, err := poll.CastVote(req.PollOptionID, user.ID)
	if err != nil {
		return err
	}

	err = h.db.Transaction(func(tx *gorm.DB) error {
		tx.Clauses(
			clause.OnConflict{
				Columns:   []clause.Column{{Name: "poll_id"}, {Name: "user_id"}},
				DoUpdates: clause.Assignments(map[string]interface{}{"poll_option_id": req.PollOptionID}),
			},
		).Create(vote)

		return nil
	})

	return err
}
