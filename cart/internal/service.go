package internal

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/google/uuid"
	_ "github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/madrabit/mini-market/cart/internal/common"
)

type Service struct {
	repo      Repo
	validator Validator
}

type Repo interface {
	BeginTransaction() (tx *sqlx.Tx, err error)
	FindItemById(tx *sqlx.Tx, productID uuid.UUID) (bool, error)
	GetCart(userID uuid.UUID) (Cart, error)
	AddToCart(tx *sqlx.Tx, item AddToCartRequest) error
	UpdateCart(tx *sqlx.Tx, item UpdateCartItemRequest) error
	DeleteProduct(id uuid.UUID) error
}

type Validator interface {
	Validate(request any) error
}

func NewService(repo Repo, validator Validator) *Service {
	return &Service{repo, validator}
}

func (s *Service) Add(item AddToCartRequest) (err error) {
	if err = s.validator.Validate(item); err != nil {
		return &common.RequestValidationError{Message: err.Error()}
	}
	tx, err := s.repo.BeginTransaction()
	if err != nil {
		return fmt.Errorf("cart service: add product: error starting transaction")
	}
	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			err = fmt.Errorf("cart service: add product:  panic add product: %v", p)
			return
		}
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				err = fmt.Errorf("rollback failed: original error: %w", err)
			}
			return
		}
		if commitErr := tx.Commit(); commitErr != nil {
			err = fmt.Errorf("cart service: add product: committing transaction failed: %w", commitErr)
		}
	}()
	isExists, err := s.repo.FindItemById(tx, item.ProductId)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("cart service: add product: error checking exists of product")
	}
	if isExists {
		return &common.AlreadyExistsError{Message: fmt.Sprintf("product with id %s already exists", item.ProductId)}
	}
	err = s.repo.AddToCart(tx, item)
	if err != nil {
		return fmt.Errorf("cart service: add product: error adding product")
	}
	return nil
}

func (s *Service) GetCart(userID uuid.UUID) (Cart, error) {
	if userID == uuid.Nil {
		return Cart{}, errors.New("cart service: invalid user")
	}
	cart, err := s.repo.GetCart(userID)
	if err != nil {
		return Cart{}, fmt.Errorf("cart service: failed to get cart: %w", err)
	}
	return cart, nil
}

func (s *Service) UpdateCart(item UpdateCartItemRequest) error {
	if err := s.validator.Validate(item); err != nil {
		return &common.RequestValidationError{Message: err.Error()}
	}
	tx, err := s.repo.BeginTransaction()
	if err != nil {
		return fmt.Errorf("cart service: update cart: error starting transaction")
	}
	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			err = fmt.Errorf("cart service: update cart: panic update cart: %v", p)
			return
		}
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				err = fmt.Errorf("rollback failed: original error: %w", err)
			}
			return
		}
		if commitErr := tx.Commit(); commitErr != nil {
			err = fmt.Errorf("cart service: update cart: committing transaction failed: %w", commitErr)
		}
	}()
	isExists, err := s.repo.FindItemById(tx, item.ProductId)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("cart service: add product: error checking exists of product")
	}
	if isExists {
		return &common.AlreadyExistsError{Message: fmt.Sprintf("product with id %s already exists", item.ProductId)}
	}
	err = s.repo.UpdateCart(tx, item)
	if err != nil {
		return fmt.Errorf("cart service: add product: error adding product")
	}
	return nil
}

func (s *Service) Delete(id uuid.UUID) error {
	err := s.repo.DeleteProduct(id)
	if err != nil {
		return fmt.Errorf("cart service: delete: error deleting product with id %d", id)
	}
	return nil
}
