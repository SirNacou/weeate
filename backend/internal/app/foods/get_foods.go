package foods

import (
	"context"
	"encoding/json"

	api_auth "github.com/SirNacou/weeate/backend/internal/api/auth"
	domain "github.com/SirNacou/weeate/backend/internal/domain"
	"github.com/gofrs/uuid/v5"
	"github.com/rs/zerolog"
	"github.com/samber/lo"
	"github.com/supabase-community/supabase-go"
)

type GetFoodsQuery struct{}

type GetFoodsQueryResult struct {
	ID          uuid.UUID   `json:"id"`
	Name        string      `json:"name"`
	ImageURL    string      `json:"image_url"`
	Description string      `json:"description"`
	Price       int64       `json:"price"`
	User        UserProfile `json:"user"`
}

type UserProfile struct {
	ID          uuid.UUID `json:"id"`
	AvatarURL   string    `json:"avatar_url"`
	DisplayName string    `json:"display_name"`
}
type GetFoodsQueryHandler struct {
	foodRepo       domain.FoodRepository
	supabaseClient *supabase.Client
}

func NewGetFoodsQueryHandler(foodRepo domain.FoodRepository, client *supabase.Client) GetFoodsQueryHandler {
	return GetFoodsQueryHandler{
		foodRepo:       foodRepo,
		supabaseClient: client,
	}
}

func (h *GetFoodsQueryHandler) Handle(ctx context.Context, query GetFoodsQuery) ([]GetFoodsQueryResult, error) {
	res, err := h.foodRepo.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	uniqueUserIDs := lo.UniqMap(res, func(x domain.Food, i int) string { return x.UserID.String() })

	// Fetch user profiles from Supabase
	userProfilesMap := make(map[uuid.UUID]api_auth.UserProfile)
	if len(uniqueUserIDs) > 0 {
		userProfilesBytes, _, err := h.supabaseClient.From("user_profiles").
			Select("*", "", false).
			In("id", uniqueUserIDs).
			Execute()
		if err != nil {
			zerolog.Ctx(ctx).Error().Err(err).Msg("Failed to fetch user profiles from Supabase")
			return nil, err
		}
		var userProfiles []api_auth.UserProfile
		if err := json.Unmarshal(userProfilesBytes, &userProfiles); err == nil {
			for _, profile := range userProfiles {
				userProfilesMap[profile.ID] = profile
			}
		}
	}

	results := make([]GetFoodsQueryResult, 0, len(res))
	for _, food := range res {
		result := GetFoodsQueryResult{
			ID:          food.ID,
			Name:        food.Name,
			ImageURL:    food.ImageURL,
			Description: food.Description,
			Price:       food.Price,
			User: UserProfile{
				ID: food.UserID,
			},
		}

		// Populate user details if available
		if profile, exists := userProfilesMap[food.UserID]; exists {
			result.User.DisplayName = profile.DisplayName
			result.User.AvatarURL = profile.AvatarURL
		}

		results = append(results, result)
	}

	return results, nil
}
