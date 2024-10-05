package repository

import (
	"context"

	"github.com/nurcahyaari/ecommerce/infrastructure/database"
	"github.com/nurcahyaari/ecommerce/src/domain/entity"
	"github.com/nurcahyaari/ecommerce/src/domain/repository"
)

type AuthRepositoryWrite struct {
	db *database.SQLDatabase
}

func NewAuthRepositoryWrite(db *database.SQLDatabase) repository.AuthRepositoryWriter {
	return &AuthRepositoryWrite{
		db: db,
	}
}

func (r *AuthRepositoryWrite) WriteAuthRefreshToken(ctx context.Context, authRefreshToken entity.AuthRefreshToken) error
