package internal

import "github.com/google/uuid"

type Item struct {
	Id        uuid.UUID
	Name      string
	UnitPrice int64
}

type GetItemsRequest struct {
	ProductIds []uuid.UUID
}

type GetItemsResponse struct {
	Items []Item
}

type AddItemRequest struct {
	Name  string
	Price int64
}

type UpdateItemRequest struct {
	Id    uuid.UUID
	Name  string
	Price int64
}

type RemoveItemRequest struct {
	Id uuid.UUID
}
