package repository

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/nurcahyaari/ecommerce/infrastructure/database"
	"github.com/nurcahyaari/ecommerce/src/domain/entity"
	"github.com/nurcahyaari/ecommerce/src/domain/repository"
	"golang.org/x/sync/errgroup"
)

type OrderReceiptRepositoryRead struct {
	db *database.SQLDatabase
}

func NewOrderReceiptRepositoryRead(db *database.SQLDatabase) repository.OrderReceiptRepositoryReader {
	return &OrderReceiptRepositoryRead{
		db: db,
	}
}

func (r *OrderReceiptRepositoryRead) GetOrderReceipts(ctx context.Context, filter entity.OrderReceiptFilter) (entity.OrderReceipts, entity.Pagination, error) {
	var (
		orderReceipts entity.OrderReceipts
		count         int
		query         = OrderReceiptQuery.Select
		queryCount    = OrderReceiptQuery.Count
		argsSelect    = make([]interface{}, 0)
		argsCount     = make([]interface{}, 0)
	)

	whereClause, args, err := filter.ComposeFilter()
	if err != nil {
		return nil, entity.Pagination{}, err
	}

	argsSelect = append(argsSelect, args...)
	argsCount = append(argsCount, args...)

	query = query + " " + whereClause
	queryCount = queryCount + " " + whereClause

	if filter.Pagination != nil {
		limitClause, limitArgs := filter.Pagination.Pagination()

		if len(limitArgs) > 0 {
			argsSelect = append(argsSelect, limitArgs...)
		}
		query = query + " " + limitClause
	}

	errGroup, ctx := errgroup.WithContext(ctx)

	errGroup.Go(func() error {
		query, args, err := sqlx.In(query, argsSelect...)
		if err != nil {
			return err
		}

		err = r.db.DB.SelectContext(ctx, &orderReceipts, query, args...)
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
		return orderReceipts, entity.Pagination{}, err
	}

	pagination := entity.Pagination{}
	if filter.Pagination != nil {
		pagination = entity.NewPagination(count, filter.Pagination.Size)
	}

	return orderReceipts, pagination, err
}
