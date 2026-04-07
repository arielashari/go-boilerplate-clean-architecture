package configs

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	App       AppConfig       `mapstructure:"app"`
	Database  DatabaseConfig  `mapstructure:"database"`
	Redis     RedisConfig     `mapstructure:"redis"`
	JWT       JWTConfig       `mapstructure:"jwt"`
	RateLimit RateLimitConfig `mapstructure:"ratelimit"`
	CORS      CORSConfig      `mapstructure:"cors"`
	SMTP      SMTPConfig      `mapstructure:"smtp"`
	S3        S3Config        `mapstructure:"s3"`
}

type AppConfig struct {
	Name     string `mapstructure:"name"`
	Port     int    `mapstructure:"port"`
	Env      string `mapstructure:"env"`
	LogLevel string `mapstructure:"log_level"`
}

type DatabaseConfig struct {
	Postgres PostgresConfig `mapstructure:"postgres"`
}

type PostgresConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Name     string `mapstructure:"name"`
	SSLMode  string `mapstructure:"ssl_mode"`
}

type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

type JWTConfig struct {
	Secret               string `mapstructure:"secret"`
	AccessExpireMinutes  int    `mapstructure:"access_expire_minutes"`
	RefreshExpireMinutes int    `mapstructure:"refresh_expire_minutes"`
}

type RateLimitConfig struct {
	MaxRequests       int `mapstructure:"max_requests"`
	ExpirationSeconds int `mapstructure:"expiration_seconds"`
}

type CORSConfig struct {
	AllowOrigins []string `mapstructure:"allow_origins"`
	AllowHeaders []string `mapstructure:"allow_headers"`
}

type SMTPConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	From     string `mapstructure:"from"`
	FromName string `mapstructure:"from_name"`
}

type S3Config struct {
	AccessKeyID          string `mapstructure:"access_key_id"`
	SecretAccessKey      string `mapstructure:"secret_access_key"`
	Region               string `mapstructure:"region"`
	Bucket               string `mapstructure:"bucket"`
	BaseURL              string `mapstructure:"base_url"`
	PresignExpiryMinutes int    `mapstructure:"presign_expiry_minutes"`
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
