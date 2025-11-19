package internal

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/madrabit/mini-market/notification/internal/common"
	"go.uber.org/zap"
	"net/http"
)

/*
POST /internal/v1/notify
→ отправить конкретное сообщение пользователю.
POST /notify/test
*/

type Controller struct {
	svc    Svc
	logger *common.Logger
}

func NewController(svc Svc, logger common.Logger) *Controller {
	return &Controller{svc: svc, logger: &logger}
}

type Svc interface {
	Notify(req NotificationRequest) (NotificationResponse, error)
	NotifyTest(to string) (NotificationResponse, error)
	GetUserNotifications(id uuid.UUID) ([]NotificationResponse, error)
}

func (c *Controller) Routes() chi.Router {
	r := chi.NewRouter()
	//отправить конкретное сообщение пользователю
	r.Post("/notify", c.Notify)
	//отправить конкретное сообщение пользователю
	r.Post("/notify/test/{to}", c.NotifyTest)
	//История уведомлений
	r.Get("/notifications/{userID}", c.GetUserNotifications)
	return r
}

func (c *Controller) Notify(w http.ResponseWriter, r *http.Request) {
	defer func() {
		err := r.Body.Close()
		if err != nil {
			c.logger.Error("failed to notify", zap.Error(err))
		}
	}()
	var req NotificationRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		c.logger.Error("failed to notify", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	resp, err := c.svc.Notify(req)
	if err != nil {
		c.logger.Error("failed to notify", zap.Error(err))
		common.ErrResponse(w, http.StatusBadRequest, error.Error(err))
		return
	}
	common.OkResponse(w, resp)
}

func (c *Controller) NotifyTest(w http.ResponseWriter, r *http.Request) {
	to := r.URL.Query().Get("to")
	if to == "" {
		c.logger.Warn("invalid param")
		common.ErrResponse(w, http.StatusBadRequest, "invalid param")
		return
	}
	resp, err := c.svc.NotifyTest(to)
	if err != nil {
		c.logger.Error("failed to notify", zap.Error(err))
		common.ErrResponse(w, http.StatusBadRequest, error.Error(err))
		return
	}
	common.OkResponse(w, resp)
}

func (c *Controller) GetUserNotifications(w http.ResponseWriter, r *http.Request) {
	req := r.URL.Query().Get("userID")
	id, err := uuid.Parse(req)
	if err != nil || id == uuid.Nil {
		c.logger.Warn("invalid param")
		common.ErrResponse(w, http.StatusBadRequest, "invalid param")
		return
	}
	notifies, err := c.svc.GetUserNotifications(id)
	if err != nil {
		c.logger.Error("failed to get product by ID", zap.Error(err))
		common.ErrResponse(w, http.StatusBadRequest, error.Error(err))
		return
	}
	common.OkResponse(w, notifies)
}
