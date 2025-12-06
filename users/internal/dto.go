package internal

import "github.com/google/uuid"

type User struct {
	Id           uuid.UUID
	Name         string
	Email        string
	PasswordHash string
	Roles        []Role
}

type CreateUserReq struct {
	Name     string
	Email    string
	Password string
	Role     uuid.UUID `json:"role_id"`
}

type UpdateUserReq struct {
	Name     string
	Email    string
	Password string
	Role     Role
}

type DeleteUserReq struct {
	UserID uuid.UUID
}

type SetUserRoleReq struct {
	UserID uuid.UUID
	Role   Role
}

type UserResponse struct {
	ID    uuid.UUID `json:"id"`
	Name  string    `json:"name"`
	Email string    `json:"email"`
	Role  []Role    `json:"role"`
}

type ListUsersRequest struct {
	IDs []uuid.UUID `json:"ids"`
}

type ListUsersResponse struct {
	Users []UserResponse `json:"users"`
}

type Role struct {
	Id   uuid.UUID
	Name string
}

type CreateRoleReq struct {
	Name string
}

type UpdateRoleReq struct {
	Name string
}

type DeleteRoleReq struct {
	RoleID uuid.UUID
}

type RoleResponse struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}
