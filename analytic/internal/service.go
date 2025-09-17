package internal

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/madrabit/mini-market/analytic/internal/common"
)

type Service struct {
	repo      Repo
	validator Validator
}

type Repo interface {
	FindEventById(tx *sqlx.Tx, eventID uuid.UUID) (bool, error)
	BeginTransaction() (tx *sqlx.Tx, err error)
	AddEvent(tx *sqlx.Tx, event Event) error
	GetOrdersSummary(from, to string) (OrdersSummary, error)
	GetTopProducts(limit int64) (TopProducts, error)
	GetDailySails(from, to string) (DailySales, error)
	GetSearchTrends(limit int64) (SearchTrends, error)
	GetFailSearch(limit int64) (FailedSearch, error)
	GetAvgRatingByProduct(productID uuid.UUID) (ProductAvgRating, error)
	GetTopProductsByRating(limit int64) (TopProducts, error)
}

type Validator interface {
	Validate(request any) error
}

func NewService(repo Repo, validator Validator) *Service {
	return &Service{repo, validator}
}

func (s *Service) Add(event Event) (err error) {
	if err = s.validator.Validate(event); err != nil {
		return &common.RequestValidationError{Message: err.Error()}
	}
	tx, err := s.repo.BeginTransaction()
	if err != nil {
		return fmt.Errorf("analytic service: add event: error starting transaction")
	}
	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			err = fmt.Errorf("analytic service: add event:  panic add employee: %v", p)
			return
		}
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				err = fmt.Errorf("rollback failed: original error: %w", err)
			}
			return
		}
		if commitErr := tx.Commit(); commitErr != nil {
			err = fmt.Errorf("analytic service: add event: committing transaction failed: %w", commitErr)
		}
	}()
	isExists, err := s.repo.FindEventById(tx, event.ID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("analytic service: add event: error checking exists of event")
	}
	if isExists {
		return &common.AlreadyExistsError{Message: fmt.Sprintf("event with id %s already exists", event.ID)}
	}
	err = s.repo.AddEvent(tx, event)
	if err != nil {
		return fmt.Errorf("analytic service: add event: error adding event")
	}
	return nil
}
