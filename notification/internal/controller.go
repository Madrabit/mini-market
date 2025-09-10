package internal

import (
	"github.com/go-chi/chi/v5"
	"net/http"
)

/*
POST /internal/v1/notify
→ отправить конкретное сообщение пользователю.
POST /notify/test
*/

type Controller struct {
}

func NewController() *Controller {
	return &Controller{}
}

func (c *Controller) Routes() chi.Router {
	r := chi.NewRouter()
	//отправить конкретное сообщение пользователю
	r.Post("/notify", c.Notify)
	//отправить конкретное сообщение пользователю
	r.Post("/notify/test", c.NotifyTest)
	//История уведомлений
	r.Get("/notifications/{userID}", c.GetUserNotifications)
	return r
}

func (c *Controller) Notify(writer http.ResponseWriter, request *http.Request) {

}

func (c *Controller) NotifyTest(writer http.ResponseWriter, request *http.Request) {

}

func (c *Controller) GetUserNotifications(writer http.ResponseWriter, request *http.Request) {

}
