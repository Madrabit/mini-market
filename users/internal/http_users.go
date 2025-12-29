package internal

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/madrabit/mini-market/users/internal/common"
	"go.uber.org/zap"
	"net/http"
	"strconv"
	"time"
)

type ControllerUsers struct {
	svc    SvcUsers
	logger *common.Logger
}

func NewControllerUsers(svc SvcUsers, logger *common.Logger) *ControllerUsers {
	return &ControllerUsers{svc: svc, logger: logger}
}

type SvcUsers interface {
	CreateUser(ctx context.Context, req CreateUserReq) error
	UpdateUser(ctx context.Context, id uuid.UUID, req UpdateUserReq) error
	DeleteUser(ctx context.Context, DeleteUserReq uuid.UUID) error
	GetUserByID(ctx context.Context, userID uuid.UUID) (User, error)
	GetUsersByIds(ctx context.Context, req ListUsersRequest) (ListUsersResponse, error)
	GetUsersPaginated(ctx context.Context, lastId uuid.UUID, pageSize int64) (ListUsersResponse, error)
}

func (c *ControllerUsers) Routes() chi.Router {
	r := chi.NewRouter()
	//Создание пользователя
	r.Post("/", c.CreateUser)
	//Обновление пользователя
	r.Patch("/{userID}", c.UpdateUser)
	//Обновление пользователя
	r.Delete("/{userID}", c.DeleteUser)
	// выводить всех пользователей по пагинации
	r.Get("/", c.GetUsersPaginated)
	// получить список пользователей по id
	r.Post("/search/", c.GetUsersByIds)
	// получить одного пользователя
	r.Get("/{userID}", c.GetUserByID)
	return r
}

func (c *ControllerUsers) CreateUser(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), time.Second)
	defer cancel()
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
	err = c.svc.CreateUser(ctx, req)
	if err != nil {
		c.logger.Error("failed to create user", zap.Error(err))
		common.ErrResponse(w, http.StatusBadRequest, error.Error(err))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func (c *ControllerUsers) UpdateUser(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), time.Second)
	defer cancel()
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
	err = c.svc.UpdateUser(ctx, userID, req)
	if err != nil {
		c.logger.Error("failed to update user", zap.Error(err))
		common.ErrResponse(w, http.StatusBadRequest, error.Error(err))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func (c *ControllerUsers) DeleteUser(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), time.Second)
	defer cancel()
	user := chi.URLParam(r, "userID")
	userID, err := uuid.Parse(user)
	if err != nil || userID == uuid.Nil {
		c.logger.Warn("invalid param")
		common.ErrResponse(w, http.StatusBadRequest, "invalid param")
		return
	}
	err = c.svc.DeleteUser(ctx, userID)
	if err != nil {
		c.logger.Error("failed to delete user", zap.Error(err))
		common.ErrResponse(w, http.StatusBadRequest, error.Error(err))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func (c *ControllerUsers) GetUsersByIds(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), time.Second)
	defer cancel()
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
	users, err := c.svc.GetUsersByIds(ctx, req)
	if err != nil {
		c.logger.Error("failed to create user", zap.Error(err))
		common.ErrResponse(w, http.StatusBadRequest, error.Error(err))
		return
	}
	common.OkResponse(w, users)

}

func (c *ControllerUsers) GetUsersPaginated(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), time.Second)
	defer cancel()
	var (
		lastId   uuid.UUID
		pageSize int64
		err      error
	)
	lastIdParam := r.URL.Query().Get("lastId")
	if lastIdParam != "" {
		lastId, err := uuid.Parse(lastIdParam)
		if err != nil || lastId == uuid.Nil {
			c.logger.Warn("invalid param")
			common.ErrResponse(w, http.StatusBadRequest, "invalid param")
			return
		}
	}
	pageSizeStr := r.URL.Query().Get("pageSize")
	if pageSizeStr != "" {
		pageSize, err = strconv.ParseInt(pageSizeStr, 10, 64)
		if err != nil {
			c.logger.Error("get page of users: wrong pageSize", zap.Error(err))
			common.ErrResponse(w, http.StatusBadRequest, err.Error())
			return
		}
	}
	users, err := c.svc.GetUsersPaginated(ctx, lastId, pageSize)
	var reqErr *common.RequestValidationError
	var notFoundErr *common.NotFoundError
	if err != nil {
		c.logger.Error("get page of users", zap.Error(err))
		switch {
		case errors.As(err, &notFoundErr):
			common.OkResponse(w, ListUsersResponse{})
			return
		case errors.As(err, &reqErr):
			common.ErrResponse(w, http.StatusBadRequest, err.Error())
			return
		default:
			common.ErrResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
	}
	common.OkResponse(w, users)
}

func (c *ControllerUsers) GetUserByID(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), time.Second)
	defer cancel()
	user := chi.URLParam(r, "userID")
	userID, err := uuid.Parse(user)
	if err != nil || userID == uuid.Nil {
		c.logger.Warn("invalid param")
		common.ErrResponse(w, http.StatusBadRequest, "invalid param")
		return
	}
	resp, err := c.svc.GetUserByID(ctx, userID)
	if err != nil {
		c.logger.Error("failed to get user by id", zap.Error(err))
		common.ErrResponse(w, http.StatusBadRequest, error.Error(err))
		return
	}
	common.OkResponse(w, resp)
}
