package services

import (
	"context"

	"github.com/SirNacou/weeate/backend/internal/common/infrastructure/centrifugo"
	auth_domain "github.com/SirNacou/weeate/backend/internal/features/auth/domain"
	"github.com/SirNacou/weeate/backend/internal/features/polls/domain"
	"github.com/gofrs/uuid/v5"
	"gorm.io/gorm"
)

type CastVoteCommand struct {
	PollID       uuid.UUID `json:"poll_id" format:"uuid" doc:"The ID of the poll to vote in"`
	PollOptionID uuid.UUID `json:"poll_option_id" format:"uuid" doc:"The ID of the poll option to vote for"`
}
type CastVoteCommandHandler struct {
	db         *gorm.DB
	centrifugo *centrifugo.CentrifugoClient
}

func NewCastVoteCommandHandler(db *gorm.DB, centrifugo *centrifugo.CentrifugoClient) *CastVoteCommandHandler {
	return &CastVoteCommandHandler{
		db:         db,
		centrifugo: centrifugo,
	}
}

func (h *CastVoteCommandHandler) Handle(ctx context.Context, req CastVoteCommand) error {
	user, ok := auth_domain.UserFromContext(ctx)
	if !ok {
		return auth_domain.ErrUserNotFoundInContext
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
	var existingVotePtr *domain.Vote

	err = h.db.Where("poll_id = ? AND user_id = ?", req.PollID, user.ID).First(&existingVote).Error
	if err == nil {
		existingVotePtr = &existingVote
	} else if err != gorm.ErrRecordNotFound {
		return err
	}

	result, err := poll.CastVote(user.ID, req.PollOptionID, existingVotePtr)
	if err != nil {
		return err
	}

	err = h.db.Transaction(func(tx *gorm.DB) error {
		if result.VoteToDelete != nil {
			if err := tx.Delete(result.VoteToDelete).Error; err != nil {
				return err
			}
		}
		if result.VoteToSave != nil {
			if existingVotePtr != nil {
				if err := tx.Save(result.VoteToSave).Error; err != nil {
					return err
				}
			} else {
				if err := tx.Create(result.VoteToSave).Error; err != nil {
					return err
				}
			}
		}
		return nil
	})
	if err != nil {
		return err
	}

	publishToWebsocket(ctx, h.centrifugo, user, result)

	return nil
}

func publishToWebsocket(ctx context.Context,
	cent *centrifugo.CentrifugoClient,
	user *auth_domain.User,
	result domain.CastVoteResult,
) {
	var voteEvent *centrifugo.PublicPollsData

	if result.VoteToDelete != nil {
		voteEvent = &centrifugo.PublicPollsData{
			Type: centrifugo.VoteRemoved,
			Data: &struct {
				PollID   uuid.UUID `json:"poll_id"`
				UserID   string    `json:"user_id"`
				OptionID uuid.UUID `json:"option_id"`
			}{
				PollID:   result.VoteToDelete.PollID,
				UserID:   result.VoteToDelete.UserID,
				OptionID: result.VoteToDelete.PollOptionID,
			},
		}
	} else if result.OldOptionID != uuid.Nil {
		voteEvent = &centrifugo.PublicPollsData{
			Type: centrifugo.VoteMoved,
			Data: &struct {
				PollID          uuid.UUID `json:"poll_id"`
				UserID          string    `json:"user_id"`
				UserDisplayName string    `json:"user_display_name"`
				UserAvatarURL   string    `json:"user_avatar_url"`
				OldOptionID     uuid.UUID `json:"old_option_id"`
				NewOptionID     uuid.UUID `json:"new_option_id"`
			}{
				PollID:          result.VoteToSave.PollID,
				UserID:          result.VoteToSave.UserID,
				UserDisplayName: user.AppMetadata.DisplayName,
				UserAvatarURL:   user.AppMetadata.AvatarURL,
				OldOptionID:     result.OldOptionID,
				NewOptionID:     result.VoteToSave.PollOptionID,
			},
		}
	} else {
		voteEvent = &centrifugo.PublicPollsData{
			Type: centrifugo.VoteAdded,
			Data: &struct {
				PollID          uuid.UUID `json:"poll_id"`
				UserID          string    `json:"user_id"`
				UserDisplayName string    `json:"user_display_name"`
				UserAvatarURL   string    `json:"user_avatar_url"`
				OptionID        uuid.UUID `json:"option_id"`
			}{
				PollID:          result.VoteToSave.PollID,
				UserID:          result.VoteToSave.UserID,
				UserDisplayName: user.AppMetadata.DisplayName,
				UserAvatarURL:   user.AppMetadata.AvatarURL,
				OptionID:        result.VoteToSave.PollOptionID,
			},
		}
	}

	cent.PublishPublicPolls(ctx, voteEvent)
}
