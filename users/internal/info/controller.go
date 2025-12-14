package info

import (
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
	"github.com/madrabit/mini-market/users/internal/common"
	"net/http"
	"time"
)

type Controller struct {
	logger *common.Logger
	cfg    *common.Config
	db     *sqlx.DB
}

func NewInfoController(logger *common.Logger, cfg *common.Config, db *sqlx.DB) *Controller {
	return &Controller{
		logger: logger,
		cfg:    cfg,
		db:     db,
	}
}

func (c *Controller) Routes() chi.Router {
	r := chi.NewRouter()
	r.Get("/info", c.GetInfo)
	r.Get("/health", c.GetHealth)
	r.Get("/ready", c.GetReady)
	return r
}

type Response struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

func (c *Controller) GetInfo(w http.ResponseWriter, r *http.Request) {
	resp := Response{
		c.cfg.AppName,
		c.cfg.AppVersion,
	}
	common.OkResponse(w, resp)
	return
}

func (c *Controller) GetHealth(w http.ResponseWriter, r *http.Request) {
	common.OkResponse(w, map[string]string{"status": "alive"})
	return
}

func (c *Controller) GetReady(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), time.Second)
	defer cancel()
	if err := c.db.PingContext(ctx); err != nil {
		common.ErrResponse(w, http.StatusServiceUnavailable, "db not ready")
		return
	}
	common.OkResponse(w, map[string]string{"status": "ready"})
	return
}
