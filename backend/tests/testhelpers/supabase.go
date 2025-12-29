package testhelpers

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/SirNacou/weeate/backend/internal/common/infrastructure/supabase"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

// SetupSupabaseTestDB creates a Postgres container with the user_profiles table
// and returns a SupabaseService configured to use it directly via raw SQL
func SetupSupabaseTestDB(ctx context.Context) (db *sql.DB, cleanup func(), err error) {
	// Start Postgres container
	pgContainer, err := postgres.Run(ctx,
		"postgres:18-alpine",
		testcontainers.WithName("supabase_test_postgres"),
		postgres.WithDatabase("testdb"),
		postgres.WithUsername("user"),
		postgres.WithPassword("password"),
		testcontainers.WithWaitStrategy(
			wait.ForListeningPort("5432/tcp").
				WithStartupTimeout(60*time.Second),
		),
	)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to start postgres container: %w", err)
	}

	// Get connection string
	connStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get connection string: %w", err)
	}

	// Connect to database
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to connect to db: %w", err)
	}

	// Create user_profiles table
	_, err = db.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS user_profiles (
			id UUID PRIMARY KEY,
			avatar_url TEXT,
			display_name TEXT,
			created_at TIMESTAMPTZ DEFAULT NOW()
		)
	`)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create user_profiles table: %w", err)
	}

	cleanup = func() {
		db.Close()
		pgContainer.Terminate(ctx)
	}

	return db, cleanup, nil
}

// MockSupabaseService creates a mock Supabase service for testing
// Note: This is a placeholder - for true integration tests with Supabase API,
// you would need to either:
// 1. Use Supabase local development (supabase/edge-runtime)
// 2. Point to a test Supabase project
// 3. Mock the Supabase client entirely
func SetupMockSupabaseService(ctx context.Context, db *sql.DB) (*supabase.SupabaseService, error) {
	// For now, we return nil since SupabaseService requires the actual Supabase SDK
	// In a real scenario, you'd either:
	// - Create a TestSupabaseService that implements the same interface but uses raw SQL
	// - Or use the actual Supabase local setup
	return nil, fmt.Errorf("mock supabase service not implemented - use direct DB access for tests")
}
