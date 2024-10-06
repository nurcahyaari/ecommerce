package repository

import (
	"context"

	"github.com/jmoiron/sqlx"
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

func (r *ProductStockRepositoryWrite) BeginTx(ctx context.Context) (repository.ProductStockRepositoryWriter, error) {
	tx, err := r.db.DB.BeginTxx(ctx, nil)
	if err != nil {
		return nil, err
	}

	db := database.SQLDatabase{
		Tx: tx,
	}

	return NewProductStockRepositoryWrite(&db), nil
}

func (r *ProductStockRepositoryWrite) Commit(ctx context.Context) error {
	if r.db.Tx == nil {
		return nil
	}

	if err := r.db.Tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (r *ProductStockRepositoryWrite) Rollback(ctx context.Context) error {
	if r.db.Tx == nil {
		return nil
	}

	if err := r.db.Tx.Rollback(); err != nil {
		return err
	}

	return nil
}

// ReserveStock will move the stock on hand to stock reserved and decrease the stock on hand
func (r *ProductStockRepositoryWrite) ReserveStock(ctx context.Context, request entity.ProductStock) error {
	var (
		query = ProductStockQuery.ReserveStock
		stmt  *sqlx.Stmt
		err   error
	)

	if request.IsRevert {
		query = ProductStockQuery.RevertStock
	}

	if r.db.Tx != nil {
		stmt, err = r.db.Tx.Preparex(query)
	} else {
		stmt, err = r.db.DB.Preparex(query)
	}
	if err != nil {
		return err
	}

	_, err = stmt.ExecContext(ctx, request.ReserveStockArgs()...)
	if err != nil {
		return err
	}

	return nil
}

func (r *ProductStockRepositoryWrite) ReserveStocks(ctx context.Context, request entity.ProductStocks) (entity.ReserveStocks, error) {
	repo, err := r.BeginTx(ctx)
	if err != nil {
		return entity.ReserveStocks{}, err
	}

	var (
		resps = make(entity.ReserveStocks, 0)
	)

	defer func() {
		if p := recover(); p != nil {
			repo.Rollback(ctx)
			panic(p)
		} else if err != nil {
			repo.Rollback(ctx)
		}
	}()

	for _, r := range request {
		resp := entity.ReserveStock{
			ProductStock: r,
		}
		err = repo.ReserveStock(ctx, r)
		if err != nil {
			resp.Message = err.Error()
			resp.Success = false
			continue
		}

		resp.Message = "success to reserve stock"
		resp.Success = true
	}

	if err != nil {
		return resps, err
	}

	err = repo.Commit(ctx)
	if err != nil {
		return entity.ReserveStocks{}, err
	}

	return resps, nil
}

// UpdateStock will decrease the stock reserved
func (r *ProductStockRepositoryWrite) UpdateStock(ctx context.Context, request entity.ProductStock) error {
	var (
		query = ProductStockQuery.UpdateStock
		stmt  *sqlx.Stmt
		err   error
	)

	if r.db.Tx != nil {
		stmt, err = r.db.Tx.Preparex(query)
	} else {
		stmt, err = r.db.DB.Preparex(query)
	}
	if err != nil {
		return err
	}

	_, err = stmt.ExecContext(ctx, request.UpdateStockArgs()...)
	if err != nil {
		return err
	}

	return nil
}

func (r *ProductStockRepositoryWrite) UpdateStocks(ctx context.Context, request entity.ProductStocks) (entity.ReserveStocks, error) {
	return entity.ReserveStocks{}, nil
}
