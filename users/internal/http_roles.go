package internal

import (
	"context"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/madrabit/mini-market/users/internal/common"
	"go.uber.org/zap"
	"net/http"
	"time"
)

type ControllerRoles struct {
	svc    SvcRoles
	logger *common.Logger
}

func NewControllerRoles(svc SvcRoles, logger *common.Logger) *ControllerRoles {
	return &ControllerRoles{svc: svc, logger: logger}
}

type SvcRoles interface {
	CreateRole(ctx context.Context, req CreateRoleReq) error
	UpdateRole(ctx context.Context, id uuid.UUID, req UpdateRoleReq) error
	DeleteRole(ctx context.Context, id uuid.UUID) error
	GetUsersByRole(ctx context.Context, role string) (ListUsersResponse, error)
	GetRoleByName(ctx context.Context, name string) (Role, error)
}

func (c *ControllerRoles) Routes() chi.Router {
	r := chi.NewRouter()
	//Создание роли
	r.Post("/", c.CreateRole)
	//Обновление роли
	r.Patch("/{roleID}", c.UpdateRole)
	//Обновление роли
	r.Delete("/{roleID}", c.DeleteRole)
	// получить список пользователей по роли
	r.Get("/{role}", c.GetUsersByRole)
	return r
}

func (c *ControllerRoles) CreateRole(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), time.Second)
	defer cancel()
	defer func() {
		if err := r.Body.Close(); err != nil {
			c.logger.Error("failed to close request body", zap.Error(err))
		}
	}()
	var req CreateRoleReq
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		c.logger.Error("failed to decode create role request", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = c.svc.CreateRole(ctx, req)
	if err != nil {
		c.logger.Error("failed to create role", zap.Error(err))
		common.ErrResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func (c *ControllerRoles) UpdateRole(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), time.Second)
	defer cancel()
	roleId := chi.URLParam(r, "roleID")
	id, err := uuid.Parse(roleId)
	if err != nil || id == uuid.Nil {
		c.logger.Warn("invalid param")
		common.ErrResponse(w, http.StatusBadRequest, "invalid param")
		return
	}
	defer func() {
		if err := r.Body.Close(); err != nil {
			c.logger.Error("failed to close request body", zap.Error(err))
		}
	}()
	var req UpdateRoleReq
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		c.logger.Error("failed to decode update role request", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = c.svc.UpdateRole(ctx, id, req)
	if err != nil {
		c.logger.Error("failed to update user", zap.Error(err))
		common.ErrResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func (c *ControllerRoles) DeleteRole(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), time.Second)
	defer cancel()
	role := chi.URLParam(r, "roleID")
	id, err := uuid.Parse(role)
	if err != nil || id == uuid.Nil {
		c.logger.Warn("invalid param")
		common.ErrResponse(w, http.StatusBadRequest, "invalid param")
		return
	}
	err = c.svc.DeleteRole(ctx, id)
	if err != nil {
		c.logger.Error("failed to delete user", zap.Error(err))
		common.ErrResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func (c *ControllerRoles) GetUsersByRole(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), time.Second)
	defer cancel()
	role := chi.URLParam(r, "role")
	if role == "" {
		c.logger.Warn("invalid param")
		common.ErrResponse(w, http.StatusBadRequest, "invalid param")
		return
	}
	users, err := c.svc.GetUsersByRole(ctx, role)
	if err != nil {
		c.logger.Error("failed to retrieve users by role", zap.Error(err))
		common.ErrResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	w.Header().Set("Content-Type", "application/json")
	common.OkResponse(w, users)
}
