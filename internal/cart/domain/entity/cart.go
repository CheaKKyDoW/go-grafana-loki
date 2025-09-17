package entity

import (
	"time"

	"github.com/google/uuid"
)

type Cart struct {
	ID        uuid.UUID  `json:"id" db:"id"`
	UserID    uuid.UUID  `json:"user_id" db:"user_id"`
	Items     []CartItem `json:"items"`
	Total     float64    `json:"total" db:"total"`
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt time.Time  `json:"updated_at" db:"updated_at"`
}

type CartItem struct {
	ID        uuid.UUID `json:"id" db:"id"`
	CartID    uuid.UUID `json:"cart_id" db:"cart_id"`
	ProductID uuid.UUID `json:"product_id" db:"product_id"`
	Quantity  int       `json:"quantity" db:"quantity"`
	Price     float64   `json:"price" db:"price"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type AddItemRequest struct {
	ProductID uuid.UUID `json:"product_id" binding:"required"`
	Quantity  int       `json:"quantity" binding:"required,min=1"`
	Price     float64   `json:"price" binding:"required,min=0"`
}

type UpdateItemRequest struct {
	Quantity int `json:"quantity" binding:"required,min=1"`
}

type CartResponse struct {
	ID        uuid.UUID  `json:"id"`
	UserID    uuid.UUID  `json:"user_id"`
	Items     []CartItem `json:"items"`
	Total     float64    `json:"total"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}
