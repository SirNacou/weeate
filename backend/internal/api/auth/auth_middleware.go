package auth

import (
	"context"
	"crypto/ecdsa"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"strings"
	"time"

	config "github.com/SirNacou/weeate/backend/internal/infrastructure/configs"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/lestrrat-go/httprc/v3"
	"github.com/lestrrat-go/httprc/v3/errsink"
	"github.com/lestrrat-go/httprc/v3/tracesink"
	"github.com/lestrrat-go/jwx/v3/jwk"
)

type AuthMiddleware struct {
	cache      *jwk.Cache
	cookieName string
}

func NewAuthMiddleware(ctx context.Context, config config.Config) (*AuthMiddleware, error) {
	jwkCache, err := jwk.NewCache(ctx, httprc.NewClient(
		httprc.WithErrorSink(errsink.NewSlog(slog.Default())),
		httprc.WithHTTPClient(http.DefaultClient),
		httprc.WithTraceSink(tracesink.NewSlog(slog.Default())),
	))
	if err != nil {
		return nil, err
	}

	if err = jwkCache.Register(ctx, config.SUPABASE_AUTH_URL, jwk.WithMaxInterval(10*time.Minute)); err != nil {
		return nil, err
	}
	return &AuthMiddleware{
		cache:      jwkCache,
		cookieName: config.SUPABASE_COOKIE_AUTH_NAME,
	}, nil
}

func (m *AuthMiddleware) Handle(c *fiber.Ctx) error {
	if c.IsFromLocal() && (c.Path() == "/docs" || c.Path() == "openapi.json") {
		return c.Next()
	}
	authCookie := c.Cookies(m.cookieName)

	if !strings.HasPrefix(authCookie, "base64-") {
		return c.SendStatus(http.StatusUnauthorized)
	}
	b64String := strings.TrimPrefix(authCookie, "base64-")

	jsonBytes, err := base64.URLEncoding.DecodeString(b64String)
	if err != nil {
		return c.Status(http.StatusUnauthorized).SendString(err.Error())
	}

	var session SupabaseSession
	if err = json.Unmarshal(jsonBytes, &session); err != nil {
		return c.Status(http.StatusUnauthorized).SendString(err.Error())
	}

	token, err := jwt.Parse(session.AccessToken, func(t *jwt.Token) (any, error) {
		ctx := context.Background()
		iss, err := t.Claims.GetIssuer()
		if err != nil {
			return nil, err
		}

		pubKeyUrl, err := url.JoinPath(iss, ".well-known/jwks.json")
		if err != nil {
			return nil, err
		}

		set, err := m.cache.Lookup(ctx, pubKeyUrl)
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

	c.Context().SetUserValue("user", session.User)

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		c.Context().SetUserValue("user_claims", claims)
	} else {
		slog.Warn("user claims not found")
	}

	return c.Next()
}
