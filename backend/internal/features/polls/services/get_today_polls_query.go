package services

import (
	"context"
	"time"

	"github.com/SirNacou/weeate/backend/internal/features/auth"
	food_domain "github.com/SirNacou/weeate/backend/internal/features/foods/domain"
	"github.com/SirNacou/weeate/backend/internal/features/polls/domain"
	"github.com/SirNacou/weeate/backend/internal/features/polls/infrastructure/data"
	"github.com/samber/lo"
	"gorm.io/gorm"
)

type (
	GetActivePollsQuery         struct{}
	GetActivePollsQueryResponse struct {
		ID                string              `json:"id"`
		Creator           auth.UserProfile    `json:"creator"`
		ScheduledClosesAt time.Time           `json:"scheduled_closes_at"`
		FinalTotalPrice   int64               `json:"final_total_price"`
		Strategy          domain.PollStrategy `json:"strategy"`
		ClosedAt          *time.Time          `json:"closed_at"`
		PollOptions       []PollOption        `json:"poll_options"`
	}

	PollOption struct {
		ID              string `json:"id"`
		Food            Food   `json:"food"`
		PriceAtCreation int64  `json:"price_at_creation"`
		Votes           []Vote `json:"votes"`
	}

	Food struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}

	Vote struct {
		Voter auth.UserProfile `json:"voter"`
	}
)

type GetActivePollsQueryHandler struct {
	// Add necessary dependencies here, e.g., database connection
	db              *gorm.DB
	supabaseService *auth.SupabaseService
}

func NewGetActivePollsQueryHandler(db *gorm.DB, s *auth.SupabaseService) *GetActivePollsQueryHandler {
	return &GetActivePollsQueryHandler{
		db:              db,
		supabaseService: s,
	}
}

func (h *GetActivePollsQueryHandler) Handle(ctx context.Context, query GetActivePollsQuery) ([]GetActivePollsQueryResponse, error) {
	now := time.Now()
	polls, err := data.GetPollAggregate(h.db).
		Where("created_at >= ? AND created_at < ?", now.Truncate(24*time.Hour), now.Truncate(24*time.Hour).Add(24*time.Hour)).
		Find(ctx)
	if err != nil {
		return nil, err
	}

	// 1. Collect all unique food and user IDs from the polls.
	var (
		foodIDs []string
		userIDs []string
	)
	for _, p := range polls {
		for _, opt := range p.PollOptions {
			foodIDs = append(foodIDs, opt.FoodID.String())
			for _, vote := range opt.Votes {
				userIDs = append(userIDs, vote.UserID)
			}
		}
	}
	uniqueFoodIDs := lo.Uniq(foodIDs)
	uniqueUserIDs := lo.Uniq(userIDs)

	// 2. Fetch all required foods and user profiles in single queries.
	foods, err := gorm.G[food_domain.Food](h.db).Where("id IN (?)", uniqueFoodIDs).Find(ctx)
	if err != nil {
		return nil, err
	}
	userProfiles, err := h.supabaseService.GetUserProfilesByIDs(uniqueUserIDs)
	if err != nil {
		return nil, err
	}

	// 3. Create maps for efficient O(1) lookups.
	foodMap := lo.KeyBy(foods, func(f food_domain.Food) string { return f.ID.String() })
	userProfileMap := lo.KeyBy(userProfiles, func(up auth.UserProfile) string { return up.ID.String() })

	// 4. Assemble the response DTO using the maps for fast data retrieval.
	res := make([]GetActivePollsQueryResponse, 0, len(polls))
	for _, poll := range polls {
		res = append(res, GetActivePollsQueryResponse{
			ID:                poll.ID.String(),
			ScheduledClosesAt: poll.ScheduledClosesAt,
			Strategy:          poll.Strategy,
			PollOptions: lo.Map(poll.PollOptions, func(option domain.PollOption, _ int) PollOption {
				food := foodMap[option.FoodID.String()]
				return PollOption{
					ID: option.ID.String(),
					Food: Food{
						ID:   food.ID.String(),
						Name: food.Name,
					},
					PriceAtCreation: option.PriceAtCreation,
					Votes: lo.Map(option.Votes, func(vote domain.Vote, _ int) Vote {
						return Vote{
							Voter: userProfileMap[vote.UserID],
						}
					}),
				}
			}),
		})
	}
	return res, nil
}
