package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"golang.org/x/sync/errgroup"

	"github.com/SirNacou/weeate/backend/internal/common/api"
	"github.com/SirNacou/weeate/backend/internal/common/infrastructure/bus"
	"github.com/SirNacou/weeate/backend/internal/common/infrastructure/centrifugo"
	"github.com/SirNacou/weeate/backend/internal/common/infrastructure/configs"
	"github.com/SirNacou/weeate/backend/internal/features/auth"
	"github.com/SirNacou/weeate/backend/internal/features/foods"
	"github.com/SirNacou/weeate/backend/internal/features/orders"
	"github.com/SirNacou/weeate/backend/internal/features/polls"
	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humafiber"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	slogfiber "github.com/samber/slog-fiber"
	"gorm.io/gorm"
)

type Server struct {
	config           configs.Env
	db               *gorm.DB
	supabaseService  *auth.SupabaseService
	bus              *bus.Bus
	centrifugoClient *centrifugo.CentrifugoClient
}

func NewServer(
	config configs.Env,
	db *gorm.DB,
	supabaseService *auth.SupabaseService,
	bus *bus.Bus,
	centrifugoClient *centrifugo.CentrifugoClient,
) *Server {
	return &Server{
		config:           config,
		db:               db,
		supabaseService:  supabaseService,
		bus:              bus,
		centrifugoClient: centrifugoClient,
	}
}

func (s *Server) buildHandler(ctx context.Context) (http.Handler, []func(context.Context), error) {
	// Setup Fiber app
	app := fiber.New(fiber.Config{})

	app.Use(slogfiber.NewWithConfig(slog.Default(), slogfiber.Config{
		WithRequestID:      true,
		WithTraceID:        true,
		WithSpanID:         true,
		WithRequestHeader:  true,
		WithResponseHeader: true,
		WithRequestBody:    true,
		WithResponseBody:   true,
	}))
	app.Use(recover.New())
	app.Use(api.CORSMiddleware(s.config.GO_ENV))
	authMiddleware, err := auth.AuthMiddleware(ctx, s.config.SUPABASE_AUTH_URL, s.config.SUPABASE_COOKIE_AUTH_NAME)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to initialize auth middleware: %v", err)
	}
	app.Use(authMiddleware)

	api := humafiber.New(app, huma.DefaultConfig("Weeate API", "v1.0.0"))

	authModule := auth.NewAuthModule(s.config)
	authModule.RegisterAPI(api)

	foodsModule := foods.NewFoodsModule(s.db, s.supabaseService)
	foodsModule.RegisterAPI(api)

	pollsModule := polls.NewPollsModule(s.bus, s.db, s.supabaseService, s.centrifugoClient)
	pollsModule.RegisterAPI(api)

	ordersModule := orders.NewOrdersModule(s.bus, s.db, s.supabaseService)
	ordersModule.RegisterAPI(api)

	huma.Get(api, "/", func(ctx context.Context, i *struct{}) (*auth.User, error) {
		user, ok := auth.UserFromContext(ctx)
		if !ok {
			return nil, huma.Error401Unauthorized("Unauthorized", fmt.Errorf("user not found in context"))
		}
		return user, nil
	})

	jobs := []func(context.Context){}
	jobs = append(jobs, authModule.Jobs()...)
	jobs = append(jobs, foodsModule.Jobs()...)
	jobs = append(jobs, pollsModule.Jobs()...)
	jobs = append(jobs, ordersModule.Jobs()...)

	return api.Adapter(), jobs, nil
}

func (s *Server) Run(ctx context.Context) error {
	server, jobs, err := s.buildHandler(ctx)
	if err != nil {
		return fmt.Errorf("failed to create server: %w", err)
	}

	g, gCtx := errgroup.WithContext(ctx)

	// Start Bus
	g.Go(func() error {
		return s.bus.Start(gCtx)
	})

	// Start Jobs
	for _, job := range jobs {
		g.Go(func() error {
			job(gCtx)
			return nil
		})
	}

	// Start HTTP server
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%v", s.config.PORT),
		Handler: server,
	}

	g.Go(func() error {
		slog.Info("Starting server", "port", s.config.PORT)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			return err
		}
		return nil
	})

	// Graceful shutdown
	g.Go(func() error {
		<-gCtx.Done()
		slog.Info("Shutting down server...")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		return srv.Shutdown(shutdownCtx)
	})

	return g.Wait()
}
