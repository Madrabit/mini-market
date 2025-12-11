package internal

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/madrabit/mini-market/users/internal/common"
)

type UserService struct {
	userRepo  UserRepo
	roleRepo  RoleRepo
	validator Validator
}

type UserRepo interface {
	BeginTransaction() (*sqlx.Tx, error)
	CreateUser(tx *sqlx.Tx, user User) error
	AddUserRoles(tx *sqlx.Tx, userId uuid.UUID, roles []uuid.UUID) error
	UpdateUser(user User) error
	DeleteUser(userID uuid.UUID) error
	GetUserByID(userID uuid.UUID) (User, error)
	GetUsersByIds(IDs []uuid.UUID) ([]User, error)
	GetUsersByRole(role string) ([]User, error)
	GetUserRoles(userID uuid.UUID) ([]Role, error)
}

func NewUserService(userRepo UserRepo, roleRepo RoleRepo, validator Validator) *UserService {
	return &UserService{
		userRepo:  userRepo,
		roleRepo:  roleRepo,
		validator: validator,
	}
}

func (s *UserService) CreateUser(req CreateUserReq) (err error) {
	if err := s.validator.Validate(req); err != nil {
		return &common.RequestValidationError{Message: err.Error()}
	}
	id := uuid.New()
	password, err := common.HashPassword(req.Password)
	if err != nil {
		return fmt.Errorf("user service: failed to hash password: %w", err)
	}
	const defaultRoleName = "basic"
	role, err := s.roleRepo.GetRoleByName(defaultRoleName)
	if err != nil {
		return fmt.Errorf("user service: create user: failed to get role by name: %w", err)
	}
	user := User{
		Id:           id,
		Name:         req.Name,
		Email:        req.Email,
		PasswordHash: password,
		Roles:        []Role{role},
	}
	tx, err := s.userRepo.BeginTransaction()
	if err != nil {
		return fmt.Errorf("user service: create user: error starting transaction: %w", err)
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()
	err = s.userRepo.CreateUser(tx, user)
	if err != nil {
		return fmt.Errorf("user service: failed to create user: %w", err)
	}
	err = s.userRepo.AddUserRoles(tx, user.Id, []uuid.UUID{role.Id})
	if err != nil {
		return fmt.Errorf("user service: failed add roles to user: %w", err)
	}
	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("user service: failed to commit transaction: %w", err)
	}
	return nil
}

func (s *UserService) UpdateUser(id uuid.UUID, req UpdateUserReq) error {
	if err := s.validator.Validate(req); err != nil {
		return &common.RequestValidationError{Message: err.Error()}
	}
	user := User{
		Id:    id,
		Name:  req.Name,
		Email: req.Email,
	}
	err := s.userRepo.UpdateUser(user)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("user service: user not found: %w", err)
		}
		return fmt.Errorf("user service: failed to update user: %w", err)
	}
	return nil
}

func (s *UserService) DeleteUser(userID uuid.UUID) error {
	if err := s.userRepo.DeleteUser(userID); err != nil {
		return fmt.Errorf("user service: failed to delete user: %w", err)
	}
	return nil
}

func (s *UserService) GetUserByID(userID uuid.UUID) (User, error) {
	user, err := s.userRepo.GetUserByID(userID)
	if err != nil {
		return User{}, fmt.Errorf("user service: failed to get user by ID: %w", err)
	}
	return user, nil
}

func (s *UserService) GetUserRoles(userID uuid.UUID) ([]Role, error) {
	roles, err := s.userRepo.GetUserRoles(userID)
	if err != nil {
		return nil, fmt.Errorf("user service: failed to get user roles: %w", err)
	}
	return roles, nil
}

func (s *UserService) GetUsersByIds(IDs ListUsersRequest) (ListUsersResponse, error) {
	if err := s.validator.Validate(IDs); err != nil {
		return ListUsersResponse{}, &common.RequestValidationError{Message: err.Error()}
	}
	users, err := s.userRepo.GetUsersByIds(IDs.IDs)
	if err != nil {
		return ListUsersResponse{}, fmt.Errorf("user service: failed to get user by IDs: %w", err)
	}
	uResp := make([]UserResponse, 0, len(users))
	for _, user := range users {
		uResp = append(uResp, UserResponse{
			ID:    user.Id,
			Name:  user.Name,
			Email: user.Email,
		})
	}
	response := ListUsersResponse{
		uResp,
	}
	return response, nil
}

func (s *UserService) GetUsersByRole(role string) (ListUsersResponse, error) {
	users, err := s.userRepo.GetUsersByRole(role)
	if err != nil {
		return ListUsersResponse{}, fmt.Errorf("user service: failed to get users by role: %w", err)
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
