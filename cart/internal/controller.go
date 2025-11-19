package internal

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/madrabit/mini-market/cart/internal/common"
	"go.uber.org/zap"
	"net/http"
)

/*
TODO
 - добавить в корзину
      запросить цену из Catalog
      запросить количество из Inventory и наличие
 - обновить в корзине кол-во товара
       обновить в инвентаре
 - удалить из корзины товар
       обновить в инвентаре
*/

type Controller struct {
	svc    Svc
	logger *common.Logger
}

func NewController(svc Svc, logger common.Logger) *Controller {
	return &Controller{svc: svc, logger: &logger}
}

type Svc interface {
	GetCart(userID uuid.UUID) (Cart, error)
	AddToCart(item AddToCartRequest) error
	UpdateCart(item UpdateCartItemRequest) error
	DeleteProduct(id uuid.UUID) error
}

func (c *Controller) Routes() chi.Router {
	r := chi.NewRouter()
	//Вернуть корзину
	r.Get("/", c.GetCart)
	//добавить в корзину
	r.Post("/items", c.AddToCart)
	//обновить товар корзине
	r.Patch("/items/{productID}", c.UpdateCart)
	// удалить товар из корзины
	r.Delete("/items/{productID}", c.DeleteProduct)
	return r
}

func (c *Controller) GetCart(w http.ResponseWriter, r *http.Request) {
	mockID, err := uuid.Parse("1b98a34a-cbcf-4e24-a4b8-2a218f5b68fc")
	if err != nil {
		c.logger.Error("failed to create uuid", zap.Error(err))
		common.ErrResponse(w, http.StatusBadRequest, error.Error(err))
		return
	}
	cart, err := c.svc.GetCart(mockID)
	if err != nil {
		c.logger.Error("failed to get cart", zap.Error(err))
		common.ErrResponse(w, http.StatusBadRequest, error.Error(err))
		return
	}
	common.OkResponse(w, cart)
}

func (c *Controller) AddToCart(w http.ResponseWriter, r *http.Request) {
	defer func() {
		err := r.Body.Close()
		if err != nil {
			c.logger.Error("failed to close body", zap.Error(err))
		}
	}()
	var item AddToCartRequest
	err := json.NewDecoder(r.Body).Decode(&item)
	if err != nil {
		if err != nil {
			c.logger.Error("failed add to cart", zap.Error(err))
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		err := c.svc.AddToCart(item)
		if err != nil {
			c.logger.Error("failed add to cart", zap.Error(err))
			common.ErrResponse(w, http.StatusBadRequest, error.Error(err))
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
	}
}

func (c *Controller) UpdateCart(w http.ResponseWriter, r *http.Request) {
	defer func() {
		err := r.Body.Close()
		if err != nil {
			c.logger.Error("failed to close body", zap.Error(err))
		}
	}()
	var item UpdateCartItemRequest
	err := json.NewDecoder(r.Body).Decode(&item)
	if err != nil {
		c.logger.Error("failed update cart", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = c.svc.UpdateCart(item)
	if err != nil {
		c.logger.Error("failed to update cart", zap.Error(err))
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
