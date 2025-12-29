package api

import (
	"context"

	"github.com/SirNacou/weeate/backend/internal/common/api"
	"github.com/SirNacou/weeate/backend/internal/common/infrastructure/centrifugo"
	"github.com/SirNacou/weeate/backend/internal/common/infrastructure/configs"
	"github.com/SirNacou/weeate/backend/internal/common/infrastructure/supabase"
	"github.com/SirNacou/weeate/backend/internal/features/auth/domain"
	"github.com/SirNacou/weeate/backend/internal/features/auth/services"
)

type AuthEndpoint struct {
	env                     configs.Config
	getImagekitTokenHandler *services.GetIKTokenQueryHandler
	supabase                *supabase.SupabaseService
	setAvatarHandler        *services.SetAvatarHandler
}

func NewAuthEndpoint(env configs.Config, imagekitHandler *services.GetIKTokenQueryHandler, setAvatarHandler *services.SetAvatarHandler, supabase *supabase.SupabaseService) *AuthEndpoint {
	return &AuthEndpoint{
		env:                     env,
		getImagekitTokenHandler: imagekitHandler,
		setAvatarHandler:        setAvatarHandler,
		supabase:                supabase,
	}
}

type GetCentrifugoJWTResponse struct {
	Token string `json:"token" doc:"The generated JWT token for Centrifugo authentication"`
}

func (e *AuthEndpoint) GetCentrifugoToken(ctx context.Context, _ *struct{}) (*api.Response[GetCentrifugoJWTResponse], error) {
	user, ok := domain.UserFromContext(ctx)
	if !ok {
		return nil, domain.ErrUserNotFoundInContext
	}

	token, err := centrifugo.GenerateCentrifugoToken(user.ID,
		user.Email,
		[]byte(e.env.CENTRIFUGO_HMAC_SECRET),
	)

	return api.NewResponse(&GetCentrifugoJWTResponse{
		Token: token,
	}), err
}

func (e *AuthEndpoint) GetImageKitToken(ctx context.Context, _ *struct{}) (*api.Response[services.ImageKitAuthParametersResponse], error) {
	res, err := e.getImagekitTokenHandler.Handle(ctx, nil)
	if err != nil {
		return nil, err
	}
	return api.NewResponse(res), nil
}

func (e *AuthEndpoint) SetAvatar(ctx context.Context, req *struct{ Body services.SetAvatarRequest }) (*api.Response[struct{}], error) {
	err := e.setAvatarHandler.Handle(ctx, &req.Body)
	if err != nil {
		return nil, err
	}
	return api.NewResponse(&struct{}{}), nil
}
