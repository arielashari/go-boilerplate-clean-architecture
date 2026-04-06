package usecase_test

import (
	"context"
	"testing"

	"github.com/Primuse-Pte-Ltd/go-boilerplate-clean-architecture/internal/entity"
	"github.com/Primuse-Pte-Ltd/go-boilerplate-clean-architecture/internal/model"
	"github.com/Primuse-Pte-Ltd/go-boilerplate-clean-architecture/internal/usecase"
	"github.com/Primuse-Pte-Ltd/go-boilerplate-clean-architecture/internal/usecase/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupRoleUseCase(t *testing.T) (usecase.RoleUseCase, *mocks.MockRoleRepository) {
	mockRepo := new(mocks.MockRoleRepository)
	uc := usecase.NewRoleUseCase(mockRepo)
	t.Cleanup(func() { mockRepo.AssertExpectations(t) })
	return uc, mockRepo
}

func TestRoleUseCase_Create(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		uc, mockRepo := setupRoleUseCase(t)

		mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*entity.Role")).
			Return(&entity.Role{ID: "role-1", Name: "admin"}, nil)

		result, err := uc.Create(context.Background(), &model.CreateRoleRequest{Name: "admin"})

		assert.NoError(t, err)
		assert.Equal(t, "admin", result.Name)
	})

	t.Run("conflict", func(t *testing.T) {
		uc, mockRepo := setupRoleUseCase(t)

		mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*entity.Role")).
			Return(nil, entity.ErrConflict)

		result, err := uc.Create(context.Background(), &model.CreateRoleRequest{Name: "admin"})

		assert.Nil(t, result)
		assert.ErrorIs(t, err, entity.ErrConflict)
	})
}

func TestRoleUseCase_GetByID(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		uc, mockRepo := setupRoleUseCase(t)

		mockRepo.On("GetByID", mock.Anything, "role-1").
			Return(&entity.Role{ID: "role-1", Name: "admin"}, nil)

		result, err := uc.GetByID(context.Background(), "role-1")

		assert.NoError(t, err)
		assert.Equal(t, "role-1", result.ID)
	})

	t.Run("not found", func(t *testing.T) {
		uc, mockRepo := setupRoleUseCase(t)

		mockRepo.On("GetByID", mock.Anything, "missing").Return(nil, entity.ErrNotFound)

		result, err := uc.GetByID(context.Background(), "missing")

		assert.Nil(t, result)
		assert.ErrorIs(t, err, entity.ErrNotFound)
	})
}

func TestRoleUseCase_Update(t *testing.T) {
	name := "New Name"
	t.Run("success", func(t *testing.T) {
		uc, mockRepo := setupRoleUseCase(t)

		mockRepo.On("GetByID", mock.Anything, "role-1").
			Return(&entity.Role{ID: "role-1", Name: "old"}, nil)
		mockRepo.On("Update", mock.Anything, mock.AnythingOfType("*entity.Role")).
			Return(&entity.Role{ID: "role-1", Name: "new"}, nil)

		result, err := uc.Update(context.Background(), &model.UpdateRoleRequest{ID: "role-1", Name: &name})

		assert.NoError(t, err)
		assert.Equal(t, "new", result.Name)
	})

	t.Run("not found", func(t *testing.T) {
		uc, mockRepo := setupRoleUseCase(t)

		mockRepo.On("GetByID", mock.Anything, "missing").Return(nil, entity.ErrNotFound)

		result, err := uc.Update(context.Background(), &model.UpdateRoleRequest{ID: "missing"})

		assert.Nil(t, result)
		assert.ErrorIs(t, err, entity.ErrNotFound)
	})
}

func TestRoleUseCase_Delete(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		uc, mockRepo := setupRoleUseCase(t)
		mockRepo.On("Delete", mock.Anything, "role-1").Return(nil)

		err := uc.Delete(context.Background(), "role-1")
		assert.NoError(t, err)
	})

	t.Run("not found", func(t *testing.T) {
		uc, mockRepo := setupRoleUseCase(t)
		mockRepo.On("Delete", mock.Anything, "missing").Return(entity.ErrNotFound)

		err := uc.Delete(context.Background(), "missing")
		assert.ErrorIs(t, err, entity.ErrNotFound)
	})
}

func TestRoleUseCase_List(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		uc, mockRepo := setupRoleUseCase(t)

		roles := []*entity.Role{
			{ID: "role-1", Name: "admin"},
			{ID: "role-2", Name: "user"},
		}
		mockRepo.On("List", mock.Anything, 10, 0).Return(roles, nil)
		mockRepo.On("Count", mock.Anything).Return(int64(2), nil)

		result, err := uc.List(context.Background(), 1, 10)

		assert.NoError(t, err)
		assert.Len(t, result.Items, 2)
		assert.Equal(t, int64(1), result.Metadata.TotalPage)
	})
}
