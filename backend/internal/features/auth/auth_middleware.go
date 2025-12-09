package auth

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/SirNacou/weeate/backend/internal/features/auth/domain"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/lestrrat-go/httprc/v3"
	"github.com/lestrrat-go/httprc/v3/errsink"
	"github.com/lestrrat-go/httprc/v3/tracesink"
	"github.com/lestrrat-go/jwx/v3/jwk"
)

func AuthMiddleware(ctx context.Context, authUrl, cookieName string) (fiber.Handler, error) {
	jwkCache, err := jwk.NewCache(ctx, httprc.NewClient(
		httprc.WithErrorSink(errsink.NewSlog(slog.Default())),
		httprc.WithHTTPClient(http.DefaultClient),
		httprc.WithTraceSink(tracesink.NewSlog(slog.Default())),
	))
	if err != nil {
		return nil, err
	}

	if err = jwkCache.Register(ctx, authUrl, jwk.WithMaxInterval(10*time.Minute)); err != nil {
		return nil, err
	}

	return func(c *fiber.Ctx) error {
		if strings.Contains(c.Path(), "/docs") || strings.Contains(c.Path(), "/openapi") {
			return c.Next()
		}

		authorizationHeader := c.GetReqHeaders()["Authorization"]
		if len(authorizationHeader) < 1 ||
			!strings.HasPrefix(authorizationHeader[0], "Bearer ") {
			slog.ErrorContext(ctx, "invalid authorization header")
			return c.SendStatus(http.StatusUnauthorized)
		}
		accessToken := strings.SplitN(authorizationHeader[0], " ", 2)[1]

		token, err := jwt.Parse(accessToken, func(t *jwt.Token) (any, error) {
			ctx := context.Background()
			iss, err := t.Claims.GetIssuer()
			if err != nil {
				return nil, err
			}

			pubKeyUrl, err := url.JoinPath(iss, ".well-known/jwks.json")
			if err != nil {
				return nil, err
			}

			set, err := jwkCache.Lookup(ctx, pubKeyUrl)
			if err != nil {
				return nil, err
			}

			keyID, ok := t.Header["kid"].(string)
			if !ok {
				return nil, errors.New("expecting JWT header to have a key ID in the kid field")
			}

			key, found := set.LookupKeyID(keyID)
			if !found {
				return nil, fmt.Errorf("unable to find key %q", keyID)
			}

			publicKey, err := key.PublicKey()
			if err != nil {
				return nil, fmt.Errorf("unable to extract public key: %w", err)
			}

			ecdsaPubKey := &ecdsa.PublicKey{}
			jwk.Export(publicKey, ecdsaPubKey)

			return ecdsaPubKey, nil
		})
		if err != nil {
			return c.Status(http.StatusUnauthorized).SendString(err.Error())
		}

		if !token.Valid {
			return c.SendStatus(http.StatusUnauthorized)
		}

		ctx := c.UserContext()

		// Extract user info from JWT claims
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			user := &User{
				ID:    claims["sub"].(string),
				Email: claims["email"].(string),
			}
			if role, ok := claims["role"].(string); ok {
				user.Role = role
			}
			ctx = WithUser(ctx, user)
			ctx = WithUserClaims(ctx, claims)
		} else {
			slog.Warn("user claims not found")
		}

		c.SetUserContext(ctx)

		return c.Next()
	}, nil
}
