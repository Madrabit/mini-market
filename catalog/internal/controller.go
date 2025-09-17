package internal

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/madrabit/mini-market/catalog/internal/common"
	"go.uber.org/zap"
	"net/http"
	"strconv"
)

/*
TODO
 -  Добавить в каталог
 - Обновить товар в каталоге
 - Удалить товар из каталога
 - Вернуть по массиву Id перечень товаров с именем и ценой
*/

type Controller struct {
	svc    Svc
	logger *common.Logger
}

func NewController(svc Svc, logger common.Logger) *Controller {
	return &Controller{svc: svc, logger: &logger}
}

type Svc interface {
	GetCatalog(limit int64, cursorID uuid.UUID) (GetCatalogResponse, error)
	AddProduct(item AddItemRequest) error
	UpdateProduct(item UpdateItemRequest) error
	DeleteProduct(id uuid.UUID) error
	GetProductById(id uuid.UUID) (Item, error)
}

func (c *Controller) Routes() chi.Router {
	r := chi.NewRouter()
	//Вернуть весь каталог
	r.Get("", c.GetCatalog)
	//Вернуть товар по id
	r.Get("/{productID}", c.GetProductById)
	//добавить в товар каталог
	r.Post("", c.AddProduct)
	//обновить товар в каталоге
	r.Patch("/{productID}", c.UpdateProduct)
	// удалить товар из каталога
	r.Delete("/{productID}", c.DeleteProduct)
	return r
}

func (c *Controller) GetCatalog(w http.ResponseWriter, r *http.Request) {
	limit := r.URL.Query().Get("limit")
	id := r.URL.Query().Get("cursorID")
	cursor, errCursor := uuid.Parse(id)
	lim, errLimit := strconv.Atoi(limit)
	if errLimit != nil || lim <= 0 || errCursor != nil || id == "" || cursor == uuid.Nil {
		c.logger.Warn("invalid param")
		common.ErrResponse(w, http.StatusBadRequest, "invalid param")
		return
	}
	catalog, err := c.svc.GetCatalog(int64(lim), cursor)
	if err != nil {
		c.logger.Error("failed to get catalog", zap.Error(err))
		common.ErrResponse(w, http.StatusBadRequest, error.Error(err))
		return
	}
	common.OkResponse(w, catalog)
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
		c.logger.Error("failed add to cart", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = c.svc.AddProduct(item)
	if err != nil {
		c.logger.Error("failed add to catalog", zap.Error(err))
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
		c.logger.Error("failed to delete item from cart", zap.Error(err))
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
