package internal

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/madrabit/mini-market/analytic/internal/common"
	"go.uber.org/zap"
	"log"
	"net/http"
)

type Controller struct {
	logger *common.Logger
}

func NewController() *Controller {
	return &Controller{}
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
	var event Event
	err := json.NewDecoder(r.Body).Decode(&event)
	if err != nil {

	}

}

func (c *Controller) getFileByProducts(w http.ResponseWriter, r *http.Request) {
	var request ProductsReq
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		c.logger.Error("failed to get employees", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	c.logger.Debug("get employees: received request", zap.Any("request", request))
	products, err := c.svc.FindByProducts(request.Products)
	if err != nil {
		c.logger.Error("failed to get employees", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	fileName := "emp.xls"
	w.Header().Set("Content-Disposition", "attachment; filename=\""+fileName+"\"")
	w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	w.Header().Set("Content-Length", fmt.Sprintf("%d", len(products)))
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(products); err != nil {
		c.logger.Error("failed to write response", zap.Error(err))
		log.Printf("failed to write response: %v", err)
	}
	c.logger.Info("successfully retrieve employees list")
}

func (c *Controller) GetOrdersSummary(w http.ResponseWriter, r *http.Request) {
	from := r.URL.Query().Get("from")
	to := r.URL.Query().Get("to")
}

func (c *Controller) GetTopProducts(w http.ResponseWriter, r *http.Request) {
	limit := r.URL.Query().Get("limit")
}

func (c *Controller) GetDailySales(w http.ResponseWriter, r *http.Request) {
	from := r.URL.Query().Get("from")
	to := r.URL.Query().Get("to")
}

func (c *Controller) GetSearchTrends(w http.ResponseWriter, r *http.Request) {
	limit := r.URL.Query().Get("limit")
}

func (c *Controller) GetFailSearch(w http.ResponseWriter, r *http.Request) {
	limit := r.URL.Query().Get("limit")
}

func (c *Controller) GetAvgRatingByProduct(w http.ResponseWriter, r *http.Request) {
	productID := r.URL.Query().Get("product_id")
}

func (c *Controller) GetTopByRating(w http.ResponseWriter, r *http.Request) {
	limit := r.URL.Query().Get("limit")
}
