package auth

import app_auth "github.com/SirNacou/weeate/backend/internal/usecase/auth"

type AuthEndpoint struct {
	registerCH   app_auth.RegisterCommandHandler
	loginCH      app_auth.LoginCommandHandler
	deleteUserCH app_auth.DeleteUserCommandHandler
}

func NewAuthEndpoint(registerCH app_auth.RegisterCommandHandler, loginCH app_auth.LoginCommandHandler, deleteUserCH app_auth.DeleteUserCommandHandler) *AuthEndpoint {
	return &AuthEndpoint{
		registerCH:   registerCH,
		loginCH:      loginCH,
		deleteUserCH: deleteUserCH,
	}
}
