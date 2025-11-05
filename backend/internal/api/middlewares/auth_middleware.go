package middlewares

import (
	"context"
	"log/slog"
	"time"

	config "github.com/SirNacou/weeate/backend/internal/infrastructure/configs"
	"github.com/labstack/echo/v4"
	"github.com/lestrrat-go/httprc/v3"
	"github.com/lestrrat-go/httprc/v3/errsink"
	"github.com/lestrrat-go/httprc/v3/tracesink"
	"github.com/lestrrat-go/jwx/v3/jwk"
)

type AuthMiddleware struct {
	c *jwk.Cache
}

func NewAuthMiddleware(ctx *context.Context, cfg *config.Config, e *echo.Echo) (*AuthMiddleware, error) {
	jwkCache, err := jwk.NewCache(*ctx, httprc.NewClient(
		httprc.WithErrorSink(errsink.NewSlog(slog.Default())),
		httprc.WithHTTPClient(e.AutoTLSManager.Client.HTTPClient),
		httprc.WithTraceSink(tracesink.NewSlog(slog.Default())),
	))
	if err != nil {
		return nil, err
	}

	if err = jwkCache.Register(*ctx, cfg.SUPABASE_AUTH_URL, jwk.WithMaxInterval(10*time.Minute)); err != nil {
		return nil, err
	}
	return &AuthMiddleware{
		c: jwkCache,
	}, nil
}