package model

import "time"

type UserResponse struct {
	ID          string        `json:"id"`
	FirstName   string        `json:"first_name"`
	LastName    string        `json:"last_name"`
	Email       string        `json:"email"`
	PhonePrefix string        `json:"phone_prefix"`
	PhoneNumber string        `json:"phone_number"`
	RoleID      string        `json:"role_id"`
	Role        *RoleResponse `json:"role,omitempty"`
	CreatedAt   time.Time     `json:"created_at"`
	UpdatedAt   time.Time     `json:"updated_at"`
	DeletedAt   *time.Time    `json:"deleted_at"`
}

type CreateUserRequest struct {
	FirstName   string `json:"first_name" validate:"required"`
	LastName    string `json:"last_name" validate:"required"`
	Email       string `json:"email" validate:"required,email"`
	Password    string `json:"password" validate:"required,min=8"`
	PhonePrefix string `json:"phone_prefix" validate:"required"`
	PhoneNumber string `json:"phone_number" validate:"required,numeric"`
	RoleID      string `json:"role_id" validate:"required"`
}

type UpdateUserRequest struct {
	ID          string  `json:"-"`
	FirstName   *string `json:"first_name,omitempty"`
	LastName    *string `json:"last_name,omitempty"`
	Email       *string `json:"email,omitempty" validate:"omitempty,email"`
	PhonePrefix *string `json:"phone_prefix,omitempty"`
	PhoneNumber *string `json:"phone_number,omitempty" validate:"omitempty,numeric"`
}

type GetAllUsersRequest struct {
	Page    int    `query:"page" validate:"required,min=1"`
	Limit   int    `query:"limit" validate:"required,min=1"`
	Search  string `query:"search"`
	RoleID  string `query:"role_id"`
	SortBy  string `query:"sort_by"`
	SortDir string `query:"sort_dir" validate:"omitempty,oneof=asc desc"`
}
