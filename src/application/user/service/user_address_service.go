package service

import (
	"context"

	"github.com/nurcahyaari/ecommerce/config"
	"github.com/nurcahyaari/ecommerce/src/domain/repository"
	"github.com/nurcahyaari/ecommerce/src/domain/service"
	"github.com/nurcahyaari/ecommerce/src/transferobject"
	"github.com/rs/zerolog"
)

type UserAddressService struct {
	cfg              *config.Config
	log              zerolog.Logger
	repositoryReader repository.UserAddressRepositoryReader
}

func NewUserAddressService(
	cfg *config.Config,
	log zerolog.Logger,
	repositoryReader repository.UserAddressRepositoryReader,
) service.UserAddressServicer {
	return &UserAddressService{
		cfg:              cfg,
		log:              log,
		repositoryReader: repositoryReader,
	}
}

func (s *UserAddressService) GetUserAddress(ctx context.Context, request transferobject.RequestSearchUserAddress) (transferobject.ResponseGetUserAddress, error) {
	filter, err := request.UserFilter()
	if err != nil {
		s.log.Error().
			Err(err).
			Msg("GetUserAddress.UserFilter")
		return transferobject.ResponseGetUserAddress{}, err
	}

	users, _, err := s.repositoryReader.FindUserAddresses(ctx, filter)
	if err != nil {
		s.log.Error().
			Err(err).
			Msg("GetUserAddress.FindUserAddresses")
		return transferobject.ResponseGetUserAddress{}, err
	}

	userAddress, status := users.One()
	if !status {
		s.log.Warn().
			Any("request", request).
			Msg("GetUserAddress.One")
		return transferobject.ResponseGetUserAddress{}, nil
	}

	return transferobject.NewResponseGetUserAddress(userAddress), nil
}
func (s *UserAddressService) GetUserAddresses(ctx context.Context, request transferobject.RequestSearchUserAddress) (transferobject.ResponseSearchUserAddress, error) {
	filter, err := request.UserFilter()
	if err != nil {
		s.log.Error().
			Err(err).
			Msg("GetUserAddresses.UserFilter")
		return transferobject.ResponseSearchUserAddress{}, err
	}
	userAddress, pagination, err := s.repositoryReader.FindUserAddresses(ctx, filter)
	if err != nil {
		s.log.Error().
			Err(err).
			Msg("GetUserAddresses.FindUser")
		return transferobject.ResponseSearchUserAddress{}, err
	}

	return transferobject.NewResponseSearchUserAddress(userAddress, pagination), nil
}
