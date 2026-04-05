package configs

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	App      AppConfig
	Database DatabaseConfig
	Redis    RedisConfig
	JWT      JWTConfig
}

type AppConfig struct {
	Name     string
	Port     int
	Env      string
	LogLevel string `mapstructure:"log_level"`
}

type DatabaseConfig struct {
	Postgres PostgresConfig
}

type PostgresConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Name     string
	SSLMode  string `mapstructure:"ssl_mode"`
}

type RedisConfig struct {
	Host     string
	Port     int
	Password string
	DB       int
}

type JWTConfig struct {
	Secret               string `mapstructure:"secret"`
	AccessExpireMinutes  int    `mapstructure:"access_expire_minutes"`
	RefreshExpireMinutes int    `mapstructure:"refresh_expire_minutes"`
}

func Load() (*Config, error) {
	env := viper.GetString("APP_ENV")
	if env == "" {
		env = "dev"
	}

	viper.SetConfigName("config." + env)
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./configs")
	viper.AddConfigPath(".")

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &cfg, nil
}
