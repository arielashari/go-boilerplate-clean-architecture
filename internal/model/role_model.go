package model

import "time"

type RoleResponse struct {
	ID        string     `json:"id"`
	Name      string     `json:"name"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}

type CreateRoleRequest struct {
	Name string `json:"name" validate:"required"`
}

type UpdateRoleRequest struct {
	ID   string  `json:"-"`
	Name *string `json:"name,omitempty"`
}
