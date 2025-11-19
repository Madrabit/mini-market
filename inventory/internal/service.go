package internal

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/madrabit/mini-market/inventory/internal/common"
)

type Service struct {
	repo      Repo
	validator Validator
}

type Repo interface {
	BeginTransaction() (tx *sqlx.Tx, err error)
	FindItemById(tx *sqlx.Tx, productID uuid.UUID) (bool, error)
	GetProductsByIds(IDs []uuid.UUID) (ListItemsResponse, error)
	AddProduct(tx *sqlx.Tx, item AddItemRequest) error
	UpdateProduct(tx *sqlx.Tx, item UpdateItemRequest) error
	DeleteProduct(id uuid.UUID) error
	GetProductById(id uuid.UUID) (Item, error)
	ReserveProducts(tx *sqlx.Tx, item ReserveItemRequest) error
	ReleaseProducts(tx *sqlx.Tx, item ReliesItemRequest) error
}

type Validator interface {
	Validate(request any) error
}

func NewService(repo Repo, validator Validator) *Service {
	return &Service{repo, validator}
}

func (s *Service) AddProduct(item AddItemRequest) (err error) {
	if err = s.validator.Validate(item); err != nil {
		return &common.RequestValidationError{Message: err.Error()}
	}
	tx, err := s.repo.BeginTransaction()
	if err != nil {
		return fmt.Errorf("inventory service: add product: error starting transaction")
	}
	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			err = fmt.Errorf("inventory service: add product:  panic add product: %v", p)
			return
		}
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				err = fmt.Errorf("rollback failed: original error: %w", err)
			}
			return
		}
		if commitErr := tx.Commit(); commitErr != nil {
			err = fmt.Errorf("inventory service: add product: committing transaction failed: %w", commitErr)
		}
	}()
	isExists, err := s.repo.FindItemById(tx, item.Id)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("inventory service: add product: error checking exists of product")
	}
	if isExists {
		return &common.AlreadyExistsError{Message: fmt.Sprintf("product with id %s already exists", item.Id)}
	}
	err = s.repo.AddProduct(tx, item)
	if err != nil {
		return fmt.Errorf("inventory service: add product: error adding product")
	}
	return nil
}

func (s *Service) GetProductsByIds(IDs []uuid.UUID) (ListItemsResponse, error) {
	if len(IDs) == 0 {
		return ListItemsResponse{}, errors.New("inventory service: empty IDs provided")
	}
	product, err := s.repo.GetProductsByIds(IDs)
	if err != nil {
		return ListItemsResponse{}, fmt.Errorf("inventory service: failed to get product by ids: %w", err)
	}
	return product, nil
}

func (s *Service) GetProductById(productId uuid.UUID) (Item, error) {
	if productId == uuid.Nil {
		return Item{}, errors.New("inventory service: invalid id")
	}
	product, err := s.repo.GetProductById(productId)
	if err != nil {
		return Item{}, fmt.Errorf("inventory service: failed to get product by id: %w", err)
	}
	return product, nil
}

func (s *Service) UpdateProduct(item UpdateItemRequest) error {
	if err := s.validator.Validate(item); err != nil {
		return &common.RequestValidationError{Message: err.Error()}
	}
	tx, err := s.repo.BeginTransaction()
	if err != nil {
		return fmt.Errorf("inventory service: update product: error starting transaction")
	}
	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			err = fmt.Errorf("inventory service: update product: panic update produdct: %v", p)
			return
		}
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				err = fmt.Errorf("rollback failed: original error: %w", err)
			}
			return
		}
		if commitErr := tx.Commit(); commitErr != nil {
			err = fmt.Errorf("inventory service: update product: committing transaction failed: %w", commitErr)
		}
	}()
	isExists, err := s.repo.FindItemById(tx, item.Id)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("catalog service: update product: error checking exists of product")
	}
	if isExists {
		return &common.AlreadyExistsError{Message: fmt.Sprintf("product with id %s already exists", item.Id)}
	}
	err = s.repo.UpdateProduct(tx, item)
	if err != nil {
		return fmt.Errorf("inventory service: update product: error update product")
	}
	return nil
}

func (s *Service) DeleteProduct(id uuid.UUID) error {
	err := s.repo.DeleteProduct(id)
	if err != nil {
		return fmt.Errorf("inventory service: delete: error deleting product with id %d", id)
	}
	return nil
}

func (s *Service) ReserveProducts(item ReserveItemRequest) error {
	if err := s.validator.Validate(item); err != nil {
		return &common.RequestValidationError{Message: err.Error()}
	}
	tx, err := s.repo.BeginTransaction()
	if err != nil {
		return fmt.Errorf("inventory service: reserve product: error starting transaction")
	}
	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			err = fmt.Errorf("inventory service: reserve product:  panic reserve product: %v", p)
			return
		}
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				err = fmt.Errorf("rollback failed: original error: %w", err)
			}
			return
		}
		if commitErr := tx.Commit(); commitErr != nil {
			err = fmt.Errorf("inventory service: reserve product: committing transaction failed: %w", commitErr)
		}
	}()
	isExists, err := s.repo.FindItemById(tx, item.Id)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("inventory service: reserve product: error checking exists of product")
	}
	if isExists {
		return &common.AlreadyExistsError{Message: fmt.Sprintf("product with id %s already exists", item.Id)}
	}
	err = s.repo.ReserveProducts(tx, item)
	if err != nil {
		return fmt.Errorf("inventory service: reserve product: error adding product")
	}
	return nil
}

func (s *Service) ReleaseProducts(item ReliesItemRequest) error {
	if err := s.validator.Validate(item); err != nil {
		return &common.RequestValidationError{Message: err.Error()}
	}
	tx, err := s.repo.BeginTransaction()
	if err != nil {
		return fmt.Errorf("inventory service: release product: error starting transaction")
	}
	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			err = fmt.Errorf("inventory service: release product:  panic add product: %v", p)
			return
		}
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				err = fmt.Errorf("rollback failed: original error: %w", err)
			}
			return
		}
		if commitErr := tx.Commit(); commitErr != nil {
			err = fmt.Errorf("inventory service: release product: committing transaction failed: %w", commitErr)
		}
	}()
	isExists, err := s.repo.FindItemById(tx, item.Id)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("inventory service: release product: error checking exists of product")
	}
	if isExists {
		return &common.AlreadyExistsError{Message: fmt.Sprintf("product with id %s already exists", item.Id)}
	}
	err = s.repo.ReleaseProducts(tx, item)
	if err != nil {
		return fmt.Errorf("inventory service: release product: error adding product")
	}
	return nil
}
