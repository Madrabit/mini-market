package internal

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/madrabit/mini-market/order/internal/common"
	"time"
)

type Service struct {
	repo      Repo
	validator Validator
}

type Repo interface {
	BeginTransaction() (tx *sqlx.Tx, err error)
	FindItemById(tx *sqlx.Tx, productID uuid.UUID) (bool, error)
	CreateOrder(tx *sqlx.Tx, req CreatOrderRequest) error
	GetStatus(user, order uuid.UUID) (StatusResponse, error)
	UpdatePaymentStatus(req UpdatePaymentStatusRequest) error
}

type Validator interface {
	Validate(request any) error
}

func NewService(repo Repo, validator Validator) *Service {
	return &Service{repo, validator}
}

func (s *Service) CreateOrder(req CreatOrderRequest) (OrderResponse, error) {
	if err := s.validator.Validate(req); err != nil {
		return OrderResponse{}, &common.RequestValidationError{Message: err.Error()}
	}
	tx, err := s.repo.BeginTransaction()
	if err != nil {
		return OrderResponse{}, fmt.Errorf("order service: create order: error starting transaction")
	}
	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			err = fmt.Errorf("order service: create order: panic add product: %v", p)
			return
		}
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				err = fmt.Errorf("rollback failed: original error: %w", err)
			}
			return
		}
		if commitErr := tx.Commit(); commitErr != nil {
			err = fmt.Errorf("order service: create order: committing transaction failed: %w", commitErr)
		}
	}()
	order := OrderResponse{
		ID:         uuid.New(),
		UserId:     req.UserID,
		Status:     New,
		GrandTotal: 0,
		Created:    time.Now(),
		Items:      []ItemResponse{},
	}
	isExists, err := s.repo.FindItemById(tx, order.ID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return OrderResponse{}, fmt.Errorf("order service: create order: error checking exists of order")
	}
	if isExists {
		return OrderResponse{}, &common.AlreadyExistsError{Message: fmt.Sprintf("order with id %s already exists", order.ID)}
	}
	err = s.repo.CreateOrder(tx, req)
	if err != nil {
		return OrderResponse{}, fmt.Errorf("order service: create order: error adding order")
	}
	return order, nil
}

func (s *Service) GetStatus(user, order uuid.UUID) (StatusResponse, error) {
	if user == uuid.Nil || order == uuid.Nil {
		return StatusResponse{}, errors.New("order service: get status: invalid id")
	}
	status, err := s.repo.GetStatus(user, order)
	if err != nil {
		return StatusResponse{}, fmt.Errorf("order service: get status: failed to get order status: %w", err)
	}
	return status, nil
}

func (s *Service) UpdatePaymentStatus(req UpdatePaymentStatusRequest) error {
	if err := s.validator.Validate(req); err != nil {
		return &common.RequestValidationError{Message: err.Error()}
	}
	err := s.repo.UpdatePaymentStatus(req)
	if err != nil {
		return fmt.Errorf("order service: update payment status: error update product")
	}
	return nil
}
