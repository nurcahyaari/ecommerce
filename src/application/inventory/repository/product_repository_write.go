package repository

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/nurcahyaari/ecommerce/infrastructure/database"
	"github.com/nurcahyaari/ecommerce/src/domain/entity"
	"github.com/nurcahyaari/ecommerce/src/domain/repository"
)

type ProductRepositoryWrite struct {
	db *database.SQLDatabase
}

func NewProductRepositoryWrite(db *database.SQLDatabase) repository.ProductRepositoryWriter {
	return &ProductRepositoryWrite{
		db: db,
	}
}

func (r *ProductRepositoryWrite) MoveWarehouse(ctx context.Context, product entity.Product) error {
	query, args, err := sqlx.Named(Query.UpdateWarehouse, product)
	if err != nil {
		return err
	}

	query, args, err = sqlx.In(query, args...)
	if err != nil {
		return err
	}

	query = r.db.DB.Rebind(query)

	stmt, err := r.db.DB.PrepareContext(ctx, query)
	if err != nil {
		return err
	}

	_, err = stmt.ExecContext(ctx, args...)
	if err != nil {
		return err
	}
	return nil
}
