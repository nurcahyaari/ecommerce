package repository

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/nurcahyaari/ecommerce/infrastructure/database"
	"github.com/nurcahyaari/ecommerce/src/domain/entity"
	"github.com/nurcahyaari/ecommerce/src/domain/repository"
	"golang.org/x/sync/errgroup"
)

type ProductRepositoryRead struct {
	db *database.SQLDatabase
}

func NewProductRepositoryRead(db *database.SQLDatabase) repository.ProductRepositoryReader {
	return &ProductRepositoryRead{
		db: db,
	}
}

func (r *ProductRepositoryRead) FindProduct(ctx context.Context, filter entity.ProductFilter) (entity.Products, entity.Pagination, error) {
	var (
		products   entity.Products
		count      int
		query      = Query.Select
		queryCount = Query.Count
		argsSelect = make([]interface{}, 0)
		argsCount  = make([]interface{}, 0)
	)

	whereClause, args, err := filter.ComposeFilter()
	if err != nil {
		return nil, entity.Pagination{}, err
	}

	argsSelect = append(argsSelect, args...)
	argsCount = append(argsCount, args...)

	limitClause, limitArgs := filter.Pagination.Pagination()

	if len(limitArgs) > 0 {
		argsSelect = append(argsSelect, limitArgs...)
	}

	errGroup, ctx := errgroup.WithContext(ctx)
	query = query + " " + whereClause + " " + limitClause
	queryCount = queryCount + " " + whereClause

	errGroup.Go(func() error {
		query, args, err := sqlx.In(query, argsSelect...)
		if err != nil {
			return err
		}

		err = r.db.DB.SelectContext(ctx, &products, query, args...)
		if err != nil {
			return err
		}
		return nil
	})

	errGroup.Go(func() error {
		query, args, err := sqlx.In(queryCount, argsCount...)
		if err != nil {
			return err
		}

		err = r.db.DB.GetContext(ctx, &count, query, args...)
		if err != nil {
			return err
		}
		return nil
	})

	if err := errGroup.Wait(); err != nil {
		return products, entity.Pagination{}, err
	}

	return products, entity.NewPagination(count, filter.Size), err
}
