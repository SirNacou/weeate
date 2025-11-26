package services

import (
	"context"
	"log/slog"
	"time"

	"github.com/SirNacou/weeate/backend/internal/features/auth"
	food_domain "github.com/SirNacou/weeate/backend/internal/features/foods/domain"
	"github.com/SirNacou/weeate/backend/internal/features/polls/domain"
	"github.com/SirNacou/weeate/backend/internal/features/polls/infrastructure/data"
	"github.com/samber/lo"
	"gorm.io/gorm"
)

type (
	GetTodayPollsQuery         struct{}
	GetTodayPollsQueryResponse struct {
		ID                string              `json:"id"`
		Creator           auth.UserProfile    `json:"creator"`
		ScheduledClosesAt time.Time           `json:"scheduled_closes_at"`
		Strategy          domain.PollStrategy `json:"strategy" enum:"ORDER_CONSENSUS_ITEM,ORDER_PERSONAL_CHOICE"`
		ClosedAt          *time.Time          `json:"closed_at"`
		PollOptions       []PollOption        `json:"poll_options"`
		OrderDate         time.Time           `json:"order_date" format:"date-time"`
	}

	PollOption struct {
		ID              string `json:"id"`
		Food            Food   `json:"food"`
		PriceAtCreation int64  `json:"price_at_creation"`
		Votes           []Vote `json:"votes"`
	}

	Food struct {
		ID       string  `json:"id"`
		Name     string  `json:"name"`
		ImageURL *string `json:"image_url"`
	}

	Vote struct {
		Voter auth.UserProfile `json:"voter"`
	}
)

type GetTodayPollsQueryHandler struct {
	// Add necessary dependencies here, e.g., database connection
	db              *gorm.DB
	supabaseService *auth.SupabaseService
}

func NewGetTodayPollsQueryHandler(db *gorm.DB, s *auth.SupabaseService) *GetTodayPollsQueryHandler {
	return &GetTodayPollsQueryHandler{
		db:              db,
		supabaseService: s,
	}
}

func (h *GetTodayPollsQueryHandler) Handle(ctx context.Context, query GetTodayPollsQuery) ([]GetTodayPollsQueryResponse, error) {
	now := time.Now()
	polls, err := data.GetPollAggregate(h.db).
		Where("polls.created_at >= ? AND polls.created_at < ?", now.Truncate(24*time.Hour), now.Truncate(24*time.Hour).Add(24*time.Hour)).
		Find(ctx)
	if err != nil {
		return nil, err
	}
	if len(polls) == 0 {
		slog.InfoContext(ctx, "polls is empty", "time", now,
			"start", now.Truncate(24*time.Hour),
			"end", now.Truncate(24*time.Hour).Add(24*time.Hour))
	}

	// 1. Collect all unique food and user IDs from the polls.
	var (
		foodIDs []string
		userIDs []string
	)
	for _, p := range polls {
		userIDs = append(userIDs, p.BuyerID)
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
	var foods []food_domain.Food
	if len(uniqueFoodIDs) > 0 {
		if err := h.db.WithContext(ctx).Where("id IN ?", uniqueFoodIDs).Find(&foods).Error; err != nil {
			return nil, err
		}
	}

	var userProfiles []auth.UserProfile
	if len(uniqueUserIDs) > 0 {
		profiles, err := h.supabaseService.GetUserProfilesByIDs(uniqueUserIDs...)
		if err != nil {
			return nil, err
		}
		userProfiles = profiles
	}

	// 3. Create maps for efficient O(1) lookups.
	foodMap := lo.KeyBy(foods, func(f food_domain.Food) string { return f.ID.String() })
	userProfileMap := lo.KeyBy(userProfiles, func(up auth.UserProfile) string { return up.ID.String() })

	// 4. Assemble the response DTO using the maps for fast data retrieval.
	res := make([]GetTodayPollsQueryResponse, 0, len(polls))
	for _, poll := range polls {
		res = append(res, GetTodayPollsQueryResponse{
			ID:                poll.ID.String(),
			Creator:           userProfileMap[poll.BuyerID],
			ScheduledClosesAt: poll.ScheduledClosesAt,
			Strategy:          poll.Strategy,
			ClosedAt:          poll.ClosedAt,
			OrderDate:         poll.OrderDate,
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
