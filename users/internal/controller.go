package internal

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/madrabit/mini-market/users/internal/common"
	"go.uber.org/zap"
	"net/http"
)

type Controller struct {
	svc    Svc
	logger *common.Logger
}

func NewController(svc Svc, logger common.Logger) *Controller {
	return &Controller{svc: svc, logger: &logger}
}

type Svc interface {
	CreateUser(req CreateUserReq) error
	UpdateUser(req UpdateUserReq) error
	DeleteUser(DeleteUserReq uuid.UUID) error
	ChangeRole(req SetUserRoleReq) error
	GetUserByID(userID uuid.UUID) (User, error)
	GetUsersByIds(IDs ListUsersRequest) (ListUsersResponse, error)
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
	r.Get("/", c.GetUsersByIds)
	// получить одного пользователя
	r.Get("/{userID}", c.GetUserByID)
	return r
}

func (c *Controller) CreateUser(w http.ResponseWriter, r *http.Request) {
	defer func() {
		err := r.Body.Close()
		if err != nil {
			c.logger.Error("failed to close request body", zap.Error(err))
		}
	}()
	var req CreateUserReq
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		c.logger.Error("failed to decode create user request", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = c.svc.CreateUser(req)
	if err != nil {
		c.logger.Error("failed to create user", zap.Error(err))
		common.ErrResponse(w, http.StatusBadRequest, error.Error(err))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func (c *Controller) UpdateUser(w http.ResponseWriter, r *http.Request) {
	defer func() {
		err := r.Body.Close()
		if err != nil {
			c.logger.Error("failed to close request body", zap.Error(err))
		}
	}()
	var req UpdateUserReq
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		c.logger.Error("failed to decode update user request", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = c.svc.UpdateUser(req)
	if err != nil {
		c.logger.Error("failed to update user", zap.Error(err))
		common.ErrResponse(w, http.StatusBadRequest, error.Error(err))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func (c *Controller) DeleteUser(w http.ResponseWriter, r *http.Request) {
	user := r.URL.Query().Get("userID")
	userID, err := uuid.Parse(user)
	if err != nil || userID == uuid.Nil {
		c.logger.Warn("invalid param")
		common.ErrResponse(w, http.StatusBadRequest, "invalid param")
		return
	}
	err = c.svc.DeleteUser(userID)
	if err != nil {
		c.logger.Error("failed to delete user", zap.Error(err))
		common.ErrResponse(w, http.StatusBadRequest, error.Error(err))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func (c *Controller) ChangeRole(w http.ResponseWriter, r *http.Request) {
	defer func() {
		err := r.Body.Close()
		if err != nil {
			c.logger.Error("failed to close request body", zap.Error(err))
		}
	}()
	var req SetUserRoleReq
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		c.logger.Error("failed to decode user role changer", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = c.svc.ChangeRole(req)
	if err != nil {
		c.logger.Error("failed to change user role", zap.Error(err))
		common.ErrResponse(w, http.StatusBadRequest, error.Error(err))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func (c *Controller) GetUsersByIds(w http.ResponseWriter, r *http.Request) {
	defer func() {
		err := r.Body.Close()
		if err != nil {
			c.logger.Error("failed to close request body", zap.Error(err))
		}
	}()
	var req ListUsersRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		c.logger.Error("failed to decode user role changer", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	users, err := c.svc.GetUsersByIds(req)
	if err != nil {
		c.logger.Error("failed to create user", zap.Error(err))
		common.ErrResponse(w, http.StatusBadRequest, error.Error(err))
		return
	}
	common.OkResponse(w, users)

}

func (c *Controller) GetUserByID(w http.ResponseWriter, r *http.Request) {
	user := r.URL.Query().Get("userID")
	userID, err := uuid.Parse(user)
	if err != nil || userID == uuid.Nil {
		c.logger.Warn("invalid param")
		common.ErrResponse(w, http.StatusBadRequest, "invalid param")
		return
	}
	resp, err := c.svc.GetUserByID(userID)
	if err != nil {
		c.logger.Error("failed to get user by id", zap.Error(err))
		common.ErrResponse(w, http.StatusBadRequest, error.Error(err))
		return
	}
	common.OkResponse(w, resp)
}
