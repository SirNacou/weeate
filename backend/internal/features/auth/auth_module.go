package auth

import (
	"context"

	"github.com/SirNacou/weeate/backend/internal/common/infrastructure/configs"
	"github.com/danielgtaylor/huma/v2"
)

type AuthModule struct {
	endpoint *AuthEndpoint
}

func NewAuthModule(env configs.Env) *AuthModule {
	authEndpoint := NewAuthEndpoint(env)
	return &AuthModule{
		endpoint: authEndpoint,
	}
}

func (m *AuthModule) RegisterAPI(api huma.API) {
	group := huma.NewGroup(api, "/auth")
	huma.Get(group, "/centrifugo/token", m.endpoint.getCentrifugoToken)
}

func (m *AuthModule) Jobs() []func(context.Context) {
	return nil
}
