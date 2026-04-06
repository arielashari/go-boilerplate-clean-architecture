package mocks

import (
	"context"
	"time"

	"github.com/stretchr/testify/mock"
)

type MockAuthRepository struct {
	mock.Mock
}

func (m *MockAuthRepository) SetSession(ctx context.Context, userID, tokenID string, duration time.Duration) error {
	args := m.Called(ctx, userID, tokenID, duration)
	return args.Error(0)
}

func (m *MockAuthRepository) CheckSession(ctx context.Context, userID, tokenID string) (bool, error) {
	args := m.Called(ctx, userID, tokenID)
	return args.Bool(0), args.Error(1)
}

func (m *MockAuthRepository) DeleteSession(ctx context.Context, userID string) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}
