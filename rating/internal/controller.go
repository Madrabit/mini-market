package internal

import (
	"github.com/go-chi/chi/v5"
	"net/http"
)

type Controller struct {
}

func NewController() *Controller {
	return &Controller{}
}

func (c *Controller) Routes() chi.Router {
	r := chi.NewRouter()
	//добавить отзыв
	r.Post("/", c.AddReview)
	//обновить отзыв
	r.Patch("/{reviewID}", c.UpdateReview)
	//удалить отзыв
	r.Delete("/{reviewID}", c.DeleteReview)
	//вернуть список отзывов на товар
	r.Get("/{productID}/reviews", c.GetReviewByProduct)
	return r
}

func (c *Controller) AddReview(writer http.ResponseWriter, request *http.Request) {

}

func (c *Controller) UpdateReview(writer http.ResponseWriter, request *http.Request) {

}

func (c *Controller) DeleteReview(writer http.ResponseWriter, request *http.Request) {

}

func (c *Controller) GetReviewByProduct(writer http.ResponseWriter, request *http.Request) {

}
