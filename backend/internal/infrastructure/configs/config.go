package config

import (
	"fmt"

	"github.com/caarlos0/env/v11"
)

type Config struct {
	PORT                 int    `env:"PORT" envDefault:"8080"`
	Timezone             string `env:"TZ" envDefault:"UTC"`
	DBHost               string `env:"DB_HOST" envDefault:"localhost"`
	DBPort               int    `env:"DB_PORT" envDefault:"5432"`
	DBUser               string `env:"DB_USER" envDefault:"weeate_user"`
	DBPassword           string `env:"DB_PASSWORD" envDefault:"weeate_password"`
	DBName               string `env:"DB_NAME" envDefault:"weeate_db"`
	SUPABASE_PROJECT_URL string `env:"SUPABASE_PROJECT_URL" required:"true"`
}

func LoadConfig() (Config, error) {
	c, err := env.ParseAs[Config]()
	if err != nil {
		return Config{}, fmt.Errorf("failed to parse config from environment variables: %w", err)
	}
	return c, nil
}
