package repository

import "github.com/nurcahyaari/ecommerce/src/domain/repository"

type UserRepositoryWrite struct {
}

func NewUserRepositoryWrite() repository.UserRepositoryWriter {
	return &UserRepositoryWrite{}
}
