package repository

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
	"github.com/nurcahyaari/ecommerce/infrastructure/database"
	"github.com/nurcahyaari/ecommerce/src/domain/entity"
	"github.com/nurcahyaari/ecommerce/src/domain/repository"
)

type WarehouseRepositoryWrite struct {
	db *database.SQLDatabase
}

func NewWarehouseRepositoryWrite(db *database.SQLDatabase) repository.WarehouseRepositoryWriter {
	return &WarehouseRepositoryWrite{
		db: db,
	}
}

func (r *WarehouseRepositoryWrite) BeginTx(ctx context.Context) (repository.WarehouseRepositoryWriter, error) {
	tx, err := r.db.DB.BeginTxx(ctx, nil)
	if err != nil {
		return nil, err
	}

	db := database.SQLDatabase{
		Tx: tx,
	}

	return NewWarehouseRepositoryWrite(&db), nil
}

func (r *WarehouseRepositoryWrite) Commit(ctx context.Context) error {
	if r.db.Tx == nil {
		return nil
	}

	if err := r.db.Tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (r *WarehouseRepositoryWrite) Rollback(ctx context.Context) error {
	if r.db.Tx == nil {
		return nil
	}

	if err := r.db.Tx.Rollback(); err != nil {
		return err
	}

	return nil
}

func (r *WarehouseRepositoryWrite) UpdateWarehouseActivationStatus(ctx context.Context, warehouse entity.Warehouse) error {
	query, args, err := sqlx.Named(WarehouseQuery.ActivationStatus, warehouse)
	if err != nil {
		return err
	}

	query, args, err = sqlx.In(query, args...)
	if err != nil {
		return err
	}

	var (
		stmt *sql.Stmt
	)

	if r.db.Tx != nil {
		query = r.db.Tx.Rebind(query)
		stmt, err = r.db.Tx.PrepareContext(ctx, query)
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

func (r *WarehouseRepositoryWrite) UpdateWarehousesActivationStatus(ctx context.Context, warehouses entity.Warehouses) error {
	repo, err := r.BeginTx(ctx)
	if err != nil {
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			repo.Rollback(ctx)
			panic(p)
		} else if err != nil {
			repo.Rollback(ctx)
		}
	}()

	for _, w := range warehouses {
		err = repo.UpdateWarehouseActivationStatus(ctx, w)
		if err != nil {
			return err
		}
	}

	err = repo.Commit(ctx)
	if err != nil {
		return err
	}
	return nil
}
