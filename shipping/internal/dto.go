package internal

import (
	"github.com/google/uuid"
	"time"
)

type Status string

const (
	Pending    Status = "pending"    // Доставка создана, ждет сборки заказа
	Processing Status = "processing" // Заказ собран, передан в доставку
	Shipped    Status = "shipped"    // Курьер забрал заказ
	InTransit  Status = "in_transit" // В пути
	Delivered  Status = "delivered"  // Доставлен
	Cancelled  Status = "cancelled"  // Отменен
)

type Delivery struct {
	ID               uuid.UUID
	UserID           uint64
	Address          string
	InventoryAddress string
	Products         []Product
}

type Product struct {
	Items     map[uuid.UUID]int64
	OrderID   uuid.UUID
	UserID    uuid.UUID
	Status    Status
	CreatedAt time.Time
	UpdatedAt time.Time
}

// CreateDeliveryRequest Запрос от сервиса заказов
type CreateDeliveryRequest struct {
	Id      uuid.UUID
	OrderID uuid.UUID
	Address string
}

type CreateDeliveryResponse struct {
	Id      uuid.UUID
	OrderID uuid.UUID
	Status  Status
	Address string
}

type NotifyOrderRequest struct {
	OrderID uuid.UUID
	UserID  uint64
}

type NotifyOrderResponse struct {
	OrderIDs uuid.UUID
	Status   Status
}
