package internal

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/madrabit/mini-market/users/internal/common"
)

type Service struct {
	repo      Repo
	validator Validator
}

type Repo interface {
	CreateUser(req CreateUserReq) error
	UpdateUser(req UpdateUserReq) error
	DeleteUser(DeleteUserReq uuid.UUID) error
	ChangeRole(req SetUserRoleReq) error
	GetUserByID(userID uuid.UUID) (User, error)
	GetUsersByIds(IDs ListUsersRequest) (ListUsersResponse, error)
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
	err := s.repo.CreateUser(req)
	if err != nil {
		return fmt.Errorf("user service: failed to create user: %w", err)
	}
	return nil
}

func (s *Service) UpdateUser(req UpdateUserReq) error {
	if err := s.validator.Validate(req); err != nil {
		return &common.RequestValidationError{Message: err.Error()}
	}
	err := s.repo.UpdateUser(req)
	if err != nil {
		return fmt.Errorf("user service: failed to update user: %w", err)
	}
	return nil
}

func (s *Service) DeleteUser(DeleteUserReq uuid.UUID) error {
	err := s.repo.DeleteUser(DeleteUserReq)
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
	users, err := s.repo.GetUsersByIds(IDs)
	if err != nil {
		return ListUsersResponse{}, fmt.Errorf("user service: failed to get user by IDs: %w", err)
	}
	return users, nil
}
