package auth

import (
	"context"

	"github.com/SirNacou/weeate/backend/internal/common/api"
	"github.com/SirNacou/weeate/backend/internal/common/infrastructure/centrifugo"
	"github.com/SirNacou/weeate/backend/internal/common/infrastructure/configs"
)

type AuthEndpoint struct {
	env configs.Config
}

func NewAuthEndpoint(env configs.Config) *AuthEndpoint {
	return &AuthEndpoint{
		env: env,
	}
}

type GetCentrifugoJWTResponse struct {
	Token string `json:"token" doc:"The generated JWT token for Centrifugo authentication"`
}

func (e *AuthEndpoint) getCentrifugoToken(ctx context.Context, _ *struct{}) (*api.Response[GetCentrifugoJWTResponse], error) {
	user, ok := UserFromContext(ctx)
	if !ok {
		return nil, ErrUserNotFoundInContext
	}

	token, err := centrifugo.GenerateCentrifugoToken(user.ID,
		user.Email,
		[]byte(e.env.CENTRIFUGO_HMAC_SECRET),
	)

	return api.NewResponse(&GetCentrifugoJWTResponse{
		Token: token,
	}), err
}
