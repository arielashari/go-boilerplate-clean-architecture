package usecase

import (
	"context"

	"github.com/Primuse-Pte-Ltd/go-boilerplate-clean-architecture/internal/entity"
	"github.com/Primuse-Pte-Ltd/go-boilerplate-clean-architecture/internal/model"
	"github.com/Primuse-Pte-Ltd/go-boilerplate-clean-architecture/internal/model/mapper"
	"golang.org/x/crypto/bcrypt"
)

type UserUseCase interface {
	Create(ctx context.Context, request *model.CreateUserRequest) (*model.UserResponse, error)
	GetByID(ctx context.Context, id string) (*model.UserResponse, error)
	Update(ctx context.Context, request *model.UpdateUserRequest) (*model.UserResponse, error)
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, request *model.GetAllUsersRequest) (*model.PaginatedWebResponse[*model.UserResponse], error)
}

type userUseCase struct {
	userRepository entity.UserPostgresRepository
}

func NewUserUseCase(userRepository entity.UserPostgresRepository) UserUseCase {
	return &userUseCase{userRepository: userRepository}
}

func (c *userUseCase) Create(ctx context.Context, request *model.CreateUserRequest) (*model.UserResponse, error) {
	user := mapper.CreateUserRequestToEntity(request)

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	user.Password = string(hashedPassword)

	createdUser, err := c.userRepository.Create(ctx, user)
	if err != nil {
		return nil, err
	}

	return mapper.UserToResponse(createdUser), nil
}

func (c *userUseCase) GetByID(ctx context.Context, id string) (*model.UserResponse, error) {
	user, err := c.userRepository.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return mapper.UserToResponse(user), nil
}

func (c *userUseCase) List(ctx context.Context, request *model.GetAllUsersRequest) (*model.PaginatedWebResponse[*model.UserResponse], error) {
	users, err := c.userRepository.List(ctx, entity.UserQuery{
		Page:    request.Page,
		Limit:   request.Limit,
		Search:  request.Search,
		RoleID:  request.RoleID,
		SortBy:  request.SortBy,
		SortDir: request.SortDir,
	})
	if err != nil {
		return nil, err
	}

	totalCount, err := c.userRepository.Count(ctx)
	if err != nil {
		return nil, err
	}

	responses := make([]*model.UserResponse, len(users))
	for i, user := range users {
		responses[i] = mapper.UserToResponse(user)
	}

	totalPage := (totalCount + int64(request.Limit) - 1) / int64(request.Limit)
	return &model.PaginatedWebResponse[*model.UserResponse]{
		Items: responses,
		Metadata: model.PageMetadata{
			Page:      request.Page,
			TotalPage: totalPage,
			TotalItem: totalCount,
			Size:      request.Limit,
		},
	}, nil
}

func (c *userUseCase) Delete(ctx context.Context, id string) error {
	return c.userRepository.Delete(ctx, id)
}

func (c *userUseCase) Update(ctx context.Context, request *model.UpdateUserRequest) (*model.UserResponse, error) {
	existingUser, err := c.userRepository.GetByID(ctx, request.ID)
	if err != nil {
		return nil, err
	}

	user := mapper.UpdateUserRequestToEntity(existingUser, request)

	updatedUser, err := c.userRepository.Update(ctx, user)
	if err != nil {
		return nil, err
	}

	return mapper.UserToResponse(updatedUser), nil
}
