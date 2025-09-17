package internal

import (
	"encoding/json"
	"github.com/google/uuid"
	"time"
)

type Event struct {
	ID         uuid.UUID
	Type       string
	Payload    json.RawMessage
	OccurredAt time.Time
}

type OrdersSummary struct {
	OrdersSum    int64
	TotalRevenue int64
	AvgBill      int64
}

type TopProducts struct {
	Products []Product
}

type Product struct {
	ID    uuid.UUID
	Name  string
	Price string
}

type DailySales struct {
	sales []Sale
}

type Sale struct {
	Date        string `json:"date"`
	TotalAmount int64  `json:"total_amount"`
	OrdersCount int64  `json:"orders_count,omitempty"`
}

type SearchTrends struct {
	Queries []string
}

type FailedSearch struct {
	Queries []string
}

type ProductAvgRating struct {
	Name string
	Avg  float64
}

type TopByRating struct {
	Top []Product
}
