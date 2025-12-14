package internal

import "github.com/google/uuid"

type User struct {
	Id           uuid.UUID `json:"id" db:"id"`
	Name         string    `json:"name"  db:"name"`
	Email        string    `json:"email" db:"email"`
	PasswordHash string    `json:"-" db:"password_hash"`
	Roles        []Role    `json:"roles" db:"-"`
}

type CreateUserReq struct {
	Name     string `json:"name" validate:"required,min=2"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

type UpdateUserReq struct {
	Name  string `json:"name" validate:"required,min=2"`
	Email string `json:"email" validate:"required,email"`
}

type UserResponse struct {
	ID    uuid.UUID `json:"id"`
	Name  string    `json:"name"`
	Email string    `json:"email"`
}

type ListUsersRequest struct {
	IDs []uuid.UUID `json:"ids" validate:"required,gt=1"`
}

type ListUsersResponse struct {
	Users []UserResponse `json:"users"`
}

type Role struct {
	Id   uuid.UUID `json:"id" db:"id"`
	Name string    `json:"name" db:"name"`
}

type CreateRoleReq struct {
	Name string `json:"name" validate:"required,min=2"`
}

type UpdateRoleReq struct {
	Name string `json:"name" validate:"required,min=2"`
}

type RoleResponse struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}
