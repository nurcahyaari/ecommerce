package service

import (
	"context"

	"github.com/nurcahyaari/ecommerce/src/transferobject"
)

type ProductServicer interface {
	GetProduct(ctx context.Context, request transferobject.RequestSearchProduct) (transferobject.ResponseGetProduct, error)
	SearchProducts(ctx context.Context, request transferobject.RequestSearchProduct) (transferobject.ResponseSearchProduct, error)
	MoveWarehouse(ctx context.Context, request transferobject.RequestMoveWarehouse) (transferobject.Product, error)
	AddReserveStock(ctx context.Context, request transferobject.RequestReserveStoct) (transferobject.ResponseReserveStock, error)
}
