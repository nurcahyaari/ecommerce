package repository

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/nurcahyaari/ecommerce/infrastructure/database"
	"github.com/nurcahyaari/ecommerce/src/domain/entity"
	"github.com/nurcahyaari/ecommerce/src/domain/repository"
	"golang.org/x/sync/errgroup"
)

type UserRepositoryRead struct {
	db *database.SQLDatabase
}

func NewUserRepositoryRead(db *database.SQLDatabase) repository.UserRepositoryReader {
	return &UserRepositoryRead{
		db: db,
	}
}

func (r *UserRepositoryRead) FindUsers(ctx context.Context, filter entity.UserFilter) (entity.Users, entity.Pagination, error) {
	var (
		users      entity.Users
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

		err = r.db.DB.SelectContext(ctx, &users, query, args...)
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
		return users, entity.Pagination{}, err
	}

	return users, entity.NewPagination(count, filter.Size), err
}
