package internal

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/madrabit/mini-market/rating/internal/common"
)

type Service struct {
	repo      Repo
	validator Validator
}

type Repo interface {
	BeginTransaction() (tx *sqlx.Tx, err error)
	FindItemById(tx *sqlx.Tx, productID uuid.UUID) (bool, error)
	AddReview(tx *sqlx.Tx, review AddReviewRequest) error
	UpdateReview(item UpdateReviewRequest) error
	DeleteReview(user, order uuid.UUID) error
	GetReviewsByProduct(productID uuid.UUID) (ProductReviewsResponse, error)
}

type Validator interface {
	Validate(request any) error
}

func NewService(repo Repo, validator Validator) *Service {
	return &Service{repo, validator}
}

func (s *Service) AddProduct(review AddReviewRequest) (err error) {
	if err = s.validator.Validate(review); err != nil {
		return &common.RequestValidationError{Message: err.Error()}
	}
	tx, err := s.repo.BeginTransaction()
	if err != nil {
		return fmt.Errorf("review service: add review: error starting transaction")
	}
	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			err = fmt.Errorf("review service: add review: panic add review: %v", p)
			return
		}
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				err = fmt.Errorf("rollback failed: original error: %w", err)
			}
			return
		}
		if commitErr := tx.Commit(); commitErr != nil {
			err = fmt.Errorf("review service: add review: committing transaction failed: %w", commitErr)
		}
	}()
	isExists, err := s.repo.FindItemById(tx, review.ProductID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("review service: add review: error checking exists of review")
	}
	if isExists {
		return &common.AlreadyExistsError{Message: fmt.Sprintf("review with id %s already exists", review.ProductID)}
	}
	err = s.repo.AddReview(tx, review)
	if err != nil {
		return fmt.Errorf("review service: add review:: error adding review")
	}
	return nil
}
