package services

import (
	"context"

	"github.com/SirNacou/weeate/backend/internal/features/auth"
	"github.com/SirNacou/weeate/backend/internal/features/polls/domain"
	"github.com/gofrs/uuid/v5"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type CastVoteCommand struct {
	PollID       uuid.UUID `json:"poll_id" format:"uuid" doc:"The ID of the poll to vote in"`
	PollOptionID uuid.UUID `json:"poll_option_id" format:"uuid" doc:"The ID of the poll option to vote for"`
}
type CastVoteCommandHandler struct {
	db *gorm.DB
}

func NewCastVoteCommandHandler(db *gorm.DB) *CastVoteCommandHandler {
	return &CastVoteCommandHandler{
		db: db,
	}
}

func (h *CastVoteCommandHandler) Handle(ctx context.Context, req CastVoteCommand) error {
	user, ok := auth.UserFromContext(ctx)
	if !ok {
		return auth.ErrUserNotFoundInContext
	}

	poll, err := gorm.G[domain.Poll](h.db).
		Where("id = ?", req.PollID).
		Preload("PollOptions", nil).
		First(ctx)
	if err != nil {
		return err
	}

	if poll.ClosedAt != nil {
		return domain.ErrPollAlreadyClosed
	}

	// Check if user already voted for the selected option
	var existingVote domain.Vote
	err = h.db.Where("poll_id = ? AND user_id = ?", req.PollID, user.ID).First(&existingVote).Error
	if err == nil {
		if existingVote.PollOptionID == req.PollOptionID {
			return domain.ErrUserAlreadyVoted
		}
	} else if err != gorm.ErrRecordNotFound {
		return err
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
