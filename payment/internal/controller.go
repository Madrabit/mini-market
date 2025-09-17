package internal

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/madrabit/mini-market/payment/internal/common"
	"go.uber.org/zap"
	"net/http"
)

/*
 TODO
Действующие лица:

Order Service: Ядро, создает заказы.

Payment Service: Ваш сервис. Его задачи: принимать команду на оплату, взаимодействовать с внешним провайдером, обрабатывать его вебхуки и уведомлять Order Service о результате.

External Payment Provider: Ваша же заглушка, которая максимально упрощенно имитирует API и поведение реального провайдера (например, Tinkoff, Stripe).

У него есть своя "платежная форма" и он умеет отправлять вебхуки.

Последовательность действий (Полный цикл)
1. Инициирование платежа (Order Service -> Payment Service)

Order Service меняет статус заказа на "pending_payment".

Он вызывает POST /payments вашего Payment Service.

Тело запроса: { "order_id": "order_123", "amount": 10000, "currency": "RUB", "description": "Order #123" }

2. Взаимодействие с провайдером (Payment Service -> External Payment Provider)

Payment Service получает запрос, сохраняет его у себя в БД в статусе pending.

Он вызывает API внешнего провайдера (вашу заглушку), например, POST /v1/init (имитация Init у Tinkoff).

Тело запроса: { "Amount": 10000, "OrderId": "order_123", "SuccessURL": "...", "FailURL": "...", "CallbackURL": "https://payment-service.ru/api/webhooks/external" } <-- Ключевой момент! CallbackURL — это адрес Payment Service, куда провайдер будет слать вебхук.

3. Эмутация оплаты пользователем (Ручная)

Вы пишете простейшую HTML-страницу для External Payment Provider с кнопками "Оплатить успешно" и "Оплатить с ошибкой".

Или просто делаете POST-запрос на эндпоинт External Payment Provider для подтверждения платежа: POST /admin/confirm_payment { "OrderId": "order_123", "Status": "CONFIRMED" }.

4. Вебхук от провайдера (External Payment Provider -> Payment Service)

После "оплаты" External Payment Provider сам делает POST-запрос (вебхук) на тот CallbackURL, который ему передал Payment Service на шаге 2.

Тело вебхука (пример):

json
{
  "TerminalKey": "123",
  "OrderId": "order_123",
  "Success": true,
  "Status": "CONFIRMED",
  "PaymentId": "ext_pay_789",
  "Amount": 10000
}
5. Обработка вебхука (Payment Service)

Payment Service принимает вебхук. ВАЖНО: В реальном мире здесь проверяют подпись запроса, чтобы убедиться, что это зовет именно провайдер, а не злоумышленник.

Payment Service обновляет статус платежа в своей БД на succeeded.

Он уведомляет Order Service о результате.

6. Уведомление о результате (Payment Service -> Order Service)

Вот тут есть два подхода:

Синхронно (RPC): Payment Service сразу же сам вызывает API Order Service: PATCH /orders/order_123/status { "status": "paid" }.

Асинхронно (Events): Payment Service публикует событие payment.succeeded в брокер (Kafka, RabbitMQ), а Order Service на него подписан и обрабатывает.

Order Service получает уведомление и меняет статус заказа на paid.

Итог и почему это круто
Вы создаете не просто монолит, а распределенную систему из трех сервисов, которые общаются друг с другом через HTTP-вызовы и вебхуки.

Order Service вообще не знает о существовании External Payment Provider. Он знает только о вашем Payment Service.

Payment Service абстрагирует всю сложность работы с разными провайдерами. В будущем вы можете добавить еще один провайдер (например, "ЮKassa"), и Order Service даже не узнает об этом — изменения будут только в Payment Service.

External Payment Provider — это черный ящик, который работает ровно так же, как и в реальности.

Это идеальная практика для pet-проекта. Вы на собственном опыте поймете все сложности и паттерны (вебхуки, идемпотентность, retry logic) распределенных систем.
*/

type Controller struct {
	svc    Svc
	logger *common.Logger
}

func NewController(svc Svc, logger common.Logger) *Controller {
	return &Controller{svc: svc, logger: &logger}
}

type Svc interface {
	CreateOrder(req PaymentRequest) (CreatePaymentResponse, error)
	PSPWebhook(req PSPWebhookRequest) error
	GetStatus(userID, orderID uuid.UUID) (PaymentStatusResponse, error)
}

func (c *Controller) Routes() chi.Router {
	r := chi.NewRouter()
	//получить заказ от Order сервиса
	r.Post("/", c.CreateOrder)
	//получает от PSP по вебхуку что оплата прошла
	r.Post("/payment", c.PSPWebhook)
	// Получить статус оплаты
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
	var req PaymentRequest
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

func (c *Controller) PSPWebhook(w http.ResponseWriter, r *http.Request) {
	defer func() {
		err := r.Body.Close()
		if err != nil {
			c.logger.Error("failed to catch webhook", zap.Error(err))
		}
	}()
	var req PSPWebhookRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		c.logger.Error("failed to decode webhook", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = c.svc.PSPWebhook(req)
	if err != nil {
		c.logger.Error("failed to process webhook", zap.Error(err))
		common.ErrResponse(w, http.StatusBadRequest, error.Error(err))
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (c *Controller) GetStatus(w http.ResponseWriter, r *http.Request) {
	user := r.URL.Query().Get("userID")
	order := r.URL.Query().Get("orderID")
	userID, err := uuid.Parse(user)
	orderID, err := uuid.Parse(order)
	if err != nil || userID == uuid.Nil {
		c.logger.Warn("invalid param")
		common.ErrResponse(w, http.StatusBadRequest, "invalid param")
		return
	}
	status, err := c.svc.GetStatus(userID, orderID)
	if err != nil {
		c.logger.Error("failed to process webhook", zap.Error(err))
		common.ErrResponse(w, http.StatusBadRequest, error.Error(err))
		return
	}
	common.OkResponse(w, status)
}
