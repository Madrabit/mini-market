package internal

import (
	"github.com/google/uuid"
	"time"
)

type Status string

const (
	Captured   Status = "captured"   // спинсаны
	Authorized Status = "authorized" // деньги заблокированы
	Pending    Status = "pending"
	Rejected   Status = "rejected"
	Failed     Status = "failed"
	Canceled   Status = "canceled"
)

type Payment struct {
	ID         uuid.UUID
	UserID     uuid.UUID
	OrderID    uuid.UUID
	Amount     int64
	Currency   string
	Status     Status
	ExternalID string // id транзакции в PSP (если есть)
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type PaymentRequest struct {
	UserID   uuid.UUID
	OrderID  string `json:"orderId"` // Важно: uuid.UUID -> string провайдеры ожидают строку
	Amount   int64
	Currency string
}

type CreatePaymentResponse struct {
	PaymentID uuid.UUID // ID платежа в вашей системе
	Status    Status    // "pending", "requires_action"
	Amount    int64     // Сумма к оплате (может отличаться от заказа)
	Currency  string    // Валюта оплаты
}

// { "Amount": 10000, "OrderId": "order_123", "SuccessURL": "...", "FailURL": "...", "CallbackURL": "https://payment-service.ru/api/webhooks/external" }
// Ключевой момент! CallbackURL — это адрес Payment Service, куда провайдер будет слать вебхук.
type PSPRequest struct {
	Amount        int64
	OrderID       uuid.UUID
	CallbackURL   string
	SuccessURL    string // https://my-shop.ru/order/success?order_id=123
	FailURL       string // https://my-shop.ru/order/faild?order_id=123
	Description   string // "Оплата заказа №123 в интернет-магазине 'MyShop'. Товары: Кроссовки Nike, Футболка Adidas."
	CustomerEmail string
}

type PaymentStatusResponse struct {
	OrderID   uuid.UUID
	PaymentID uuid.UUID
	Status    Status
	Amount    int64
	Currency  string
}

type PSPWebhookRequest struct {
	EventID   string    `json:"eventId"`   // ID события у провайдера
	PaymentID string    `json:"paymentId"` // ExternalID (их ID платежа) это айди вебхука, чтобы его не дублировать например
	OrderID   string    `json:"orderId"`
	Status    Status    `json:"status"` // Но часто статус у них свой, строка: "succeeded", "failed"
	Amount    int64     `json:"amount"`
	Currency  string    `json:"currency"`
	Signature string    `json:"signature"` // Цифровая подпись для проверки валидности вебхука
	CreatedAt time.Time `json:"created_at"`
}
