package usecase_test

import (
	"context"
	"testing"

	"github.com/Primuse-Pte-Ltd/go-boilerplate-clean-architecture/configs"
	"github.com/Primuse-Pte-Ltd/go-boilerplate-clean-architecture/internal/entity"
	"github.com/Primuse-Pte-Ltd/go-boilerplate-clean-architecture/internal/model"
	"github.com/Primuse-Pte-Ltd/go-boilerplate-clean-architecture/internal/usecase"
	"github.com/Primuse-Pte-Ltd/go-boilerplate-clean-architecture/internal/usecase/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupAuthUseCase(t *testing.T) (usecase.AuthUseCase, *mocks.MockAuthRepository, *mocks.MockUserRepository) {
	mockAuthRepo := new(mocks.MockAuthRepository)
	mockUserRepo := new(mocks.MockUserRepository)
	mockTransactor := new(mocks.MockTransactor)
	mockMailer := new(mocks.MockEmailSender)
	cfg := &configs.Config{
		JWT: configs.JWTConfig{
			Secret:               "test-secret",
			AccessExpireMinutes:  15,
			RefreshExpireMinutes: 10080,
		},
	}
	uc := usecase.NewAuthUseCase(mockAuthRepo, mockUserRepo, mockTransactor, mockMailer, cfg)
	t.Cleanup(func() {
		mockAuthRepo.AssertExpectations(t)
		mockUserRepo.AssertExpectations(t)
	})
	return uc, mockAuthRepo, mockUserRepo
}

func TestAuthUseCase_Register(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		uc, _, mockUserRepo := setupAuthUseCase(t)

		mockUserRepo.On("GetByEmail", mock.Anything, "test@test.com").Return(nil, entity.ErrNotFound)
		mockUserRepo.On("Create", mock.Anything, mock.AnythingOfType("*entity.User")).
			Return(&entity.User{ID: "user-1", Email: "test@test.com", FirstName: "John", LastName: "Doe"}, nil)

		result, err := uc.Register(context.Background(), &model.RegisterRequest{
			Email:     "test@test.com",
			Password:  "password123",
			FirstName: "John",
			LastName:  "Doe",
		})

		assert.NoError(t, err)
		assert.Equal(t, "test@test.com", result.Email)
	})

	t.Run("email already exists", func(t *testing.T) {
		uc, _, mockUserRepo := setupAuthUseCase(t)

		mockUserRepo.On("GetByEmail", mock.Anything, "taken@test.com").
			Return(&entity.User{ID: "existing", Email: "taken@test.com"}, nil)

		result, err := uc.Register(context.Background(), &model.RegisterRequest{
			Email:    "taken@test.com",
			Password: "password123",
		})

		assert.Nil(t, result)
		assert.ErrorIs(t, err, entity.ErrEmailAlreadyExists)
	})
}

func TestAuthUseCase_Login(t *testing.T) {
	t.Run("user not found", func(t *testing.T) {
		uc, _, mockUserRepo := setupAuthUseCase(t)

		mockUserRepo.On("GetUserForAuth", mock.Anything, "ghost@test.com").
			Return(nil, entity.ErrNotFound)

		result, err := uc.Login(context.Background(), &model.LoginRequest{
			Email:    "ghost@test.com",
			Password: "password123",
		})

		assert.Nil(t, result)
		assert.ErrorIs(t, err, entity.ErrInvalidCredentials)
	})
}

func TestAuthUseCase_Logout(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		uc, mockAuthRepo, _ := setupAuthUseCase(t)

		mockAuthRepo.On("DeleteSession", mock.Anything, "user-1").Return(nil)

		err := uc.Logout(context.Background(), "user-1")
		assert.NoError(t, err)
	})
}
