package auth

import "context"

type (
	userKey       struct{}
	userClaimsKey struct{}
)

func WithUser(ctx context.Context, u *User) context.Context {
	return context.WithValue(ctx, userKey{}, u)
}

func UserFromContext(ctx context.Context) (*User, bool) {
	val := ctx.Value(userKey{})
	user, ok := val.(*User)
	return user, ok
}

func WithUserClaims(ctx context.Context, claims map[string]interface{}) context.Context {
	return context.WithValue(ctx, userClaimsKey{}, claims)
}

func UserClaimsFromContext(ctx context.Context) (map[string]any, bool) {
	val := ctx.Value(userClaimsKey{})
	claims, ok := val.(map[string]any)
	return claims, ok
}
