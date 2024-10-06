package repository

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
	"github.com/nurcahyaari/ecommerce/infrastructure/database"
	"github.com/nurcahyaari/ecommerce/src/domain/entity"
	"github.com/nurcahyaari/ecommerce/src/domain/repository"
)

type OrderDetailRepositoryWrite struct {
	db *database.SQLDatabase
}

func NewOrderDetailRepositoryWrite(db *database.SQLDatabase) repository.OrderDetailRepositoryWriter {
	return &OrderDetailRepositoryWrite{
		db: db,
	}
}

func (r *OrderDetailRepositoryWrite) BeginTx(ctx context.Context) (repository.OrderDetailRepositoryWriter, error) {
	tx, err := r.db.DB.BeginTxx(ctx, nil)
	if err != nil {
		return nil, err
	}

	db := database.SQLDatabase{
		Tx: tx,
	}

	return NewOrderDetailRepositoryWrite(&db), nil
}

func (r *OrderDetailRepositoryWrite) Commit(ctx context.Context) error {
	if r.db.Tx == nil {
		return nil
	}

	if err := r.db.Tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (r *OrderDetailRepositoryWrite) Rollback(ctx context.Context) error {
	if r.db.Tx == nil {
		return nil
	}

	if err := r.db.Tx.Rollback(); err != nil {
		return err
	}

	return nil
}

func (r *OrderDetailRepositoryWrite) CreateOrderDetails(ctx context.Context, orderDetails entity.OrderDetails) (entity.OrderDetails, error) {
	if len(orderDetails) == 0 {
		return entity.OrderDetails{}, nil
	}

	for i, orderDetail := range orderDetails {
		query, args, err := sqlx.Named(OrderDetailQuery.Insert, orderDetail)
		if err != nil {
			return entity.OrderDetails{}, err
		}

		query, args, err = sqlx.In(query, args...)
		if err != nil {
			return entity.OrderDetails{}, err
		}

		stmt := &sql.Stmt{}
		if r.db.Tx != nil {
			query = r.db.DB.Rebind(query)
			stmt, err = r.db.DB.PrepareContext(ctx, query)
			if err != nil {
				return entity.OrderDetails{}, err
			}
		} else {
			query = r.db.DB.Rebind(query)
			stmt, err = r.db.DB.PrepareContext(ctx, query)
			if err != nil {
				return entity.OrderDetails{}, err
			}
		}

		result, err := stmt.ExecContext(ctx, args...)
		if err != nil {
			return entity.OrderDetails{}, err
		}

		orderDetailId, err := result.LastInsertId()
		if err != nil {
			return entity.OrderDetails{}, err
		}

		orderDetails[i].Id = orderDetailId
	}

	return orderDetails, nil
}
