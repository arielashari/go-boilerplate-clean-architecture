package database

import (
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type DatabaseInfrastructure struct {
	Postgres *pgxpool.Pool
	Redis    *redis.Client
}

func (i *DatabaseInfrastructure) CloseAll() {
	if i.Postgres != nil {
		i.Postgres.Close()
	}

	if i.Redis != nil {
		if err := i.Redis.Close(); err != nil {
			log.Printf("Error closing Redis client: %v", err)
		}
	}
}
