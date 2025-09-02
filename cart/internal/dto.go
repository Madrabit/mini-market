package internal

import "github.com/google/uuid"

type Cart struct {
	Userid uuid.UUID
	Id     uuid.UUID
	Items  []Item
}

type Item struct {
	Id  uuid.UUID
	Qty int64
}

type CatalogResponse struct {
	ItemsIds []uuid.UUID
}

type CatalogRequest struct {
	Id    uuid.UUID
	Name  string
	Price int64
}

type InventoryResponse struct {
	ItemsIds []uuid.UUID
}

type InventoryRequest struct {
}
