package services

import (
	"context"
	"log/slog"

	"github.com/SirNacou/weeate/backend/internal/features/auth"
	"github.com/SirNacou/weeate/backend/internal/features/foods/domain"
	"github.com/gofrs/uuid/v5"
	"github.com/samber/lo"
	"gorm.io/gorm"
)

type GetFoodsQuery struct {
	UserID uuid.UUID `query:"user_id,omitempty" format:"uuid" doc:"The ID of the user whose foods to retrieve (optional)"`
}

type GetFoodsQueryResponse struct {
	ID          uuid.UUID        `json:"id"`
	Name        string           `json:"name"`
	ImageURL    string           `json:"image_url"`
	Description string           `json:"description"`
	Price       int64            `json:"price"`
	User        auth.UserProfile `json:"user"`
}

type GetFoodsQueryHandler struct {
	db              *gorm.DB
	supabaseService *auth.SupabaseService
}

func NewGetFoodsQueryHandler(db *gorm.DB, s *auth.SupabaseService) *GetFoodsQueryHandler {
	return &GetFoodsQueryHandler{
		db:              db,
		supabaseService: s,
	}
}

func (h *GetFoodsQueryHandler) Handle(ctx context.Context, query GetFoodsQuery) ([]GetFoodsQueryResponse, error) {
	foods := []domain.Food{}
	err := error(nil)

	foodsTable := gorm.G[domain.Food](h.db)
	if !query.UserID.IsNil() {
		foods, err = foodsTable.Where("user_id = ?", query.UserID).Find(ctx)
	} else {
		foods, err = foodsTable.Find(ctx)
	}

	if err != nil {
		return nil, err
	}

	// 1. Collect unique user IDs from the foods.
	uniqueUserIDs := lo.Uniq(lo.Map(foods, func(f domain.Food, _ int) string {
		return f.UserID.String()
	}))

	// 2. Fetch user profiles from Supabase.
	userProfiles := []auth.UserProfile{}
	if len(uniqueUserIDs) > 0 {
		profiles, err := h.supabaseService.GetUserProfilesByIDs(uniqueUserIDs)
		if err != nil {
			slog.ErrorContext(ctx, "Failed to fetch user profiles from Supabase", "error", err)
			// We can choose to return the error or continue with partial data.
			// Continuing for now, but this could be a hard error.
		}
		userProfiles = profiles
	}

	// 3. Create a map for efficient O(1) user profile lookups.
	userProfileMap := lo.KeyBy(userProfiles, func(p auth.UserProfile) uuid.UUID {
		return p.ID
	})

	// 4. Assemble the response DTO.
	results := make([]GetFoodsQueryResponse, 0, len(foods))
	for _, food := range foods {
		result := GetFoodsQueryResponse{
			ID:          food.ID,
			Name:        food.Name,
			ImageURL:    food.ImageURL,
			Description: food.Description,
			Price:       food.Price,
		}

		// Populate user details if available, otherwise default to a user with only an ID.
		if profile, ok := userProfileMap[food.UserID]; ok {
			result.User = profile
		} else {
			result.User = auth.UserProfile{ID: food.UserID}
		}

		results = append(results, result)
	}

	return results, nil
}
