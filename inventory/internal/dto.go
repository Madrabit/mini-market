package internal

import "github.com/google/uuid"

type Warehouse struct {
	items map[uuid.UUID]Item
}

type Item struct {
	ID        uuid.UUID
	Qty       int64
	Reserved  int64
	Available int64
}

type ListItemsRequest struct {
	IDs []uuid.UUID
}

type ListItemsResponse struct {
	Items []Item
}

type AddItemRequest struct {
	Id  uuid.UUID
	Qty int64
}

type UpdateItemRequest struct {
	Id  uuid.UUID
	Qty int64
}

type DeleteItemRequest struct {
	Id uuid.UUID
}

type ReserveItemRequest struct {
	Id      uuid.UUID
	Qty     int64
	OrderID uuid.UUID // под какой заказ резерв
}

type ReliesItemRequest struct {
	Id  uuid.UUID
	Qty int64
}
