package repository

import (
	"context"

	"github.com/nurcahyaari/ecommerce/src/domain/entity"
)

type CartRepositorier interface {
	FindCart(ctx context.Context, filter entity.CartFilter) (entity.Carts, error)
	UpsertCart(ctx context.Context, cart entity.Cart) error
	DeleteCart(ctx context.Context, filter entity.CartFilter) error
}
