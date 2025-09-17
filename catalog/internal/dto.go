package internal

import "github.com/google/uuid"

type Item struct {
	Id        uuid.UUID
	Name      string
	UnitPrice int64
}

type GetCatalogRequest struct {
	Limit    int64
	CursorID uuid.UUID
}

type GetCatalogResponse struct {
	Items        []Item
	NextCursorID uuid.UUID
	HasMore      bool
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
