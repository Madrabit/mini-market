package internal

import (
	"github.com/google/uuid"
	"time"
)

type Rating struct {
	ID          uuid.UUID
	ProductID   uuid.UUID
	Stars       int
	Description string
	UserID      uuid.UUID
}

type AddReviewRequest struct {
	ProductID   uuid.UUID
	Stars       int
	Description string
	UserID      uuid.UUID
}

type UpdateReviewRequest struct {
	RatingId    uuid.UUID
	Stars       int
	Description string
	UserID      uuid.UUID
}

type DeleteReviewRequest struct {
	RatingId uuid.UUID
}

type ReviewResponse struct {
	ID          uuid.UUID `json:"id"`
	ProductID   uuid.UUID `json:"product_id"`
	UserID      uuid.UUID `json:"user_id"` // Можно возвращать, если нет требований к анонимности
	Stars       int       `json:"stars"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"` // Полезно добавлять
	UpdatedAt   time.Time `json:"updated_at"`
}

type ProductReviewsResponse struct {
	ProductID uuid.UUID        `json:"product_id"`
	Average   float64          `json:"average_rating"` // Средний рейтинг — очень важно!
	Count     int              `json:"reviews_count"`  // Количество отзывов
	Reviews   []ReviewResponse `json:"reviews"`        // Слайс самих отзывов
}
