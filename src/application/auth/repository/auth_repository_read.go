package repository

import (
	"context"

	"github.com/nurcahyaari/ecommerce/infrastructure/database"
	"github.com/nurcahyaari/ecommerce/src/domain/entity"
	"github.com/nurcahyaari/ecommerce/src/domain/repository"
)

type AuthRepositoryRead struct {
	db *database.SQLDatabase
}

func NewAuthRepositoryRead(db *database.SQLDatabase) repository.AuthRepositoryReader {
	return &AuthRepositoryRead{
		db: db,
	}
}

func (r *AuthRepositoryRead) FindAuthRefreshToken(ctx context.Context, filter entity.AuthRefreshTokenFilter) (entity.AuthRefreshToken, error)
