package repository

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
	"github.com/nurcahyaari/ecommerce/infrastructure/database"
	"github.com/nurcahyaari/ecommerce/src/domain/entity"
	"github.com/nurcahyaari/ecommerce/src/domain/repository"
)

type OrderRepositoryWrite struct {
	db *database.SQLDatabase
}

func NewOrderRepositoryWrite(db *database.SQLDatabase) repository.OrderRepositoryWriter {
	return &OrderRepositoryWrite{
		db: db,
	}
}

func (r *OrderRepositoryWrite) BeginTx(ctx context.Context) (repository.OrderRepositoryWriter, error) {
	tx, err := r.db.DB.BeginTxx(ctx, nil)
	if err != nil {
		return nil, err
	}

	db := database.SQLDatabase{
		DB: r.db.DB,
		Tx: tx,
	}

	return NewOrderRepositoryWrite(&db), nil
}

func (r *OrderRepositoryWrite) Commit(ctx context.Context) error {
	if r.db.Tx == nil {
		return nil
	}

	if err := r.db.Tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (r *OrderRepositoryWrite) Rollback(ctx context.Context) error {
	if r.db.Tx == nil {
		return nil
	}

	if err := r.db.Tx.Rollback(); err != nil {
		return err
	}

	return nil
}

func (r *OrderRepositoryWrite) CreateOrder(ctx context.Context, data entity.Order) (entity.Order, error) {
	query, args, err := sqlx.Named(OrderQuery.Create, data)
	if err != nil {
		return entity.Order{}, nil
	}

	query, args, err = sqlx.In(query, args...)
	if err != nil {
		return entity.Order{}, err
	}

	stmt := &sql.Stmt{}
	if r.db.Tx != nil {
		query = r.db.DB.Rebind(query)
		stmt, err = r.db.DB.PrepareContext(ctx, query)
		if err != nil {
			return entity.Order{}, err
		}
	} else {
		query = r.db.DB.Rebind(query)
		stmt, err = r.db.DB.PrepareContext(ctx, query)
		if err != nil {
			return entity.Order{}, err
		}
	}

	res, err := stmt.ExecContext(ctx, args...)
	if err != nil {
		return entity.Order{}, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return entity.Order{}, err
	}

	data.Id = id

	return data, nil
}

func (r *OrderRepositoryWrite) UpdateOrderStatus(ctx context.Context, order entity.Order) error {
	query, args, err := sqlx.Named(OrderQuery.UpdateOrderStatus, order)
	if err != nil {
		return nil
	}

	query, args, err = sqlx.In(query, args...)
	if err != nil {
		return err
	}

	stmt := &sql.Stmt{}
	if r.db.Tx != nil {
		query = r.db.DB.Rebind(query)
		stmt, err = r.db.DB.PrepareContext(ctx, query)
		if err != nil {
			return err
		}
	} else {
		query = r.db.DB.Rebind(query)
		stmt, err = r.db.DB.PrepareContext(ctx, query)
		if err != nil {
			return err
		}
	}

	_, err = stmt.ExecContext(ctx, args...)
	if err != nil {
		return err
	}

	return nil
}

func (r *OrderRepositoryWrite) UpdateOrdersStatus(ctx context.Context, orders entity.Orders) error {
	orderRepo, err := r.BeginTx(ctx)
	if err != nil {
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			orderRepo.Rollback(ctx)
			panic(p)
		} else if err != nil {
			orderRepo.Rollback(ctx)
		}
	}()

	for _, order := range orders {
		err = orderRepo.UpdateOrderStatus(ctx, order)
		if err != nil {
			return err
		}
	}

	return nil
}
