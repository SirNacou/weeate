package configs

import (
	"fmt"
	"strconv"

	"github.com/caarlos0/env/v11"
)

const (
	EnvDevelopment = "development"
	EnvProduction  = "production"
)

type Env struct {
	PORT                      int    `env:"PORT" envDefault:"8080"`
	Timezone                  string `env:"TZ" envDefault:"UTC"`
	DB                        db
	SUPABASE_URL              string `env:"SUPABASE_URL" required:"true"`
	SUPABASE_AUTH_URL         string `env:"SUPABASE_AUTH_URL" required:"true"`
	SUPABASE_API_KEY          string `env:"SUPABASE_API_KEY" required:"true"`
	SUPABASE_COOKIE_AUTH_NAME string `env:"SUPABASE_COOKIE_AUTH_NAME" required:"true"`
	GO_ENV                    string `env:"GO_ENV" envDefault:"development"`
	IMAGE_KIT_API_KEY         string `env:"IMAGE_KIT_API_KEY" required:"true"`
	IMAGEKIT_URL              string `env:"IMAGEKIT_URL" required:"true"`
}

type db struct {
	Host     string `env:"DB_HOST" envDefault:"localhost"`
	Port     int    `env:"DB_PORT" envDefault:"5432"`
	User     string `env:"DB_USER" envDefault:"weeate_user"`
	Password string `env:"DB_PASSWORD" envDefault:"weeate_password"`
	Name     string `env:"DB_NAME" envDefault:"weeate_db"`
}

func (e *Env) GetDBDsn() string {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=%s",
		e.DB.Host,
		e.DB.User,
		e.DB.Password,
		e.DB.Name,
		strconv.Itoa(e.DB.Port),
		e.Timezone,
	)
	return dsn
}

func LoadEnv() (Env, error) {
	e := Env{}
	err := env.Parse(&e)
	return e, err
}
