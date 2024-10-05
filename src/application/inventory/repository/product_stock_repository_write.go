package repository

import (
	"context"

	"github.com/nurcahyaari/ecommerce/infrastructure/database"
	"github.com/nurcahyaari/ecommerce/src/domain/entity"
	"github.com/nurcahyaari/ecommerce/src/domain/repository"
)

type ProductStockRepositoryWrite struct {
	db *database.SQLDatabase
}

func NewProductStockRepositoryWrite(db *database.SQLDatabase) repository.ProductStockRepositoryWriter {
	return &ProductStockRepositoryWrite{
		db: db,
	}
}

// ReserveStock will move the stock on hand to stock reserved and decrease the stock on hand
func (r *ProductStockRepositoryWrite) ReserveStock(ctx context.Context, request entity.ProductStock) error {
	return nil
}

// UpdateStock will decrease the stock reserved
func (r *ProductStockRepositoryWrite) UpdateStock(ctx context.Context, request entity.ProductStock) error {
	return nil
}
