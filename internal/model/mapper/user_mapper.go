package mapper

import (
	"github.com/Primuse-Pte-Ltd/go-boilerplate-clean-architecture/internal/entity"
	"github.com/Primuse-Pte-Ltd/go-boilerplate-clean-architecture/internal/model"
	"github.com/google/uuid"
	"github.com/jinzhu/copier"
)

func UserToResponse(user *entity.User) *model.UserResponse {
	var response model.UserResponse
	copier.Copy(&response, user)
	return &response
}

func CreateUserRequestToEntity(request *model.CreateUserRequest) *entity.User {
	var user entity.User
	copier.CopyWithOption(&user, request, copier.Option{
		DeepCopy: true,
	})
	user.ID = uuid.NewString()

	return &user
}

func UpdateUserRequestToEntity(existingUser *entity.User, request *model.UpdateUserRequest) *entity.User {
	copier.CopyWithOption(existingUser, request, copier.Option{
		IgnoreEmpty: true,
		DeepCopy:    true,
	})
	return existingUser
}
