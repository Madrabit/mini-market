package internal

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/madrabit/mini-market/payment/internal/common"
)

type Service struct {
	repo      Repo
	validator Validator
}

type Repo interface {
	BeginTransaction() (tx *sqlx.Tx, err error)
	FindItemById(tx *sqlx.Tx, productID uuid.UUID) (bool, error)
	CreateOrder(tx *sqlx.Tx, req PaymentRequest) (CreatePaymentResponse, error)
	PSPWebhook(tx *sqlx.Tx, req PSPWebhookRequest) error
	GetStatus(userID, orderID uuid.UUID) (PaymentStatusResponse, error)
}

type Validator interface {
	Validate(request any) error
}

func NewService(repo Repo, validator Validator) *Service {
	return &Service{repo, validator}
}

func (s *Service) CreateOrder(req PaymentRequest) (CreatePaymentResponse, error) {
	if err := s.validator.Validate(req); err != nil {
		return CreatePaymentResponse{}, &common.RequestValidationError{Message: err.Error()}
	}
	tx, err := s.repo.BeginTransaction()
	if err != nil {
		return CreatePaymentResponse{}, fmt.Errorf("payment service: create order: error starting transaction")
	}
	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			err = fmt.Errorf("payment service: create order: panic create order: %v", p)
			return
		}
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				err = fmt.Errorf("rollback failed: original error: %w", err)
			}
			return
		}
		if commitErr := tx.Commit(); commitErr != nil {
			err = fmt.Errorf("payment service: create order: committing transaction failed: %w", commitErr)
		}
	}()
	order := CreatePaymentResponse{
		PaymentID: uuid.New(),
		Status:    Pending,
		Amount:    req.Amount,
		Currency:  req.Currency,
	}
	isExists, err := s.repo.FindItemById(tx, order.PaymentID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return CreatePaymentResponse{}, fmt.Errorf("payment service: create order: error checking exists of order")
	}
	if isExists {
		return CreatePaymentResponse{}, &common.AlreadyExistsError{Message: fmt.Sprintf("payment with id %s already exists", order.PaymentID)}
	}
	order, err = s.repo.CreateOrder(tx, req)
	if err != nil {
		return CreatePaymentResponse{}, fmt.Errorf("payment service: create order: error adding order")
	}
	return order, nil
}

func (s *Service) GetStatus(userID, orderID uuid.UUID) (PaymentStatusResponse, error) {
	if userID == uuid.Nil || orderID == uuid.Nil {
		return PaymentStatusResponse{}, errors.New("payment service: get status: invalid id")
	}
	status, err := s.repo.GetStatus(userID, orderID)
	if err != nil {
		return PaymentStatusResponse{}, fmt.Errorf("payment service: get status: failed to get payment status: %w", err)
	}
	return status, nil
}

func (s *Service) PSPWebhook(req PSPWebhookRequest) (PaymentStatusResponse, error) {
	if err := s.validator.Validate(req); err != nil {
		return PaymentStatusResponse{}, &common.RequestValidationError{Message: err.Error()}
	}
	tx, err := s.repo.BeginTransaction()
	if err != nil {
		return PaymentStatusResponse{}, fmt.Errorf("payment service: pspwebhook: error starting transaction")
	}
	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			err = fmt.Errorf("payment service: pspwebhook: panic creat psp webhook: %v", p)
			return
		}
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				err = fmt.Errorf("rollback failed: original error: %w", err)
			}
			return
		}
		if commitErr := tx.Commit(); commitErr != nil {
			err = fmt.Errorf("payment service: pspwebhook: committing transaction failed: %w", commitErr)
		}
	}()
	resp := PaymentStatusResponse{
		PaymentID: uuid.New(),
		Status:    Authorized,
		Amount:    req.Amount,
		Currency:  req.Currency,
	}
	err = s.repo.PSPWebhook(tx, req)
	if err != nil {
		return PaymentStatusResponse{}, fmt.Errorf("payment service: pspwebhook: error create pspwebhook")
	}
	return resp, nil
}
