package internal

import (
	"github.com/go-chi/chi/v5"
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
}

func NewController() *Controller {
	return &Controller{}
}

func (c *Controller) Routes() chi.Router {
	r := chi.NewRouter()
	//От Order прилетает заказ. Доставка вносит заказ и ставит статус pending
	r.Post("/orders", c.CatchOrder)
	//Сервис заказов может обратиться к Shipping и получить, что заказ доставлен
	r.Get("/orders/{orderID}/orders-status", c.CheckStatus)
	return r
}

func (c *Controller) CatchOrder(writer http.ResponseWriter, request *http.Request) {

}

func (c *Controller) CheckStatus(writer http.ResponseWriter, request *http.Request) {

}
