package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type MockEmailSender struct {
	mock.Mock
}

func (m *MockEmailSender) SendVerificationEmail(ctx context.Context, to, name, otp string) error {
	return m.Called(ctx, to, name, otp).Error(0)
}

func (m *MockEmailSender) SendPasswordResetEmail(ctx context.Context, to, name, resetLink string) error {
	return m.Called(ctx, to, name, resetLink).Error(0)
}

func (m *MockEmailSender) SendWelcomeEmail(ctx context.Context, to, name string) error {
	return m.Called(ctx, to, name).Error(0)
}

func (m *MockEmailSender) SendNotificationEmail(ctx context.Context, to, name, subject, message string) error {
	return m.Called(ctx, to, name, subject, message).Error(0)
}
