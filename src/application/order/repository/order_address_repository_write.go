package repository

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
	"github.com/nurcahyaari/ecommerce/infrastructure/database"
	"github.com/nurcahyaari/ecommerce/src/domain/entity"
	"github.com/nurcahyaari/ecommerce/src/domain/repository"
)

type OrderAddressRepositoryWrite struct {
	db *database.SQLDatabase
}

func NewOrderAddressRepositoryWrite(db *database.SQLDatabase) repository.OrderAddressRepositoryWriter {
	return &OrderAddressRepositoryWrite{
		db: db,
	}
}

func (r *OrderAddressRepositoryWrite) BeginTx(ctx context.Context) (repository.OrderAddressRepositoryWriter, error) {
	tx, err := r.db.DB.BeginTxx(ctx, nil)
	if err != nil {
		return nil, err
	}

	db := database.SQLDatabase{
		Tx: tx,
	}

	return NewOrderAddressRepositoryWrite(&db), nil
}

func (r *OrderAddressRepositoryWrite) Commit(ctx context.Context) error {
	if r.db.Tx == nil {
		return nil
	}

	if err := r.db.Tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (r *OrderAddressRepositoryWrite) Rollback(ctx context.Context) error {
	if r.db.Tx == nil {
		return nil
	}

	if err := r.db.Tx.Rollback(); err != nil {
		return err
	}

	return nil
}

func (r *OrderAddressRepositoryWrite) CreateOrderAddresses(ctx context.Context, orderAddresses entity.OrderAddresses) (entity.OrderAddresses, error) {
	if len(orderAddresses) == 0 {
		return entity.OrderAddresses{}, nil
	}

	for i, orderAddress := range orderAddresses {
		query, args, err := sqlx.Named(OrderAddressQuery.Insert, orderAddress)
		if err != nil {
			return entity.OrderAddresses{}, err
		}

		query, args, err = sqlx.In(query, args...)
		if err != nil {
			return entity.OrderAddresses{}, err
		}

		stmt := &sql.Stmt{}
		if r.db.Tx != nil {
			query = r.db.DB.Rebind(query)
			stmt, err = r.db.DB.PrepareContext(ctx, query)
			if err != nil {
				return entity.OrderAddresses{}, err
			}
		} else {
			query = r.db.DB.Rebind(query)
			stmt, err = r.db.DB.PrepareContext(ctx, query)
			if err != nil {
				return entity.OrderAddresses{}, err
			}
		}

		result, err := stmt.ExecContext(ctx, args...)
		if err != nil {
			return entity.OrderAddresses{}, err
		}

		orderAddressId, err := result.LastInsertId()
		if err != nil {
			return entity.OrderAddresses{}, err
		}

		orderAddresses[i].Id = orderAddressId
	}

	return orderAddresses, nil
}
