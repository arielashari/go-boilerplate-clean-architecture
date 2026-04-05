package usecase

import (
	"context"

	"github.com/Primuse-Pte-Ltd/go-boilerplate-clean-architecture/internal/entity"
	"github.com/Primuse-Pte-Ltd/go-boilerplate-clean-architecture/internal/model"
	"github.com/Primuse-Pte-Ltd/go-boilerplate-clean-architecture/internal/model/mapper"
)

type RoleUseCase interface {
	Create(ctx context.Context, input *model.CreateRoleRequest) (*model.RoleResponse, error)
	GetByID(ctx context.Context, id string) (*model.RoleResponse, error)
	Update(ctx context.Context, input *model.UpdateRoleRequest) (*model.RoleResponse, error)
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, page, limit int) (*model.PaginatedWebResponse[*model.RoleResponse], error)
}

type roleUseCase struct {
	roleRepository entity.RolePostgresRepository
}

func NewRoleUseCase(roleRepository entity.RolePostgresRepository) RoleUseCase {
	return &roleUseCase{roleRepository: roleRepository}
}

func (c *roleUseCase) Create(ctx context.Context, request *model.CreateRoleRequest) (*model.RoleResponse, error) {
	role := mapper.CreateRoleRequestToEntity(request)

	createdRole, err := c.roleRepository.Create(ctx, role)
	if err != nil {
		return nil, err
	}

	return mapper.RoleToResponse(createdRole), nil
}

func (c *roleUseCase) GetByID(ctx context.Context, id string) (*model.RoleResponse, error) {
	role, err := c.roleRepository.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return mapper.RoleToResponse(role), nil
}

func (c *roleUseCase) List(ctx context.Context, page, limit int) (*model.PaginatedWebResponse[*model.RoleResponse], error) {
	offset := (page - 1) * limit
	roles, err := c.roleRepository.List(ctx, limit, offset)
	if err != nil {
		return nil, err
	}

	totalCount, err := c.roleRepository.Count(ctx)
	if err != nil {
		return nil, err
	}

	responseRoles := make([]*model.RoleResponse, len(roles))
	for i, role := range roles {
		responseRoles[i] = mapper.RoleToResponse(role)
	}

	totalPage := (totalCount + int64(limit) - 1) / int64(limit)
	return &model.PaginatedWebResponse[*model.RoleResponse]{
		Items: responseRoles,
		Metadata: model.PageMetadata{
			Page:      page,
			TotalPage: totalPage,
			TotalItem: totalCount,
			Size:      limit,
		},
	}, nil
}

func (c *roleUseCase) Update(ctx context.Context, request *model.UpdateRoleRequest) (*model.RoleResponse, error) {
	existingRole, err := c.roleRepository.GetByID(ctx, request.ID)
	if err != nil {
		return nil, err
	}

	updatedRoleEntity := mapper.UpdateRoleRequestToEntity(existingRole, request)

	updatedRole, err := c.roleRepository.Update(ctx, updatedRoleEntity)
	if err != nil {
		return nil, err
	}

	return mapper.RoleToResponse(updatedRole), nil
}

func (c *roleUseCase) Delete(ctx context.Context, id string) error {
	return c.roleRepository.Delete(ctx, id)
}
