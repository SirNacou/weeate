package services

import (
	"context"

	"github.com/SirNacou/weeate/backend/internal/common/infrastructure/supabase"
	"github.com/SirNacou/weeate/backend/internal/features/auth/domain"
)

type SetAvatarRequest struct {
	AvatarURL string `json:"avatarUrl" doc:"The URL of the newly set avatar"`
}

// SupabaseServiceInterface defines the interface for Supabase operations
type SupabaseServiceInterface interface {
	UpdateUserProfile(userID string, avatarURL string) error
}

type SetAvatarHandler struct {
	supabase SupabaseServiceInterface
}

func NewSetAvatarHandler(supabase SupabaseServiceInterface) *SetAvatarHandler {
	return &SetAvatarHandler{
		supabase: supabase,
	}
}

// NewSetAvatarHandlerWithService is a convenience constructor for production use
func NewSetAvatarHandlerWithService(supabase *supabase.SupabaseService) *SetAvatarHandler {
	return NewSetAvatarHandler(supabase)
}

func (h *SetAvatarHandler) Handle(ctx context.Context, req *SetAvatarRequest) error {
	user, ok := domain.UserFromContext(ctx)
	if !ok {
		return domain.ErrUserNotFoundInContext
	}

	err := h.supabase.UpdateUserProfile(user.ID, req.AvatarURL)
	if err != nil {
		return err
	}

	return nil
}
