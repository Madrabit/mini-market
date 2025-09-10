package internal

import (
	"github.com/go-chi/chi/v5"
	"net/http"
)

/*
TODO


Автообновление эластика их каталогаа
Сбросить каталог и обновить весь
*/

type Controller struct {
}

func NewController() *Controller {
	return &Controller{}
}

func (c *Controller) Routes() chi.Router {
	r := chi.NewRouter()
	//Полнотекстовый поиск
	r.Get("/", c.FullTextSearching)
	//Подсказки выпадающие в поиске
	r.Get("/hint", c.DropDownHint)
	return r
}

func (c *Controller) FullTextSearching(writer http.ResponseWriter, request *http.Request) {

}

func (c *Controller) DropDownHint(writer http.ResponseWriter, request *http.Request) {

}
