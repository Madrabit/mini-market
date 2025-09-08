package internal

import "github.com/google/uuid"

type Warehouse struct {
	items map[uuid.UUID]Item
}

type Item struct {
	Qty int64
}

type ListItemsRequest struct {
	Id uuid.UUID
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
