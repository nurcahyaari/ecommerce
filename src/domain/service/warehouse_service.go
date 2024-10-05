package service

import (
	"context"

	"github.com/nurcahyaari/ecommerce/src/transferobject"
)

type WarehouseServicer interface {
	GetWarehouse(ctx context.Context, request transferobject.RequestSearchWarehouse) (transferobject.ResponseGetWarehouse, error)
	SearchWarehouses(ctx context.Context, request transferobject.RequestSearchWarehouse) (transferobject.ResponseSearchWarehouse, error)
	OpenCloseWarehouse(ctx context.Context, request transferobject.RequestOpenCloseWarehouse) (transferobject.Warehouses, error)
}
