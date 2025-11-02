package db

import (
	"context"
	"fmt"
	"strconv"

	config "github.com/SirNacou/weeate/backend/internal/infrastructure/configs"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectToPostgres(ctx context.Context, cfg config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=%s",
		cfg.DBHost,
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBName,
		strconv.Itoa(cfg.DBPort),
		cfg.Timezone,
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	return db.WithContext(ctx), nil
}
