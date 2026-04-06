package usecase_test

import (
	"context"
	"errors"
	"testing"

	"github.com/Primuse-Pte-Ltd/go-boilerplate-clean-architecture/internal/entity"
	"github.com/Primuse-Pte-Ltd/go-boilerplate-clean-architecture/internal/model"
	"github.com/Primuse-Pte-Ltd/go-boilerplate-clean-architecture/internal/usecase"
	"github.com/Primuse-Pte-Ltd/go-boilerplate-clean-architecture/internal/usecase/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupUserUseCase(t *testing.T) (usecase.UserUseCase, *mocks.MockUserRepository) {
	mockRepo := new(mocks.MockUserRepository)
	uc := usecase.NewUserUseCase(mockRepo)
	t.Cleanup(func() { mockRepo.AssertExpectations(t) })
	return uc, mockRepo
}

func TestUserUseCase_GetByID(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		uc, mockRepo := setupUserUseCase(t)

		expected := &entity.User{ID: "user-1", Email: "test@test.com", FirstName: "John", LastName: "Doe"}
		mockRepo.On("GetByID", mock.Anything, "user-1").Return(expected, nil)

		result, err := uc.GetByID(context.Background(), "user-1")

		assert.NoError(t, err)
		assert.Equal(t, expected.ID, result.ID)
		assert.Equal(t, expected.Email, result.Email)
	})

	t.Run("not found", func(t *testing.T) {
		uc, mockRepo := setupUserUseCase(t)

		mockRepo.On("GetByID", mock.Anything, "missing").Return(nil, entity.ErrNotFound)

		result, err := uc.GetByID(context.Background(), "missing")

		assert.Nil(t, result)
		assert.ErrorIs(t, err, entity.ErrNotFound)
	})
}

func TestUserUseCase_Create(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		uc, mockRepo := setupUserUseCase(t)

		req := &model.CreateUserRequest{
			Email:     "test@test.com",
			Password:  "password123",
			FirstName: "John",
			LastName:  "Doe",
		}
		mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*entity.User")).
			Return(&entity.User{ID: "user-1", Email: req.Email, FirstName: req.FirstName, LastName: req.LastName}, nil)

		result, err := uc.Create(context.Background(), req)

		assert.NoError(t, err)
		assert.Equal(t, req.Email, result.Email)
	})

	t.Run("repository error", func(t *testing.T) {
		uc, mockRepo := setupUserUseCase(t)

		req := &model.CreateUserRequest{
			Email:    "test@test.com",
			Password: "password123",
		}
		mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*entity.User")).
			Return(nil, entity.ErrConflict)

		result, err := uc.Create(context.Background(), req)

		assert.Nil(t, result)
		assert.ErrorIs(t, err, entity.ErrConflict)
	})
}

func TestUserUseCase_Update(t *testing.T) {
	firstName := "New"
	lastName := "Name"
	t.Run("success", func(t *testing.T) {
		uc, mockRepo := setupUserUseCase(t)

		existing := &entity.User{ID: "user-1", Email: "old@test.com", FirstName: "Old", LastName: "Name"}
		updated := &entity.User{ID: "user-1", Email: "old@test.com", FirstName: "New", LastName: "Name"}

		mockRepo.On("GetByID", mock.Anything, "user-1").Return(existing, nil)
		mockRepo.On("Update", mock.Anything, mock.AnythingOfType("*entity.User")).Return(updated, nil)

		result, err := uc.Update(context.Background(), &model.UpdateUserRequest{ID: "user-1", FirstName: &firstName, LastName: &lastName})

		assert.NoError(t, err)
		assert.Equal(t, "New", result.FirstName)
	})

	t.Run("user not found", func(t *testing.T) {
		uc, mockRepo := setupUserUseCase(t)

		mockRepo.On("GetByID", mock.Anything, "missing").Return(nil, entity.ErrNotFound)

		result, err := uc.Update(context.Background(), &model.UpdateUserRequest{ID: "missing"})

		assert.Nil(t, result)
		assert.ErrorIs(t, err, entity.ErrNotFound)
	})
}

func TestUserUseCase_Delete(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		uc, mockRepo := setupUserUseCase(t)

		mockRepo.On("Delete", mock.Anything, "user-1").Return(nil)

		err := uc.Delete(context.Background(), "user-1")
		assert.NoError(t, err)
	})

	t.Run("not found", func(t *testing.T) {
		uc, mockRepo := setupUserUseCase(t)

		mockRepo.On("Delete", mock.Anything, "missing").Return(entity.ErrNotFound)

		err := uc.Delete(context.Background(), "missing")
		assert.ErrorIs(t, err, entity.ErrNotFound)
	})
}

func TestUserUseCase_List(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		uc, mockRepo := setupUserUseCase(t)

		users := []*entity.User{
			{ID: "user-1", Email: "a@test.com"},
			{ID: "user-2", Email: "b@test.com"},
		}
		mockRepo.On("List", mock.Anything, mock.AnythingOfType("entity.UserQuery")).Return(users, nil)
		mockRepo.On("Count", mock.Anything).Return(int64(2), nil)

		result, err := uc.List(context.Background(), &model.GetAllUsersRequest{
			PaginationRequest: model.PaginationRequest{Page: 1, Limit: 10},
		})

		assert.NoError(t, err)
		assert.Len(t, result.Items, 2)
		assert.Equal(t, int64(2), result.Metadata.TotalItem)
		assert.Equal(t, int64(1), result.Metadata.TotalPage)
	})

	t.Run("repository error", func(t *testing.T) {
		uc, mockRepo := setupUserUseCase(t)

		mockRepo.On("List", mock.Anything, mock.AnythingOfType("entity.UserQuery")).
			Return(nil, errors.New("db error"))

		result, err := uc.List(context.Background(), &model.GetAllUsersRequest{
			PaginationRequest: model.PaginationRequest{Page: 1, Limit: 10},
		})

		assert.Nil(t, result)
		assert.Error(t, err)
	})
}
