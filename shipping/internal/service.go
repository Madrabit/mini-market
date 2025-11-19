package internal

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/madrabit/mini-market/shipping/internal/common"
)

type Service struct {
	repo      Repo
	validator Validator
}

type Repo interface {
	CatchOrder(req CreateDeliveryRequest) (CreateDeliveryResponse, error)
	CheckStatus(userID, orderID uuid.UUID) (NotifyOrderResponse, error)
}

type Validator interface {
	Validate(request any) error
}

func NewService(repo Repo, validator Validator) *Service {
	return &Service{repo, validator}
}

func (s *Service) CatchOrder(req CreateDeliveryRequest) (CreateDeliveryResponse, error) {
	if err := s.validator.Validate(req); err != nil {
		return CreateDeliveryResponse{}, &common.RequestValidationError{Message: err.Error()}
	}
	resp, err := s.repo.CatchOrder(req)
	if err != nil {
		return CreateDeliveryResponse{}, fmt.Errorf("shipping service: failed to catch order: %w", err)
	}
	return resp, nil
}

func (s *Service) CheckStatus(userID, orderID uuid.UUID) (NotifyOrderResponse, error) {
	hint, err := s.repo.CheckStatus(userID, orderID)
	if err != nil {
		return NotifyOrderResponse{}, fmt.Errorf("shipping service: failed to check status: %w", err)
	}
	return hint, nil
}
