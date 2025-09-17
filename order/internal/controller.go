package internal

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/madrabit/mini-market/order/internal/common"
	"go.uber.org/zap"
	"net/http"
)

/*
TODO
  Создание заказа:
  Из сервсиа Cart поступают данны заказа
Перед тем как перенаправить на оплату, он может (а часто и должен) сделать синхронные проверки:
Inventory Service: Запрос на резервирование товаров (POST /reservations). Склад временно "замораживает" необходимое количество. Если товара нет — сразу возвращается ошибка пользователю. Это лучше, чем брать деньги за несуществующий товар.
Catalog Service: Запрос актуальных цен и названий товаров (POST /prices). Это нужно для формирования финансовой суммы и снапшота.
  Пользак как то попадает в оплату. Оплачивает. И в заказ возвращается факт того что оплачено.
  Order → Inventory: из сервса склада списывается количество
  Order сохраняет заказ + строки (с снапшотом). Публикует order.created/order.paid.
  Payment отправляет от себя в Order что заказ оплачен
  Приходит уведомление из сервиса Notification, что заказ оплачен
  Delivery Service: Может подписаться на order.paid, чтобы начать процесс доставки (сборка, упаковка, логистика).
  Cart Service: Подписывается на order.paid и очищает корзину пользователя.
  Послать id заказа и получить статус, в каком он состоянии

*/

type Controller struct {
	svc    Svc
	logger *common.Logger
}

func NewController(svc Svc, logger common.Logger) *Controller {
	return &Controller{svc: svc, logger: &logger}
}

type Svc interface {
	CreateOrder(req CreatOrderRequest) (OrderResponse, error)
	GetStatus(user, order uuid.UUID) (StatusResponse, error)
	UpdatePaymentStatus(req UpdatePaymentStatusRequest) error
}

func (c *Controller) Routes() chi.Router {
	r := chi.NewRouter()
	//создать заказ
	r.Post("/", c.CreateOrder)
	//получает от сервиса payment что заказ оплачен
	r.Post("/{orderID}/payment-status", c.UpdatePaymentStatus)
	// Получить статус заказа
	r.Get("{/orderID}", c.GetStatus)
	return r
}

func (c *Controller) CreateOrder(w http.ResponseWriter, r *http.Request) {
	defer func() {
		err := r.Body.Close()
		if err != nil {
			c.logger.Error("failed to create order", zap.Error(err))
		}
	}()
	var req CreatOrderRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		c.logger.Error("failed to create order", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	resp, err := c.svc.CreateOrder(req)
	if err != nil {
		c.logger.Error("failed to create order", zap.Error(err))
		common.ErrResponse(w, http.StatusBadRequest, error.Error(err))
		return
	}
	common.OkResponse(w, resp)
}

func (c *Controller) UpdatePaymentStatus(w http.ResponseWriter, r *http.Request) {
	defer func() {
		err := r.Body.Close()
		if err != nil {
			c.logger.Error("failed to update payment status", zap.Error(err))
		}
	}()
	var req UpdatePaymentStatusRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		c.logger.Error("failed to update payment status", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = c.svc.UpdatePaymentStatus(req)
	if err != nil {
		c.logger.Error("failed to update payment status", zap.Error(err))
		common.ErrResponse(w, http.StatusBadRequest, error.Error(err))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func (c *Controller) GetStatus(w http.ResponseWriter, r *http.Request) {
	orderId := r.URL.Query().Get("orderID")
	userID := r.URL.Query().Get("userID")
	user, err := uuid.Parse(userID)
	order, err := uuid.Parse(orderId)
	if err != nil || user == uuid.Nil || order == uuid.Nil {
		c.logger.Warn("invalid param")
		common.ErrResponse(w, http.StatusBadRequest, "invalid param")
		return
	}
	status, err := c.svc.GetStatus(user, order)
	if err != nil {
		c.logger.Error("failed to get order status", zap.Error(err))
		common.ErrResponse(w, http.StatusBadRequest, error.Error(err))
		return
	}
	common.OkResponse(w, status)
}
