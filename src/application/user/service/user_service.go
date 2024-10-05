package service

import (
	"context"

	"github.com/nurcahyaari/ecommerce/config"
	"github.com/nurcahyaari/ecommerce/src/domain/repository"
	"github.com/nurcahyaari/ecommerce/src/domain/service"
	"github.com/nurcahyaari/ecommerce/src/transferobject"
	"github.com/rs/zerolog"
)

type UserService struct {
	cfg              *config.Config
	log              zerolog.Logger
	RepositoryReader repository.UserRepositoryReader
}

func NewUserService(
	cfg *config.Config,
	log zerolog.Logger,
	repositoryReader repository.UserRepositoryReader,
) service.UserServicer {
	return &UserService{
		cfg:              cfg,
		log:              log,
		RepositoryReader: repositoryReader,
	}
}

func (s *UserService) GetUser(ctx context.Context, request transferobject.RequestSearchUser) (transferobject.ResponseGetUser, error) {
	filter, err := request.UserFilter()
	if err != nil {
		s.log.Error().
			Err(err).
			Msg("GetUser.UserFilter")
		return transferobject.ResponseGetUser{}, err
	}

	users, _, err := s.RepositoryReader.FindUsers(ctx, filter)
	if err != nil {
		s.log.Error().
			Err(err).
			Msg("GetUser.FindUser")
		return transferobject.ResponseGetUser{}, err
	}

	user, status := users.One()
	if !status {
		s.log.Warn().
			Any("request", request).
			Msg("GetUser.One")
		return transferobject.ResponseGetUser{}, nil
	}

	return transferobject.NewResponseGetUser(user), nil
}

func (s *UserService) SearchUsers(ctx context.Context, request transferobject.RequestSearchUser) (transferobject.ResponseSearchUser, error) {
	filter, err := request.UserFilter()
	if err != nil {
		s.log.Error().
			Err(err).
			Msg("SearchUser.UserFilter")
		return transferobject.ResponseSearchUser{}, err
	}
	users, pagination, err := s.RepositoryReader.FindUsers(ctx, filter)
	if err != nil {
		s.log.Error().
			Err(err).
			Msg("SearchUser.FindUser")
		return transferobject.ResponseSearchUser{}, err
	}

	return transferobject.NewResponseSearchUser(users, pagination), nil
}
