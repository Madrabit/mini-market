package internal

import (
	"github.com/go-chi/chi/v5"
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
}

func NewController() *Controller {
	return &Controller{}
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

func (c *Controller) GetStatus(writer http.ResponseWriter, request *http.Request) {

}

func (c *Controller) UpdatePaymentStatus(writer http.ResponseWriter, request *http.Request) {

}

func (c *Controller) CreateOrder(writer http.ResponseWriter, request *http.Request) {

}
