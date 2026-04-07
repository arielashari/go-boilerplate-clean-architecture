package entity

import (
	"context"
	"time"
)

type UserRepository interface {
	Create(ctx context.Context, user *User) (*User, error)
	GetByID(ctx context.Context, id string) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	GetUserForAuth(ctx context.Context, email string) (*User, error)
	Update(ctx context.Context, user *User) (*User, error)
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, query UserQuery) ([]*User, error)
	Count(ctx context.Context) (int64, error)
}

type RoleRepository interface {
	Create(ctx context.Context, role *Role) (*Role, error)
	GetByID(ctx context.Context, id string) (*Role, error)
	GetByName(ctx context.Context, name string) (*Role, error)
	Update(ctx context.Context, role *Role) (*Role, error)
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, limit, offset int) ([]*Role, error)
	Count(ctx context.Context) (int64, error)
}

type AuthRepository interface {
	SetSession(ctx context.Context, userID, tokenID string, duration time.Duration) error
	CheckSession(ctx context.Context, userID, tokenID string) (bool, error)
	DeleteSession(ctx context.Context, userID string) error
}

type Transactor interface {
	WithTx(ctx context.Context, fn func(ctx context.Context) error) error
}

type EmailSender interface {
	SendVerificationEmail(ctx context.Context, to, name, otp string) error
	SendPasswordResetEmail(ctx context.Context, to, name, resetLink string) error
	SendWelcomeEmail(ctx context.Context, to, name string) error
	SendNotificationEmail(ctx context.Context, to, name, subject, message string) error
}

type FileStorage interface {
	Upload(ctx context.Context, input *UploadInput) (*UploadResult, error)
	Delete(ctx context.Context, key string) error
	GetPresignedURL(ctx context.Context, key string, operation string) (string, error)
	GetPublicURL(key string) string
}
