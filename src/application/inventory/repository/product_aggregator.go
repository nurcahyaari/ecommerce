package repository

import (
	"context"

	"github.com/nurcahyaari/ecommerce/src/domain/entity"
	"github.com/nurcahyaari/ecommerce/src/domain/repository"
)

type ProductAggregate struct {
	repositoryReader      repository.ProductRepositoryReader
	stockRepositoryReader repository.ProductStockRepositoryReader
}

func NewProductAggregate(
	repositoryReader repository.ProductRepositoryReader,
	stockRepositoryReader repository.ProductStockRepositoryReader,
) repository.ProductAggregator {
	return &ProductAggregate{
		repositoryReader:      repositoryReader,
		stockRepositoryReader: stockRepositoryReader,
	}
}

func (r *ProductAggregate) FindProduct(ctx context.Context, filter entity.ProductFilter) (entity.Products, entity.Pagination, error) {
	products, pagination, err := r.repositoryReader.FindProduct(ctx, filter)
	if err != nil {
		return nil, entity.Pagination{}, err
	}

	productStock, err := r.stockRepositoryReader.FindProductStock(ctx, products.ProductStockFilter())
	if err != nil {
		return nil, entity.Pagination{}, err
	}

	products.AggregateStock(productStock.MapByProductId())

	return products, pagination, nil
}
