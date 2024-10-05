package repository

import (
	"context"

	"github.com/nurcahyaari/ecommerce/src/domain/entity"
)

type ProductRepositoryReader interface {
	FindProduct(context.Context, entity.ProductFilter) (entity.Products, entity.Pagination, error)
}

type ProductRepositoryWriter interface {
	MoveWarehouse(context.Context, entity.Product) error
}

type ProductStockRepositoryReader interface {
	FindProductStock(context.Context, entity.ProductStockFilter) (entity.ProductStocks, error)
}

type ProductStockRepositoryWriter interface {
	// ReserveStock will move the stock on hand to stock reserved and decrease the stock on hand
	ReserveStock(context.Context, entity.ProductStock) error
	// UpdateStock will decrease the stock reserved
	UpdateStock(context.Context, entity.ProductStock) error
}

type ProductAggregator interface {
	FindProduct(context.Context, entity.ProductFilter) (entity.Products, entity.Pagination, error)
}
