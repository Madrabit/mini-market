package internal

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/madrabit/mini-market/shipping/internal/common"
	"go.uber.org/zap"
	"net/http"
)

/*
TODO

От Order прилетает заказ.
Доставка вносит заказ и ставит статус pending
Потом по таймеру типа shipping происходит
Сервис заказов может обратиться к Shipping и получить, что заказ доставлен
Либо когда доставлено, отправляется уведомление через notify
*/

type Controller struct {
	svc    Svc
	logger *common.Logger
}

func NewController(svc Svc, logger common.Logger) *Controller {
	return &Controller{svc: svc, logger: &logger}
}

type Svc interface {
	CatchOrder(req CreateDeliveryRequest) (CreateDeliveryResponse, error)
	CheckStatus(userID, orderID uuid.UUID) (NotifyOrderResponse, error)
}

func (c *Controller) Routes() chi.Router {
	r := chi.NewRouter()
	//От Order прилетает заказ. Доставка вносит заказ и ставит статус pending
	r.Post("/orders", c.CatchOrder)
	//Сервис заказов может обратиться к Shipping и получить, что заказ доставлен
	r.Get("/orders/{orderID}/orders-status", c.CheckStatus)
	return r
}

func (c *Controller) CatchOrder(w http.ResponseWriter, r *http.Request) {
	defer func() {
		err := r.Body.Close()
		if err != nil {
			c.logger.Error("failed to catch webhook", zap.Error(err))
		}
	}()
	var req CreateDeliveryRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		c.logger.Error("failed to decode order", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	resp, err := c.svc.CatchOrder(req)
	if err != nil {
		c.logger.Error("failed to create order", zap.Error(err))
		common.ErrResponse(w, http.StatusBadRequest, error.Error(err))
		return
	}
	common.OkResponse(w, resp)
}

func (c *Controller) CheckStatus(w http.ResponseWriter, r *http.Request) {
	user := r.URL.Query().Get("userID")
	order := r.URL.Query().Get("orderID")
	userID, err := uuid.Parse(user)
	orderID, err := uuid.Parse(order)
	if err != nil || userID == uuid.Nil {
		c.logger.Warn("invalid param")
		common.ErrResponse(w, http.StatusBadRequest, "invalid param")
		return
	}
	status, err := c.svc.CheckStatus(userID, orderID)
	if err != nil {
		c.logger.Error("failed to check status", zap.Error(err))
		common.ErrResponse(w, http.StatusBadRequest, error.Error(err))
		return
	}
	common.OkResponse(w, status)
}
