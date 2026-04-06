package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type MockTransactor struct {
	mock.Mock
}

func (m *MockTransactor) WithTx(ctx context.Context, fn func(ctx context.Context) error) error {
	return fn(ctx)
}
