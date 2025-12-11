package internal

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/madrabit/mini-market/users/internal/common"
	"go.uber.org/zap"
	"net/http"
)

type ControllerUsers struct {
	svc    SvcUsers
	logger *common.Logger
}

func NewControllerUsers(svc SvcUsers, logger *common.Logger) *ControllerUsers {
	return &ControllerUsers{svc: svc, logger: logger}
}

type SvcUsers interface {
	CreateUser(req CreateUserReq) error
	UpdateUser(id uuid.UUID, req UpdateUserReq) error
	DeleteUser(DeleteUserReq uuid.UUID) error
	GetUserByID(userID uuid.UUID) (User, error)
	GetUsersByIds(req ListUsersRequest) (ListUsersResponse, error)
}

func (c *ControllerUsers) Routes() chi.Router {
	r := chi.NewRouter()
	//Создание пользователя
	r.Post("/", c.CreateUser)
	//Обновление пользователя
	r.Patch("/{userID}", c.UpdateUser)
	//Обновление пользователя
	r.Delete("/{userID}", c.DeleteUser)
	// получить список пользователей
	r.Get("/", c.GetUsersByIds)
	// получить одного пользователя
	r.Get("/{userID}", c.GetUserByID)
	return r
}

func (c *ControllerUsers) CreateUser(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := r.Body.Close(); err != nil {
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

func (c *ControllerUsers) UpdateUser(w http.ResponseWriter, r *http.Request) {
	user := chi.URLParam(r, "userID")
	userID, err := uuid.Parse(user)
	if err != nil || userID == uuid.Nil {
		c.logger.Warn("invalid param")
		common.ErrResponse(w, http.StatusBadRequest, "invalid param")
		return
	}
	defer func() {
		if err := r.Body.Close(); err != nil {
			c.logger.Error("failed to close request body", zap.Error(err))
		}
	}()
	var req UpdateUserReq
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		c.logger.Error("failed to decode update user request", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = c.svc.UpdateUser(userID, req)
	if err != nil {
		c.logger.Error("failed to update user", zap.Error(err))
		common.ErrResponse(w, http.StatusBadRequest, error.Error(err))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func (c *ControllerUsers) DeleteUser(w http.ResponseWriter, r *http.Request) {
	user := chi.URLParam(r, "userID")
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

func (c *ControllerUsers) GetUsersByIds(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := r.Body.Close(); err != nil {
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

func (c *ControllerUsers) GetUserByID(w http.ResponseWriter, r *http.Request) {
	user := chi.URLParam(r, "userID")
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
