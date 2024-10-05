package repository

import (
	"context"

	"github.com/nurcahyaari/ecommerce/infrastructure/database"
	"github.com/nurcahyaari/ecommerce/src/domain/entity"
)

type WarehouseRepositoryReader interface {
	FindWarehouses(ctx context.Context, filter entity.WarehouseFilter) (entity.Warehouses, entity.Pagination, error)
}

type WarehouseRepositoryWriter interface {
	database.SQLDatabaseTrx[WarehouseRepositoryWriter]
	UpdateWarehouseActivationStatus(context.Context, entity.Warehouse) error
	UpdateWarehousesActivationStatus(context.Context, entity.Warehouses) error
}
