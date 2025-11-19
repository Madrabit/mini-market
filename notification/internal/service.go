package internal

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/madrabit/mini-market/notification/internal/common"
	"time"
)

type Service struct {
	repo      Repo
	validator Validator
	sender    Sender
}

type Repo interface {
	BeginTransaction() (tx *sqlx.Tx, err error)
	FindItemById(tx *sqlx.Tx, productID uuid.UUID) (bool, error)
	Notify(req NotificationRequest) (NotificationResponse, error)
	NotifyTest(req NotificationRequest) (NotificationResponse, error)
	GetUserNotifications(id uuid.UUID) ([]NotificationResponse, error)
}

type Validator interface {
	Validate(request any) error
}

type Sender interface {
	Send(to string, message string) error
}

func NewService(repo Repo, validator Validator, sender Sender) *Service {
	return &Service{repo, validator, sender}
}

func (s *Service) Notify(req NotificationRequest) (NotificationResponse, error) {
	if err := s.validator.Validate(req); err != nil {
		return NotificationResponse{}, &common.RequestValidationError{Message: err.Error()}
	}
	err := s.sender.Send(req.To, req.Text)
	if err != nil {
		return NotificationResponse{}, fmt.Errorf("notification service: failed to send notification: %w", err)
	}
	response := NotificationResponse{
		Success:   true,
		MessageID: "0",
		Status:    "sent",
		Timestamp: time.Now(),
	}
	return response, nil
}

func (s *Service) NotifyTest(to string) (NotificationResponse, error) {
	err := s.sender.Send(to, "test notification")
	if err != nil {
		return NotificationResponse{}, fmt.Errorf("notification service: failed to send notification: %w", err)
	}
	response := NotificationResponse{
		Success:   true,
		MessageID: "0",
		Status:    "sent",
		Timestamp: time.Now(),
	}
	return response, nil
}

func (s *Service) GetUserNotifications(userID uuid.UUID) ([]NotificationResponse, error) {
	if userID == uuid.Nil {
		return []NotificationResponse{}, errors.New("notification service: invalid id")
	}
	product, err := s.repo.GetUserNotifications(userID)
	if err != nil {
		return []NotificationResponse{}, fmt.Errorf("notification service: failed to get notifucations by userID: %w", err)
	}
	return product, nil
}
