package internal

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/madrabit/mini-market/rating/internal/common"
	"go.uber.org/zap"
	"net/http"
)

type Controller struct {
	svc    Svc
	logger *common.Logger
}

func NewController(svc Svc, logger common.Logger) *Controller {
	return &Controller{svc: svc, logger: &logger}
}

type Svc interface {
	AddReview(item AddReviewRequest) error
	UpdateReview(item UpdateReviewRequest) error
	DeleteReview(user, order uuid.UUID) error
	GetReviewsByProduct(productID uuid.UUID) (ProductReviewsResponse, error)
}

func (c *Controller) Routes() chi.Router {
	r := chi.NewRouter()
	//добавить отзыв
	r.Post("/", c.AddReview)
	//обновить отзыв
	r.Patch("/{reviewID}", c.UpdateReview)
	//удалить отзыв
	r.Delete("/{reviewID}", c.DeleteReview)
	//вернуть список отзывов на товар
	r.Get("/{productID}/reviews", c.GetReviewsByProduct)
	return r
}

func (c *Controller) AddReview(w http.ResponseWriter, r *http.Request) {
	defer func() {
		err := r.Body.Close()
		if err != nil {
			c.logger.Error("failed to close body", zap.Error(err))
		}
	}()
	var item AddReviewRequest
	err := json.NewDecoder(r.Body).Decode(&item)
	if err != nil {
		c.logger.Error("failed to add review", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = c.svc.AddReview(item)
	if err != nil {
		c.logger.Error("failed to add review", zap.Error(err))
		common.ErrResponse(w, http.StatusBadRequest, error.Error(err))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func (c *Controller) UpdateReview(w http.ResponseWriter, r *http.Request) {
	defer func() {
		err := r.Body.Close()
		if err != nil {
			c.logger.Error("failed to close body", zap.Error(err))
		}
	}()
	var item UpdateReviewRequest
	err := json.NewDecoder(r.Body).Decode(&item)
	if err != nil {
		c.logger.Error("failed to update review", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = c.svc.UpdateReview(item)
	if err != nil {
		c.logger.Error("failed to update review", zap.Error(err))
		common.ErrResponse(w, http.StatusBadRequest, error.Error(err))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func (c *Controller) DeleteReview(w http.ResponseWriter, r *http.Request) {
	orderId := r.URL.Query().Get("reviewID")
	userID := r.URL.Query().Get("userID")
	user, err := uuid.Parse(userID)
	order, err := uuid.Parse(orderId)
	if err != nil || user == uuid.Nil || order == uuid.Nil {
		c.logger.Warn("invalid param")
		common.ErrResponse(w, http.StatusBadRequest, "invalid param")
		return
	}
	err = c.svc.DeleteReview(user, order)
	if err != nil {
		c.logger.Error("failed to delete review", zap.Error(err))
		common.ErrResponse(w, http.StatusBadRequest, error.Error(err))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func (c *Controller) GetReviewsByProduct(w http.ResponseWriter, r *http.Request) {
	productID := r.URL.Query().Get("productID")
	id, err := uuid.Parse(productID)
	if err != nil || id == uuid.Nil {
		c.logger.Warn("invalid param")
		common.ErrResponse(w, http.StatusBadRequest, "invalid param")
		return
	}
	reviews, err := c.svc.GetReviewsByProduct(id)
	if err != nil {
		c.logger.Error("failed to get order status", zap.Error(err))
		common.ErrResponse(w, http.StatusBadRequest, error.Error(err))
		return
	}
	common.OkResponse(w, reviews)
}
