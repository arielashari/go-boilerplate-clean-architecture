package entity

import (
	"context"
	"time"
)

type UserPostgresRepository interface {
	Create(ctx context.Context, user *User) (*User, error)
	GetByID(ctx context.Context, id string) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	GetUserForAuth(ctx context.Context, email string) (*User, error)
	Update(ctx context.Context, user *User) (*User, error)
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, query UserQuery) ([]*User, error)
	Count(ctx context.Context) (int64, error)
}

type RolePostgresRepository interface {
	Create(ctx context.Context, role *Role) (*Role, error)
	GetByID(ctx context.Context, id string) (*Role, error)
	Update(ctx context.Context, role *Role) (*Role, error)
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, limit, offset int) ([]*Role, error)
	Count(ctx context.Context) (int64, error)
}

type AuthRedisRepository interface {
	SetSession(ctx context.Context, userID, tokenID string, duration time.Duration) error
	CheckSession(ctx context.Context, userID, tokenID string) (bool, error)
	DeleteSession(ctx context.Context, userID string) error
}
