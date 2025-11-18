package data

import (
	"context"
	"fmt"

	foods_domain "github.com/SirNacou/weeate/backend/internal/features/foods/domain"
	polls_domain "github.com/SirNacou/weeate/backend/internal/features/polls/domain"
	orders_domain "github.com/SirNacou/weeate/backend/internal/features/orders/domain"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectToPostgres(ctx context.Context, dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	return db.WithContext(ctx), nil
}

func MigratePostgresDB(db *gorm.DB) error {
	if err := db.AutoMigrate(&foods_domain.Food{},
		&polls_domain.Poll{},
		&polls_domain.PollOption{},
		&polls_domain.Vote{},
		&orders_domain.Order{}); err != nil {
		return fmt.Errorf("failed to migrate database: %w", err)
	}
	return nil
}
