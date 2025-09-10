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
	//Создание пользователя
	r.Post("/", c.CreateUser)
	//Обновление пользователя
	r.Patch("/{userID}", c.UpdateUser)
	//Обновление пользователя
	r.Delete("/{userID}", c.DeleteUser)
	//Сменить роль пользователя
	r.Patch("{usersID}/roles", c.ChangeRole)
	// получить список пользователей
	r.Get("/", c.GetUsers)
	// получить одного пользователя
	r.Get("/{userID}", c.GetUserByID)
	return r
}

func (c *Controller) CreateUser(writer http.ResponseWriter, request *http.Request) {

}

func (c *Controller) UpdateUser(writer http.ResponseWriter, request *http.Request) {

}

func (c *Controller) DeleteUser(writer http.ResponseWriter, request *http.Request) {

}

func (c *Controller) ChangeRole(writer http.ResponseWriter, request *http.Request) {

}

func (c *Controller) GetUsers(writer http.ResponseWriter, request *http.Request) {

}

func (c *Controller) GetUserByID(writer http.ResponseWriter, request *http.Request) {

}
