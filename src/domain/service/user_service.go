package service

import (
	"context"

	"github.com/nurcahyaari/ecommerce/src/transferobject"
)

type UserServicer interface {
	GetUser(ctx context.Context, request transferobject.RequestSearchUser) (transferobject.ResponseGetUser, error)
	SearchUsers(ctx context.Context, request transferobject.RequestSearchUser) (transferobject.ResponseSearchUser, error)
}

type UserAddressServicer interface {
	GetUserAddress(ctx context.Context, request transferobject.RequestSearchUserAddress) (transferobject.ResponseGetUserAddress, error)
	GetUserAddresses(ctx context.Context, request transferobject.RequestSearchUserAddress) (transferobject.ResponseSearchUserAddress, error)
}
