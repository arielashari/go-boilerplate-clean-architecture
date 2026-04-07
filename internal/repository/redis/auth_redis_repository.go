package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/Primuse-Pte-Ltd/go-boilerplate-clean-architecture/internal/entity"
	"github.com/redis/go-redis/v9"
)

type authRedisRepository struct {
	client *redis.Client
}

var _ entity.AuthRepository = (*authRedisRepository)(nil)

func NewAuthRedisRepository(client *redis.Client) entity.AuthRepository {
	return &authRedisRepository{client: client}
}

func (r *authRedisRepository) SetSession(ctx context.Context, userID, tokenID string, duration time.Duration) error {
	key := fmt.Sprintf("session:%s", userID)
	return r.client.Set(ctx, key, tokenID, duration).Err()
}

func (r *authRedisRepository) CheckSession(ctx context.Context, userID, tokenID string) (bool, error) {
	key := fmt.Sprintf("session:%s", userID)
	val, err := r.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return false, nil
	}
	return val == tokenID, err
}

func (r *authRedisRepository) DeleteSession(ctx context.Context, userID string) error {
	return r.client.Del(ctx, fmt.Sprintf("session:%s", userID)).Err()
}
