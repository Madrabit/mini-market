package internal

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/madrabit/mini-market/inventory/internal/common"
	"go.uber.org/zap"
	"net/http"
)

/*
TODO
 Запросить по Ids товары с количеством
 Добавить товар с количеством
 Изменить количество
 Удалить товар совсем
*/

type Controller struct {
	svc    Svc
	logger *common.Logger
}

func NewController(svc Svc, logger common.Logger) *Controller {
	return &Controller{svc: svc, logger: &logger}
}

type Svc interface {
	GetProductsByIds(IDs []uuid.UUID) (ListItemsResponse, error)
	AddProduct(item AddItemRequest) error
	UpdateProduct(item UpdateItemRequest) error
	DeleteProduct(id uuid.UUID) error
	GetProductById(id uuid.UUID) (Item, error)
	ReserveProducts(item ReserveItemRequest) error
	ReleaseProducts(item ReliesItemRequest) error
}

func (c *Controller) Routes() chi.Router {
	r := chi.NewRouter()
	//Запросить по Ids товары с количеством
	r.Post("/bulk-get", c.GetProductsByIDs)
	// GET /api/v1/inventory/{product_id} - Получить информацию о конкретном товаре
	r.Get("/{productID}", c.GetProductById)
	//Добавить товар с количеством
	r.Post("", c.AddProduct)
	//Изменить количество
	r.Patch("/{productID}", c.UpdateProduct)
	// Удалить товар совсем
	r.Delete("/{productID}", c.DeleteProduct)
	// POST /api/v1/inventory/reserve - Зарезервировать товары на время оформления заказа
	r.Post("/reserve", c.ReserveProducts)
	// POST /api/v1/inventory/release - Освободить резерв (если заказ отменен)
	r.Post("/release", c.ReleaseProducts)
	return r
}

func (c *Controller) GetProductsByIDs(w http.ResponseWriter, r *http.Request) {
	defer func() {
		err := r.Body.Close()
		if err != nil {
			c.logger.Error("failed to close body", zap.Error(err))
		}
	}()
	var req ListItemsRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		c.logger.Error("failed to get products by IDs", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	items, err := c.svc.GetProductsByIds(req.IDs)
	if err != nil {
		c.logger.Error("failed to get items by IDs", zap.Error(err))
		common.ErrResponse(w, http.StatusBadRequest, error.Error(err))
		return
	}
	common.OkResponse(w, items)
}

func (c *Controller) AddProduct(w http.ResponseWriter, r *http.Request) {
	defer func() {
		err := r.Body.Close()
		if err != nil {
			c.logger.Error("failed to close body", zap.Error(err))
		}
	}()
	var item AddItemRequest
	err := json.NewDecoder(r.Body).Decode(&item)
	if err != nil {
		c.logger.Error("failed add product to inventory", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = c.svc.AddProduct(item)
	if err != nil {
		c.logger.Error("failed add product to inventory", zap.Error(err))
		common.ErrResponse(w, http.StatusBadRequest, error.Error(err))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func (c *Controller) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	defer func() {
		err := r.Body.Close()
		if err != nil {
			c.logger.Error("failed to close body", zap.Error(err))
		}
	}()
	var item UpdateItemRequest
	err := json.NewDecoder(r.Body).Decode(&item)
	if err != nil {
		c.logger.Error("failed update product", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = c.svc.UpdateProduct(item)
	if err != nil {
		c.logger.Error("failed to update product", zap.Error(err))
		common.ErrResponse(w, http.StatusBadRequest, error.Error(err))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func (c *Controller) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	productID := r.URL.Query().Get("productID")
	id, err := uuid.Parse(productID)
	if productID == "" {
		c.logger.Warn("empty params")
		common.ErrResponse(w, http.StatusBadRequest, "empty param")
		return
	}
	err = c.svc.DeleteProduct(id)
	if err != nil {
		c.logger.Error("failed to delete item from inventory", zap.Error(err))
		common.ErrResponse(w, http.StatusBadRequest, error.Error(err))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func (c *Controller) GetProductById(w http.ResponseWriter, r *http.Request) {
	req := r.URL.Query().Get("productID")
	id, err := uuid.Parse(req)
	if err != nil || id == uuid.Nil {
		c.logger.Warn("invalid param")
		common.ErrResponse(w, http.StatusBadRequest, "invalid param")
		return
	}
	product, err := c.svc.GetProductById(id)
	if err != nil {
		c.logger.Error("failed to get product by ID", zap.Error(err))
		common.ErrResponse(w, http.StatusBadRequest, error.Error(err))
		return
	}
	common.OkResponse(w, product)
}

func (c *Controller) ReserveProducts(w http.ResponseWriter, r *http.Request) {
	defer func() {
		err := r.Body.Close()
		if err != nil {
			c.logger.Error("failed to close body", zap.Error(err))
		}
	}()
	var item ReserveItemRequest
	err := json.NewDecoder(r.Body).Decode(&item)
	if err != nil {
		c.logger.Error("failed reserve product", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = c.svc.ReserveProducts(item)
	if err != nil {
		c.logger.Error("failed reserve product", zap.Error(err))
		common.ErrResponse(w, http.StatusBadRequest, error.Error(err))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func (c *Controller) ReleaseProducts(w http.ResponseWriter, r *http.Request) {
	defer func() {
		err := r.Body.Close()
		if err != nil {
			c.logger.Error("failed to close body", zap.Error(err))
		}
	}()
	var item ReliesItemRequest
	err := json.NewDecoder(r.Body).Decode(&item)
	if err != nil {
		c.logger.Error("failed reserve product", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = c.svc.ReleaseProducts(item)
	if err != nil {
		c.logger.Error("failed reserve product", zap.Error(err))
		common.ErrResponse(w, http.StatusBadRequest, error.Error(err))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}
