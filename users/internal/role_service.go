package internal

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/madrabit/mini-market/users/internal/common"
)

type RoleService struct {
	repo      RoleRepo
	validator Validator
}

type RoleRepo interface {
	CreateRole(role Role) error
	UpdateRole(id uuid.UUID, name string) error
	DeleteRole(id uuid.UUID) error
	GetAllRoles() ([]Role, error)
	GetRoleByName(name string) (Role, error)
	GetUsersByRole(role string) ([]User, error)
}

func NewRoleService(repo RoleRepo, validator Validator) *RoleService {
	return &RoleService{
		repo:      repo,
		validator: validator,
	}
}

func (s *RoleService) CreateRole(req CreateRoleReq) error {
	if err := s.validator.Validate(req); err != nil {
		return &common.NotFoundError{Message: err.Error()}
	}
	id := uuid.New()
	role := Role{
		Id:   id,
		Name: req.Name,
	}
	err := s.repo.CreateRole(role)
	if err != nil {
		return fmt.Errorf("role service: failed to create role: %w", err)
	}
	return nil
}

func (s *RoleService) UpdateRole(id uuid.UUID, req UpdateRoleReq) error {
	if err := s.validator.Validate(req); err != nil {
		return &common.RequestValidationError{Message: err.Error()}
	}
	err := s.repo.UpdateRole(id, req.Name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return &common.RequestValidationError{Message: "role not found"}
		}
		return fmt.Errorf("role service: failed to update role: %w", err)
	}
	return nil
}

func (s *RoleService) DeleteRole(id uuid.UUID) error {
	err := s.repo.DeleteRole(id)
	if err != nil {
		return fmt.Errorf("role service: failed to delete role: %w", err)
	}
	return nil
}

func (s *RoleService) GetRoleByName(name string) (Role, error) {
	role, err := s.repo.GetRoleByName(name)
	if err != nil {
		return Role{}, fmt.Errorf("role service: failed to get role by role name: %w", err)
	}
	return role, nil
}

func (s *RoleService) GetAllRoles() ([]Role, error) {
	roles, err := s.repo.GetAllRoles()
	if err != nil {
		return nil, fmt.Errorf("role service: failed to get all roles: %w", err)
	}
	return roles, nil
}

func (s *RoleService) GetUsersByRole(name string) (ListUsersResponse, error) {
	users, err := s.repo.GetUsersByRole(name)
	if err != nil {
		return ListUsersResponse{}, fmt.Errorf("role service: failed to get users by role: %w", err)
	}
	userResp := make([]UserResponse, 0, len(users))
	for _, u := range users {
		userResp = append(userResp, UserResponse{
			ID:    u.Id,
			Name:  u.Name,
			Email: u.Email,
		})
	}
	return ListUsersResponse{
		userResp,
	}, nil
}
