package auth

import (
	"context"

	"github.com/SirNacou/weeate/backend/internal/common/infrastructure/configs"
	"github.com/SirNacou/weeate/backend/internal/features/auth/api"
	"github.com/SirNacou/weeate/backend/internal/features/auth/services"
	"github.com/danielgtaylor/huma/v2"
	"github.com/imagekit-developer/imagekit-go/v2"
)

type AuthModule struct {
	endpoint *api.AuthEndpoint
}

func NewAuthModule(env configs.Config, ikClient *imagekit.Client) *AuthModule {
	handler := services.GetIKTokenQueryHandler(env.IMAGE_KIT_API_KEY, env.IMAGEKIT_PUBLIC_KEY)
	authEndpoint := api.NewAuthEndpoint(env, handler)
	return &AuthModule{
		endpoint: authEndpoint,
	}
}

func (m *AuthModule) RegisterAPI(api huma.API) {
	group := huma.NewGroup(api, "/auth")
	huma.Get(group, "/centrifugo/token", m.endpoint.GetCentrifugoToken)
	huma.Get(group, "/imagekit/token", m.endpoint.GetImageKitToken)
}

func (m *AuthModule) Jobs() []func(context.Context) {
	return nil
}
