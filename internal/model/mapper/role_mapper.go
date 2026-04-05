package mapper

import (
	"github.com/Primuse-Pte-Ltd/go-boilerplate-clean-architecture/internal/entity"
	"github.com/Primuse-Pte-Ltd/go-boilerplate-clean-architecture/internal/model"
	"github.com/google/uuid"
	"github.com/jinzhu/copier"
)

func RoleToResponse(role *entity.Role) *model.RoleResponse {
	var response model.RoleResponse
	copier.Copy(&response, role)
	return &response
}

func CreateRoleRequestToEntity(request *model.CreateRoleRequest) *entity.Role {
	var role entity.Role
	copier.CopyWithOption(&role, request, copier.Option{
		DeepCopy: true,
	})
	role.ID = uuid.NewString()
	return &role
}

func UpdateRoleRequestToEntity(existingRole *entity.Role, request *model.UpdateRoleRequest) *entity.Role {
	copier.CopyWithOption(existingRole, request, copier.Option{
		IgnoreEmpty: true,
		DeepCopy:    true,
	})
	return existingRole
}
