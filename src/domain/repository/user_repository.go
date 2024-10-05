package repository

import (
	"context"

	"github.com/nurcahyaari/ecommerce/src/domain/entity"
)

type UserRepositoryReader interface {
	FindUsers(ctx context.Context, filter entity.UserFilter) (entity.Users, entity.Pagination, error)
}

type UserRepositoryWriter interface {
}

type UserAddressRepositoryReader interface {
	FindUserAddresses(ctx context.Context, filter entity.UserAddressFilter) (entity.UserAddresses, entity.Pagination, error)
}

type UserAddressRepositoryWriter interface {
}
