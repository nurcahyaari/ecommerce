package service

import (
	"context"

	"github.com/nurcahyaari/ecommerce/src/transferobject"
)

type OrderServicer interface {
	CreateOrder(ctx context.Context, request transferobject.RequestCreateOrder) (transferobject.Order, error)
	ExpiredOrder(ctx context.Context) error
}
