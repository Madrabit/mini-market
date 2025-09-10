package internal

import (
	"github.com/go-chi/chi/v5"
	"net/http"
)

/*
TODO
 - добавить в корзину
      запросить цену из Catalog
      запросить количество из Inventory и наличие
 - обновить в корзине кол-во товара
       обновить в инвентаре
 - удалить из корзины товар
       обновить в инвентаре
*/

type Controller struct {
}

func NewController() *Controller {
	return &Controller{}
}

func (c *Controller) Routes() chi.Router {
	r := chi.NewRouter()
	//Вернуть корзину
	r.Get("/", c.GetCart)
	//добавить в корзину
	r.Post("/items", c.AddToCart)
	//обновить товар корзине
	r.Patch("/items/{productID}", c.UpdateCart)
	// удалить товар из корзины
	r.Delete("/items/{productID}", c.DeleteProduct)
	return r
}

func (c *Controller) GetCart(w http.ResponseWriter, r *http.Request) {

}

func (c *Controller) AddToCart(w http.ResponseWriter, r *http.Request) {
	var items AddToCartRequest
}

func (c *Controller) UpdateCart(w http.ResponseWriter, r *http.Request) {

	var item UpdateCartItemRequest
}

func (c *Controller) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	var id RemoveCartItemRequest
}
