package db

import (
	"context"
	"fmt"

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
