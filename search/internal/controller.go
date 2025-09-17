package internal

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/madrabit/mini-market/search/internal/common"
	"go.uber.org/zap"
	"net/http"
	"strconv"
)

/*
TODO


Автообновление эластика их каталогаа
Сбросить каталог и обновить весь
*/

type Controller struct {
	svc    Svc
	logger *common.Logger
}

func NewController(svc Svc, logger common.Logger) *Controller {
	return &Controller{svc: svc, logger: &logger}
}

type Svc interface {
	FullTextSearching(req SearchRequest) (SearchResponse, error)
	DropDownHint(query string, limit int64) (SuggestResponse, error)
}

func (c *Controller) Routes() chi.Router {
	r := chi.NewRouter()
	//Полнотекстовый поиск
	r.Get("/", c.FullTextSearching)
	//Подсказки выпадающие в поиске
	r.Get("/hint", c.DropDownHint)
	return r
}

func (c *Controller) FullTextSearching(w http.ResponseWriter, r *http.Request) {
	defer func() {
		err := r.Body.Close()
		if err != nil {
			c.logger.Error("failed to search text", zap.Error(err))
		}
	}()
	var req SearchRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		c.logger.Error("failed to search text", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	resp, err := c.svc.FullTextSearching(req)
	if err != nil {
		c.logger.Error("failed to search text", zap.Error(err))
		common.ErrResponse(w, http.StatusBadRequest, error.Error(err))
		return
	}
	common.OkResponse(w, resp)
}

func (c *Controller) DropDownHint(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("query")
	lim := r.URL.Query().Get("limit")
	limit, err := strconv.Atoi(lim)
	if err != nil || query == "" || lim == "" || limit <= 0 {
		c.logger.Warn("invalid param")
		common.ErrResponse(w, http.StatusBadRequest, "invalid param")
		return
	}
	suggests, err := c.svc.DropDownHint(query, int64(limit))
	if err != nil {
		c.logger.Error("failed to search suggests", zap.Error(err))
		common.ErrResponse(w, http.StatusBadRequest, error.Error(err))
		return
	}
	common.OkResponse(w, suggests)
}
