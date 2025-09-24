package internal

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/madrabit/mini-market/analytic/internal/common"
	"time"
)

type Service struct {
	repo      Repo
	validator Validator
}

type Repo interface {
	FindEventById(tx *sqlx.Tx, eventID uuid.UUID) (bool, error)
	BeginTransaction() (tx *sqlx.Tx, err error)
	AddEvent(tx *sqlx.Tx, event Event) error
	GetOrdersSummary(from, to time.Time) (OrdersSummary, error)
	GetTopProducts(limit int64) (TopProducts, error)
	GetDailySails(from, to time.Time) (DailySales, error)
	GetSearchTrends(limit int64) (SearchTrends, error)
	GetFailSearch(limit int64) (FailedSearch, error)
	GetAvgRatingByProduct(productID uuid.UUID) (ProductAvgRating, error)
	GetTopProductsByRating(limit int64) (TopProducts, error)
}

type Validator interface {
	Validate(request any) error
	ValidateVar(field any, tag string) error
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
			err = fmt.Errorf("analytic service: add event:  panic add event: %v", p)
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

const (
	dateLayout   = "2006-01-02"
	maxProducts  = 1000
	defaultLimit = 20
)

func (s *Service) GetOrdersSummary(from, to string) (OrdersSummary, error) {
	if err := s.validator.ValidateVar(from, "required,datetime="+dateLayout); err != nil {
		return OrdersSummary{}, &common.RequestValidationError{Message: "from" + err.Error()}
	}
	if err := s.validator.ValidateVar(to, "required,datetime=2006-01-02"); err != nil {
		return OrdersSummary{}, &common.RequestValidationError{Message: "to" + err.Error()}
	}
	dateFrom, err := time.Parse(dateLayout, from)
	if err != nil {
		return OrdersSummary{}, &common.RequestValidationError{Message: "from: invalid date"}
	}
	dateTo, err := time.Parse(dateLayout, to)
	if err != nil {
		return OrdersSummary{}, &common.RequestValidationError{Message: "from: invalid date"}
	}
	if dateFrom.After(dateTo) {
		return OrdersSummary{}, &common.RequestValidationError{Message: "date from > date to"}
	}
	// включить весь день
	dateFrom = dateFrom.Add(24*time.Hour - time.Nanosecond)
	ordersSummary, err := s.repo.GetOrdersSummary(dateFrom, dateTo)
	var nfError common.NotFoundError
	if err != nil {
		if errors.Is(err, &nfError) {
			return OrdersSummary{}, &common.NotFoundError{Message: "analytic service: get orders summary: order not found"}
		}
		return OrdersSummary{}, fmt.Errorf("analytic service: get orders summary: get orders summary: %w", err)
	}
	return ordersSummary, nil
}

func normalizeLimit(limit int64) int64 {
	if limit == 0 {
		return defaultLimit
	}
	if limit > maxProducts {
		limit = maxProducts
	}
	return limit
}

func (s *Service) GetTopProducts(limit int64) (TopProducts, error) {
	if err := s.validator.ValidateVar(limit, "gte=0"); err != nil {
		return TopProducts{}, &common.RequestValidationError{Message: "limit: " + err.Error()}
	}
	lim := normalizeLimit(limit)
	topProducts, err := s.repo.GetTopProducts(lim)
	if err != nil {
		return TopProducts{}, fmt.Errorf("analytic service: get top products: %w", err)
	}
	return topProducts, nil
}

func (s *Service) GetDailySails(from, to string) (DailySales, error) {
	if err := s.validator.ValidateVar(from, "required,datetime="+dateLayout); err != nil {
		return DailySales{}, &common.RequestValidationError{Message: "from" + err.Error()}
	}
	if err := s.validator.ValidateVar(to, "required,datetime="+dateLayout); err != nil {
		return DailySales{}, &common.RequestValidationError{Message: "to" + err.Error()}
	}
	dateFrom, err := time.Parse(dateLayout, from)
	if err != nil {
		return DailySales{}, &common.RequestValidationError{Message: "from: invalid date"}
	}
	dateTo, err := time.Parse(dateLayout, to)
	if err != nil {
		return DailySales{}, &common.RequestValidationError{Message: "from: invalid date"}
	}
	if dateFrom.After(dateTo) {
		return DailySales{}, &common.RequestValidationError{Message: "date from > date to"}
	}
	dateFrom = dateFrom.Add(24*time.Hour - time.Nanosecond)
	dailySales, err := s.repo.GetDailySails(dateFrom, dateTo)
	if err != nil {
		return DailySales{}, fmt.Errorf("analytic service: get daily sales %w", err)
	}
	return dailySales, nil
}

func (s *Service) GetSearchTrends(limit int64) (SearchTrends, error) {
	if err := s.validator.ValidateVar(limit, "gte=0"); err != nil {
		return SearchTrends{}, &common.RequestValidationError{Message: "limit: " + err.Error()}
	}
	lim := normalizeLimit(limit)
	trends, err := s.repo.GetSearchTrends(lim)
	if err != nil {
		return SearchTrends{}, fmt.Errorf("analytic service: get search trends: %w", err)
	}
	return trends, nil
}

func (s *Service) GetFailSearch(limit int64) (FailedSearch, error) {
	if err := s.validator.ValidateVar(limit, "gte=0"); err != nil {
		return FailedSearch{}, &common.RequestValidationError{Message: "limit: " + err.Error()}
	}
	lim := normalizeLimit(limit)
	searches, err := s.repo.GetFailSearch(lim)
	if err != nil {
		return FailedSearch{}, fmt.Errorf("analytic service: get failed searches: %w", err)
	}
	return searches, nil
}

func (s *Service) GetAvgRatingByProduct(productID uuid.UUID) (ProductAvgRating, error) {
	rating, err := s.repo.GetAvgRatingByProduct(productID)
	if err != nil {
		return ProductAvgRating{}, fmt.Errorf("analytic service: get avg rating by product %w", err)
	}
	return rating, nil
}

func (s *Service) GetTopProductsByRating(limit int64) (TopProducts, error) {
	if err := s.validator.ValidateVar(limit, "gte=0"); err != nil {
		return TopProducts{}, &common.RequestValidationError{Message: "limit: " + err.Error()}
	}
	lim := normalizeLimit(limit)
	products, err := s.repo.GetTopProductsByRating(lim)
	if err != nil {
		return TopProducts{}, fmt.Errorf("analytic service: get top products by rating: %w", err)
	}
	return products, nil
}
