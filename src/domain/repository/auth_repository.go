package repository

import (
	"context"

	"github.com/nurcahyaari/ecommerce/src/domain/entity"
)

type AuthRepositoryReader interface {
	FindAuthRefreshToken(ctx context.Context, filter entity.AuthRefreshTokenFilter) (entity.AuthRefreshToken, error)
}

type AuthRepositoryWriter interface {
	WriteAuthRefreshToken(ctx context.Context, authRefreshToken entity.AuthRefreshToken) error
}
