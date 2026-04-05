package entity

import "time"

type User struct {
	ID          string
	Email       string
	Password    string
	FirstName   string
	LastName    string
	PhonePrefix string
	PhoneNumber string
	RoleID      string
	Role        *Role
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time
}

type UserQuery struct {
	Page    int
	Limit   int
	Search  string
	RoleID  string
	SortBy  string
	SortDir string
}
