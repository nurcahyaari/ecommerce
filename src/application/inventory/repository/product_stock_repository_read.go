package repository

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/nurcahyaari/ecommerce/infrastructure/database"
	"github.com/nurcahyaari/ecommerce/src/domain/entity"
	"github.com/nurcahyaari/ecommerce/src/domain/repository"
)

type ProductStockRepositoryRead struct {
	db *database.SQLDatabase
}

func NewProductStockRepositoryRead(db *database.SQLDatabase) repository.ProductStockRepositoryReader {
	return &ProductStockRepositoryRead{
		db: db,
	}
}

func (r *ProductStockRepositoryRead) FindProductStock(ctx context.Context, filter entity.ProductStockFilter) (entity.ProductStocks, error) {
	var (
		products   entity.ProductStocks
		query      = ProductStockQuery.Select
		argsSelect = make([]interface{}, 0)
		argsCount  = make([]interface{}, 0)
	)

	whereClause, args, err := filter.ComposeFilter()
	if err != nil {
		return nil, err
	}

	argsSelect = append(argsSelect, args...)
	argsCount = append(argsCount, args...)

	query = query + " " + whereClause + " "

	query, args, err = sqlx.In(query, argsSelect...)
	if err != nil {
		return nil, err
	}

	err = r.db.DB.SelectContext(ctx, &products, query, args...)
	if err != nil {
		return nil, err
	}
	return products, err
}
