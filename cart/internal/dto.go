package internal

import "github.com/google/uuid"

type Cart struct {
	Id     uuid.UUID
	Userid uuid.UUID
	Items  map[uuid.UUID]Product
}

type Product struct {
	Id        uuid.UUID
	CartId    uuid.UUID
	ProductId uuid.UUID
	Qty       int64
}

type AddToCartRequest struct {
	ProductId uuid.UUID
	Qty       int64
}

type UpdateCartItemRequest struct {
	ProductId uuid.UUID
	Qty       int64
}

type CatalogRequest struct {
	ProductIds []uuid.UUID
}

type CatalogResponse struct {
	Products []CatalogProduct
}

type CatalogProduct struct {
	Id    uuid.UUID
	Name  string
	Price int64
}

type InventoryRequest struct {
	ItemsIds []uuid.UUID
}

type InventoryResponse struct {
	Statuses []InventoryStatus
}

type InventoryStatus struct {
	Product   uuid.UUID
	Available bool
	Quantity  int64
}
