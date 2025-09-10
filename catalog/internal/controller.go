package internal

import (
	"github.com/go-chi/chi/v5"
	"net/http"
)

/*
TODO
 -  Добавить в каталог
 - Обновить товар в каталоге
 - Удалить товар из каталога
 - Вернуть по массиву Id перечень товаров с именем и ценой
*/

type Controller struct {
}

func NewController() *Controller {
	return &Controller{}
}

func (c *Controller) Routes() chi.Router {
	r := chi.NewRouter()
	//Вернуть весь каталог
	r.Get("", c.GetCatalog)
	//Вернуть товар по id
	r.Get("/{productID}", c.GetProductById)
	//добавить в товар каталог
	r.Post("", c.AddProduct)
	//обновить товар в каталоге
	r.Patch("/{productID}", c.UpdateProduct)
	// удалить товар из каталога
	r.Delete("/{productID}", c.DeleteProduct)
	return r
}

func (c *Controller) GetCart(w http.ResponseWriter, r *http.Request) {

}
func (c *Controller) AddProduct(w http.ResponseWriter, r *http.Request) {

}
func (c *Controller) UpdateProduct(w http.ResponseWriter, r *http.Request) {

}
func (c *Controller) DeleteProduct(w http.ResponseWriter, r *http.Request) {

}

func (c *Controller) GetProductById(writer http.ResponseWriter, request *http.Request) {

}

func (c *Controller) GetCatalog(writer http.ResponseWriter, request *http.Request) {

}
