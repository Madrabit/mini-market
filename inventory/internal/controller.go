package internal

import (
	"github.com/go-chi/chi/v5"
	"net/http"
)

/*
TODO
 Запросить по Ids товары с количеством
 Добавить товар с количеством
 Изменить количество
 Удалить товар совсем
*/

type Controller struct {
}

func NewController() *Controller {
	return &Controller{}
}

func (c *Controller) Routes() chi.Router {
	r := chi.NewRouter()
	//Запросить по Ids товары с количеством
	r.Post("/bulk-get", c.GetProductsById)
	// GET /api/v1/inventory/{product_id} - Получить информацию о конкретном товаре
	r.Get("/{productID}", c.GetProduct)
	//Добавить товар с количеством
	r.Post("", c.AddProduct)
	//Изменить количество
	r.Patch("/{productID}", c.UpdateProduct)
	// Удалить товар совсем
	r.Delete("/{productID}", c.DeleteProduct)
	// POST /api/v1/inventory/reserve - Зарезервировать товары на время оформления заказа
	r.Post("/reserve", c.ReserveProducts)
	// POST /api/v1/inventory/release - Освободить резерв (если заказ отменен)
	r.Post("/release", c.ReleaseProducts)
	return r
}

func (c *Controller) GetProductsById(writer http.ResponseWriter, request *http.Request) {

}

func (c *Controller) AddProduct(writer http.ResponseWriter, request *http.Request) {

}

func (c *Controller) UpdateProduct(writer http.ResponseWriter, request *http.Request) {

}

func (c *Controller) DeleteProduct(writer http.ResponseWriter, request *http.Request) {

}

func (c *Controller) GetProduct(writer http.ResponseWriter, request *http.Request) {

}

func (c *Controller) ReserveProducts(writer http.ResponseWriter, request *http.Request) {

}

func (c *Controller) ReleaseProducts(writer http.ResponseWriter, request *http.Request) {

}
