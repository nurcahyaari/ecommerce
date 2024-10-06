package repository

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
	"github.com/nurcahyaari/ecommerce/infrastructure/database"
	"github.com/nurcahyaari/ecommerce/src/domain/entity"
	"github.com/nurcahyaari/ecommerce/src/domain/repository"
)

type OrderReceiptRepositoryWrite struct {
	db *database.SQLDatabase
}

func NewOrderReceiptRepositoryWrite(db *database.SQLDatabase) repository.OrderReceiptRepositoryWriter {
	return &OrderReceiptRepositoryWrite{
		db: db,
	}
}

func (r *OrderReceiptRepositoryWrite) BeginTx(ctx context.Context) (repository.OrderReceiptRepositoryWriter, error) {
	tx, err := r.db.DB.BeginTxx(ctx, nil)
	if err != nil {
		return nil, err
	}

	db := database.SQLDatabase{
		Tx: tx,
	}

	return NewOrderReceiptRepositoryWrite(&db), nil
}

func (r *OrderReceiptRepositoryWrite) Commit(ctx context.Context) error {
	if r.db.Tx == nil {
		return nil
	}

	if err := r.db.Tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (r *OrderReceiptRepositoryWrite) Rollback(ctx context.Context) error {
	if r.db.Tx == nil {
		return nil
	}

	if err := r.db.Tx.Rollback(); err != nil {
		return err
	}

	return nil
}

func (r *OrderReceiptRepositoryWrite) CreateOrderReceipts(ctx context.Context, orderReceipts entity.OrderReceipts) (entity.OrderReceipts, error) {

	if len(orderReceipts) == 0 {
		return entity.OrderReceipts{}, nil
	}

	for i, d := range orderReceipts {
		query, args, err := sqlx.Named(OrderReceiptQuery.Insert, d)
		if err != nil {
			return entity.OrderReceipts{}, err
		}

		query, args, err = sqlx.In(query, args...)
		if err != nil {
			return entity.OrderReceipts{}, err
		}

		stmt := &sql.Stmt{}
		if r.db.Tx != nil {
			query = r.db.DB.Rebind(query)
			stmt, err = r.db.DB.PrepareContext(ctx, query)
			if err != nil {
				return entity.OrderReceipts{}, err
			}
		} else {
			query = r.db.DB.Rebind(query)
			stmt, err = r.db.DB.PrepareContext(ctx, query)
			if err != nil {
				return entity.OrderReceipts{}, err
			}
		}

		result, err := stmt.ExecContext(ctx, args...)
		if err != nil {
			return entity.OrderReceipts{}, err
		}

		orderReceiptId, err := result.LastInsertId()
		if err != nil {
			return entity.OrderReceipts{}, err
		}

		orderReceipts[i].Id = orderReceiptId

	}

	return orderReceipts, nil
}
