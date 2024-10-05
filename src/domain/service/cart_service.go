package service

import (
	"context"

	"github.com/nurcahyaari/ecommerce/src/transferobject"
)

type CartServicer interface {
	GetCart(ctx context.Context, request transferobject.RequestGetCart) (transferobject.ResponseGetCart, error)
	AddItemToCart(ctx context.Context, request transferobject.RequestAddItemToCart) (transferobject.Cart, error)
}
