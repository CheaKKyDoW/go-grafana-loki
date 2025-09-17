package repository

import (
	"context"

	"go-clean-arch/internal/cart/domain/entity"

	"github.com/google/uuid"
)

type CartRepository interface {
	Create(ctx context.Context, cart *entity.Cart) error
	GetByUserID(ctx context.Context, userID uuid.UUID) (*entity.Cart, error)
	Update(ctx context.Context, cart *entity.Cart) error
	Delete(ctx context.Context, id uuid.UUID) error
	AddItem(ctx context.Context, item *entity.CartItem) error
	UpdateItem(ctx context.Context, item *entity.CartItem) error
	RemoveItem(ctx context.Context, cartID, itemID uuid.UUID) error
	GetItems(ctx context.Context, cartID uuid.UUID) ([]entity.CartItem, error)
}
