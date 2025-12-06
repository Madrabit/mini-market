package internal

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/madrabit/mini-market/users/internal/common"
)

type Service struct {
	repo      Repo
	validator Validator
}

type Repo interface {
	CreateUser(tx *sqlx.Tx, user User) error
	AddUserRoles(tx *sqlx.Tx, userId uuid.UUID, roles []uuid.UUID) error
	UpdateUser(user User) error
	DeleteUser(userID uuid.UUID) error
	ChangeRole(req SetUserRoleReq) error
	GetUserByID(userID uuid.UUID) (User, error)
	GetUsersByIds(IDs []uuid.UUID) ([]User, error)
	CreateRole(role Role) error
	UpdateRole(id uuid.UUID, role string) error
	DeleteRole(id uuid.UUID) error
	GetUsersByRole(role string) ([]User, error)
	GetAllRoles() ([]Role, error)
	GetRoleByName(role string) (Role, error)
	BeginTransaction() (tx *sqlx.Tx, err error)
}

type Validator interface {
	Validate(request any) error
}

func NewService(repo Repo, validator Validator) *Service {
	return &Service{repo, validator}
}

func (s *Service) CreateUser(req CreateUserReq) error {
	if err := s.validator.Validate(req); err != nil {
		return &common.RequestValidationError{Message: err.Error()}
	}
	id := uuid.New()
	password, err := common.HashPassword(req.Password)
	if err != nil {
		return fmt.Errorf("user service: failed to hash password: %w", err)
	}
	role, err := s.GetRoleByName("basic")
	if err != nil {
		return fmt.Errorf("user service: create user: failed to get role by name: %w", err)
	}
	user := User{
		id,
		req.Name,
		req.Email,
		password,
		[]Role{role},
	}
	tx, err := s.repo.BeginTransaction()
	if err != nil {
		return fmt.Errorf("user service: create user: error starting transaction")
	}
	defer tx.Rollback()
	err = s.repo.CreateUser(tx, user)
	if err != nil {
		return fmt.Errorf("user service: failed to create user: %w", err)
	}
	err = s.repo.AddUserRoles(tx, user.Id, []uuid.UUID{role.Id})
	if err != nil {
		return fmt.Errorf("user service: failed add roles to user: %w", err)
	}
	tx.Commit()
	return nil
}

func (s *Service) UpdateUser(id uuid.UUID, req UpdateUserReq) error {
	if err := s.validator.Validate(req); err != nil {
		return &common.RequestValidationError{Message: err.Error()}
	}
	passwordHash, err := common.HashPassword(req.Password)
	if err != nil {
		return fmt.Errorf("user service: failed to hash password: %w", err)
	}
	user := User{
		Id:           id,
		Name:         req.Name,
		Email:        req.Email,
		PasswordHash: passwordHash,
		//TODO доделать тут норм добавление ролей
		Roles: []Role{},
	}
	err = s.repo.UpdateUser(user)
	if err != nil {
		return fmt.Errorf("user service: failed to update user: %w", err)
	}
	return nil
}

func (s *Service) DeleteUser(userID uuid.UUID) error {
	err := s.repo.DeleteUser(userID)
	if err != nil {
		return fmt.Errorf("user service: failed to delete user: %w", err)
	}
	return nil
}

func (s *Service) ChangeRole(req SetUserRoleReq) error {
	if err := s.validator.Validate(req); err != nil {
		return &common.RequestValidationError{Message: err.Error()}
	}
	err := s.repo.ChangeRole(req)
	if err != nil {
		return fmt.Errorf("user service: failed to change user role: %w", err)
	}
	return nil
}

func (s *Service) GetUserByID(userID uuid.UUID) (User, error) {
	user, err := s.repo.GetUserByID(userID)
	if err != nil {
		return User{}, fmt.Errorf("user service: failed to get user by ID: %w", err)
	}
	return user, nil
}

func (s *Service) GetUsersByIds(IDs ListUsersRequest) (ListUsersResponse, error) {
	if err := s.validator.Validate(IDs.IDs); err != nil {
		return ListUsersResponse{}, &common.RequestValidationError{Message: err.Error()}
	}
	users, err := s.repo.GetUsersByIds(IDs.IDs)
	if err != nil {
		return ListUsersResponse{}, fmt.Errorf("user service: failed to get user by IDs: %w", err)
	}
	var uResp []UserResponse
	for _, user := range users {
		uResp = append(uResp, UserResponse{
			user.Id,
			user.Name,
			user.Email,
			user.Roles,
		})
	}
	response := ListUsersResponse{
		uResp,
	}
	return response, nil
}

//roles service

func (s *Service) CreateRole(req CreateRoleReq) error {
	if err := s.validator.Validate(req); err != nil {
		return &common.RequestValidationError{Message: err.Error()}
	}
	id := uuid.New()
	role := Role{
		id,
		req.Name,
	}
	err := s.repo.CreateRole(role)
	if err != nil {
		return fmt.Errorf("user service: failed to create role: %w", err)
	}
	return nil
}

func (s *Service) UpdateRole(id uuid.UUID, req UpdateRoleReq) error {
	if err := s.validator.Validate(req); err != nil {
		return &common.RequestValidationError{Message: err.Error()}
	}
	err := s.repo.UpdateRole(id, req.Name)
	if err != nil {
		return fmt.Errorf("user service: failed to update role: %w", err)
	}
	return nil
}

func (s *Service) DeleteRole(id uuid.UUID) error {
	err := s.repo.DeleteRole(id)
	if err != nil {
		return fmt.Errorf("user service: failed to delete role: %w", err)
	}
	return nil
}

func (s *Service) GetUsersByRole(role string) (ListUsersResponse, error) {
	users, err := s.repo.GetUsersByRole(role)
	if err != nil {
		return ListUsersResponse{}, fmt.Errorf("user service: failed to get users by role: %w", err)
	}
	var userResp []UserResponse
	for _, u := range users {
		userResp = append(userResp, UserResponse{
			u.Id,
			u.Name,
			u.Email,
			u.Roles,
		})
	}
	return ListUsersResponse{
		userResp,
	}, nil
}

func (s *Service) GetRoleByName(name string) (Role, error) {
	role, err := s.repo.GetRoleByName(name)
	if err != nil {
		return Role{}, fmt.Errorf("user service: failed to get role by roleName: %w", err)
	}
	return role, nil
}
