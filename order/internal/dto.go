package internal

import (
	"github.com/google/uuid"
	"time"
)

type Status string

const (
	New       Status = "new"
	Shipped   Status = "shipped"
	Delivered Status = "delivered"
	Paid      Status = "paid"
	Canceled  Status = "canceled"
)

type ItemRow struct {
	ID        uuid.UUID `db:"id"`
	Name      string    `db:"name"`
	Quantity  int64     `db:"quantity"`
	OrderID   uuid.UUID `db:"order_id"`
	UnitPrice int64     `db:"unit_price"`
}

type OrderRow struct {
	ID         uuid.UUID `db:"id"`
	UserID     uuid.UUID `db:"user_id"`
	CreatedAt  time.Time `db:"created_at"`
	Status     Status    `db:"status"`
	GrandTotal int64     `db:"grand_total"`
}

type ItemQty struct {
	ID       uuid.UUID `json:"id" validate:"required"`
	Quantity int       `json:"quantity" validate:"gte=1"`
}

type CreatRequest struct {
	UserID uuid.UUID `json:"user_id" validate:"required"`
	Items  []ItemQty `json:"items" validate:"min=1,dive"`
}

type ItemResponse struct {
	ID        uuid.UUID `json:"id" validate:"required"`
	Name      string    `json:"name"`
	Quantity  int       `json:"quantity"`
	UnitPrice int64     `json:"unit_price"`
}

type OrderResponse struct {
	ID         uuid.UUID      `json:"id"`
	UserId     uuid.UUID      `json:"user_id"`
	Status     Status         `json:"status"`
	GrandTotal int64          `json:"grand_total"`
	Created    time.Time      `json:"created"`
	Items      []ItemResponse `json:"items"`
}
