package testhelpers

import (
	"context"
	"fmt"
	"log"
	"time"

	food "github.com/SirNacou/weeate/backend/internal/features/foods/domain"
	poll "github.com/SirNacou/weeate/backend/internal/features/polls/domain"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	postgresDriver "gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// setupGlobalDB spins up Postgres and returns the DB + a cleanup function.
// Note: We use 'log' here because we don't have access to 't' yet.
func SetupGORM(ctx context.Context) (db *gorm.DB, cleanup func(), err error) {
	// 1. Start Container
	pgContainer, err := postgres.Run(ctx,
		"postgres:18-alpine",
		testcontainers.WithName("postgres_test"),
		postgres.WithDatabase("testdb"),
		postgres.WithUsername("user"),
		postgres.WithPassword("password"),
		testcontainers.WithWaitStrategy(
			wait.ForListeningPort("5432/tcp").
				WithStartupTimeout(60*time.Second),
		),
	)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to start container: %w", err)
	}

	// 2. Connect GORM
	connStr, _ := pgContainer.ConnectionString(ctx, "sslmode=disable")
	db, err = gorm.Open(postgresDriver.Open(connStr), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return nil, nil, fmt.Errorf("failed to connect to db: %w", err)
	}

	err = db.AutoMigrate(&poll.Poll{}, &poll.PollOption{}, &poll.Vote{},
		&food.Food{})
	if err != nil {
		return nil, nil, fmt.Errorf("failed to migrate db: %w", err)
	}

	// 3. Create a Cleanup Closure
	// This lets TestMain terminate the container when it's done
	teardown := func() {
		if err := pgContainer.Terminate(ctx); err != nil {
			log.Printf("failed to terminate container: %s", err)
		}
	}

	return db, teardown, nil
}
