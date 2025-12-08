package domain

import "context"

type (
	UserKey       struct{}
	UserClaimsKey struct{}
)

func WithUser(ctx context.Context, u *User) context.Context {
	return context.WithValue(ctx, UserKey{}, u)
}

func UserFromContext(ctx context.Context) (*User, bool) {
	val := ctx.Value(UserKey{})
	user, ok := val.(*User)
	return user, ok
}

func WithUserClaims(ctx context.Context, claims map[string]interface{}) context.Context {
	return context.WithValue(ctx, UserClaimsKey{}, claims)
}

func UserClaimsFromContext(ctx context.Context) (map[string]any, bool) {
	val := ctx.Value(UserClaimsKey{})
	claims, ok := val.(map[string]any)
	return claims, ok
}
