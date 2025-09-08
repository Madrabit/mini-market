package internal

import (
	"github.com/google/uuid"
	"time"
)

type SearchRequest struct {
	Query   string            `json:"query" validate:"required"`
	Filters map[string]string `json:"filters,omitempty"` // "brand":"Nike"
	Sort    string            `json:"sort,omitempty"`    // "price_asc"
	Limit   int               `json:"limit,omitempty"`
	Offset  int               `json:"offset,omitempty"`
}

type SearchItem struct {
	ProductID uuid.UUID `json:"product_id"`
	Name      string    `json:"name"`
	Price     int64     `json:"price"`
	Score     float64   `json:"score"` // релевантность
}

type SearchResponse struct {
	Total int          `json:"total"`
	Items []SearchItem `json:"items"`
}

//подсказки в поиске

type SuggestRequest struct {
	Query string `json:"query" validate:"required"`
	Limit int    `json:"limit,omitempty"`
}

type SuggestResponse struct {
	Suggestions []string `json:"suggestions"`
}

// Для батчевого обновления из Catalog
type ReindexRequest struct {
	Force bool `json:"force,omitempty"` // если true — полная пересборка
}

type ReindexResponse struct {
	Indexed    int `json:"indexed"`
	Skipped    int `json:"skipped"`
	DurationMs int `json:"duration_ms"`
}

// Переиндексация из каталога

type ProductForIndex struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Price     int64     `json:"price"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CatalogBatchResponse struct {
	Total      int               `json:"total"`
	Items      []ProductForIndex `json:"items"`
	NextCursor string            `json:"next_cursor,omitempty"`
}
