package internal

import (
	"encoding/json"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/madrabit/mini-market/analytic/internal/common"
	"go.uber.org/zap"
	"net/http"
	"strconv"
)

type Controller struct {
	logger *common.Logger
	svc    Svc
}

type Svc interface {
	AddEvent(event Event) error
	GetOrdersSummary(from, to string) (OrdersSummary, error)
	GetTopProducts(limit int64) (TopProducts, error)
	GetDailySails(from, to string) (DailySales, error)
	GetSearchTrends(limit int64) (SearchTrends, error)
	GetFailSearch(limit int64) (FailedSearch, error)
	GetAvgRatingByProduct(productID uuid.UUID) (ProductAvgRating, error)
	GetTopProductsByRating(limit int64) (TopProducts, error)
}

func NewController(logger *common.Logger, svc Svc) *Controller {
	return &Controller{
		logger: logger,
		svc:    svc,
	}
}

func (c *Controller) Routes() chi.Router {
	r := chi.NewRouter()
	//Единый приемник событий со всех микросервисов
	r.Post("/orders/", c.HandelEvent)
	//GET /api/v1/analytics/orders/summary?from=2025-09-01&to=2025-09-07
	//Вернуть агрегаты: количество заказов, общую выручку, средний чек.
	r.Get("/orders/summary/", c.GetOrdersSummary)
	//GET /api/v1/analytics/orders/top-products?limit=10
	//Вернуть топ-товары по количеству продаж или по сумме.
	r.Get("/orders/top-products", c.GetTopProducts)
	//GET /api/v1/analytics/orders/daily?from=2025-09-01&to=2025-09-07
	//Продажи по дням (для графика).
	r.Get("/orders/daily", c.GetDailySales)
	// GET /api/v1/analytics/search/trending?limit=10
	//Популярные поисковые запросы за период.
	r.Get("/search/trending", c.GetSearchTrends)
	//GET /api/v1/analytics/search/failed?limit=10
	//Запросы без результатов (важно для каталога, чтобы видеть, что пользователи ищут).
	r.Get("/search/failed", c.GetFailSearch)
	//GET /api/v1/analytics/ratings/summary?product_id = xxx
	//Средняя оценка и распределение звёзд.
	r.Get("/ratings/summary", c.GetAvgRatingByProduct)
	//GET /api/v1/analytics/ratings/top?limit = 10
	// Вернуть топ продуктов по звездам
	r.Get("/ratings/top", c.GetTopByRating)
	return r
}

func (c *Controller) HandelEvent(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := r.Body.Close(); err != nil {
			c.logger.Error("error closing body", zap.Error(err))
		}
	}()
	var event Event
	err := json.NewDecoder(r.Body).Decode(&event)
	if err != nil {
		c.logger.Error("failed to get event", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = c.svc.AddEvent(event)
	if err != nil {
		c.logger.Error("failed to add event", zap.Error(err))
		common.ErrResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	common.OkResponse(w, event.ID)
}

func (c *Controller) GetOrdersSummary(w http.ResponseWriter, r *http.Request) {
	from := r.URL.Query().Get("from")
	to := r.URL.Query().Get("to")
	if from == "" || to == "" {
		c.logger.Warn("empty params")
		common.ErrResponse(w, http.StatusBadRequest, "empty param or more")
		return
	}
	orders, err := c.svc.GetOrdersSummary(from, to)
	var nfErr *common.NotFoundError
	if err != nil {
		if errors.As(err, &nfErr) {
			c.logger.Warn("info not found", zap.Error(err))
			common.OkResponseMsg(w, orders, "info not found")
			return
		}
		c.logger.Error("failed to get orders summary", zap.Error(err))
		common.ErrResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	common.OkResponse(w, orders)
}

func (c *Controller) GetTopProducts(w http.ResponseWriter, r *http.Request) {
	limit := r.URL.Query().Get("limit")
	lim, err := strconv.Atoi(limit)
	if err != nil && limit == "" && lim > 0 {
		c.logger.Warn("invalid param")
		common.ErrResponse(w, http.StatusBadRequest, "invalid param")
		return
	}
	top, err := c.svc.GetTopProducts(int64(lim))
	if err != nil {
		c.logger.Error("failed to get top products", zap.Error(err))
		common.ErrResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	common.OkResponse(w, top)
}

func (c *Controller) GetDailySales(w http.ResponseWriter, r *http.Request) {
	from := r.URL.Query().Get("from")
	to := r.URL.Query().Get("to")
	if from == "" || to == "" {
		c.logger.Warn("empty params")
		common.ErrResponse(w, http.StatusBadRequest, "empty param or more")
		return
	}
	sales, err := c.svc.GetDailySails(from, to)
	if err != nil {
		c.logger.Error("failed to get daily sales", zap.Error(err))
		common.ErrResponse(w, http.StatusBadRequest, error.Error(err))
		return
	}
	common.OkResponse(w, sales)
}

func (c *Controller) GetSearchTrends(w http.ResponseWriter, r *http.Request) {
	limit := r.URL.Query().Get("limit")
	lim, err := strconv.Atoi(limit)
	if err != nil && limit == "" && lim > 0 {
		c.logger.Warn("invalid param")
		common.ErrResponse(w, http.StatusBadRequest, "invalid param")
		return
	}
	searches, err := c.svc.GetSearchTrends(int64(lim))
	if err != nil {
		c.logger.Error("failed to get searching queries", zap.Error(err))
		common.ErrResponse(w, http.StatusBadRequest, error.Error(err))
		return
	}
	common.OkResponse(w, searches)

}

func (c *Controller) GetFailSearch(w http.ResponseWriter, r *http.Request) {
	limit := r.URL.Query().Get("limit")
	lim, err := strconv.Atoi(limit)
	if err != nil && limit == "" && lim > 0 {
		c.logger.Warn("invalid param")
		common.ErrResponse(w, http.StatusBadRequest, "invalid param")
		return
	}
	searches, err := c.svc.GetFailSearch(int64(lim))
	if err != nil {
		c.logger.Error("failed to get queries", zap.Error(err))
		common.ErrResponse(w, http.StatusBadRequest, error.Error(err))
		return
	}
	common.OkResponse(w, searches)
}

func (c *Controller) GetAvgRatingByProduct(w http.ResponseWriter, r *http.Request) {
	productID := r.URL.Query().Get("product_id")
	id, err := uuid.Parse(productID)
	if productID == "" {
		c.logger.Warn("invalid param")
		common.ErrResponse(w, http.StatusBadRequest, "invalid param")
		return
	}
	rating, err := c.svc.GetAvgRatingByProduct(id)
	if err != nil {
		c.logger.Error("failed to get avg rating by product", zap.Error(err))
		common.ErrResponse(w, http.StatusBadRequest, error.Error(err))
		return
	}
	common.OkResponse(w, rating)
}

func (c *Controller) GetTopByRating(w http.ResponseWriter, r *http.Request) {
	limit := r.URL.Query().Get("limit")
	lim, err := strconv.Atoi(limit)
	if err != nil && limit == "" && lim > 0 {
		c.logger.Warn("invalid param")
		common.ErrResponse(w, http.StatusBadRequest, "invalid param")
		return
	}
	top, err := c.svc.GetTopProductsByRating(int64(lim))
	if err != nil {
		c.logger.Error("failed to get top by rating", zap.Error(err))
		common.ErrResponse(w, http.StatusBadRequest, error.Error(err))
		return
	}
	common.OkResponse(w, top)
}
