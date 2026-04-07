package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/Primuse-Pte-Ltd/go-boilerplate-clean-architecture/internal/entity"
)

type MockFileStorage struct {
	mock.Mock
}

func (m *MockFileStorage) Upload(ctx context.Context, input *entity.UploadInput) (*entity.UploadResult, error) {
	args := m.Called(ctx, input)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.UploadResult), args.Error(1)
}

func (m *MockFileStorage) Delete(ctx context.Context, key string) error {
	return m.Called(ctx, key).Error(0)
}

func (m *MockFileStorage) GetPresignedURL(ctx context.Context, key string, operation string) (string, error) {
	args := m.Called(ctx, key, operation)
	return args.String(0), args.Error(1)
}

func (m *MockFileStorage) GetPublicURL(key string) string {
	return m.Called(key).String(0)
}
