package internal

import "github.com/google/uuid"

type User struct {
	Id           uuid.UUID
	Name         string
	Email        string
	PasswordHash string
	Role         Role
}

type Role struct {
	Id   uuid.UUID
	Name string
}

type CreateUserReq struct {
	Name     string
	Email    string
	Password string
	Role     Role
}

type UpdateUserReq struct {
	UserID   uuid.UUID
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
	Role  string    `json:"role"`
}

type ListUsersResponse struct {
	Users []UserResponse `json:"users"`
}
