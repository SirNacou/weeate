package auth

import (
	"context"
	"errors"
)

// SupabaseSession represents the complete session object
// stored in the @supabase/ssr cookie.
type SupabaseSession struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int64  `json:"expires_in"`
	ExpiresAt    int64  `json:"expires_at"`
	RefreshToken string `json:"refresh_token"`
	User         User   `json:"user"`
}

// User represents the Supabase user object.
type User struct {
	ID               string       `json:"id"`
	Audience         string       `json:"aud"`
	Role             string       `json:"role"`
	Email            string       `json:"email"`
	EmailConfirmedAt string       `json:"email_confirmed_at"`
	Phone            string       `json:"phone"`
	ConfirmedAt      string       `json:"confirmed_at"`
	LastSignInAt     string       `json:"last_sign_in_at"`
	AppMetadata      AppMetadata  `json:"app_metadata"`
	UserMetadata     UserMetadata `json:"user_metadata"`
	Identities       []Identity   `json:"identities"`
	CreatedAt        string       `json:"created_at"`
	UpdatedAt        string       `json:"updated_at"`
	IsAnonymous      bool         `json:"is_anonymous"`
}

// AppMetadata contains application-specific metadata for the user.
type AppMetadata struct {
	AvatarURL   string   `json:"avatar_url"`
	DisplayName string   `json:"display_name"`
	Provider    string   `json:"provider"`
	Providers   []string `json:"providers"`
}

// UserMetadata contains custom user-specific metadata.
type UserMetadata struct {
	EmailVerified bool `json:"email_verified"`
}

// Identity represents a user's identity from a specific provider.
type Identity struct {
	IdentityID   string       `json:"identity_id"`
	ID           string       `json:"id"`
	UserID       string       `json:"user_id"`
	IdentityData IdentityData `json:"identity_data"`
	Provider     string       `json:"provider"`
	LastSignInAt string       `json:"last_sign_in_at"`
	CreatedAt    string       `json:"created_at"`
	UpdatedAt    string       `json:"updated_at"`
	Email        string       `json:"email"`
}

// IdentityData contains the raw data from an identity provider.
type IdentityData struct {
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	PhoneVerified bool   `json:"phone_verified"`
	Subject       string `json:"sub"`
}

func GetUserContext(ctx context.Context) (User, error) {
	user := ctx.Value("user")

	u, ok := user.(User)
	if !ok {
		return User{}, errors.New("invalid user in context")
	}

	return u, nil
}
