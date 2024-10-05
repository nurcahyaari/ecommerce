package service

import (
	"context"

	"github.com/nurcahyaari/ecommerce/config"
	"github.com/nurcahyaari/ecommerce/src/domain/repository"
	"github.com/nurcahyaari/ecommerce/src/domain/service"
	"github.com/nurcahyaari/ecommerce/src/transferobject"
	"github.com/rs/zerolog"
)

type WarehouseService struct {
	cfg        *config.Config
	log        zerolog.Logger
	repoReader repository.WarehouseRepositoryReader
	repoWriter repository.WarehouseRepositoryWriter
}

func NewWarehouseService(
	cfg *config.Config,
	log zerolog.Logger,
	repoReader repository.WarehouseRepositoryReader,
	repoWriter repository.WarehouseRepositoryWriter,
) service.WarehouseServicer {
	return &WarehouseService{
		cfg:        cfg,
		log:        log,
		repoReader: repoReader,
		repoWriter: repoWriter,
	}
}

func (s *WarehouseService) GetWarehouse(ctx context.Context, request transferobject.RequestSearchWarehouse) (transferobject.ResponseGetWarehouse, error) {
	filter, err := request.WarehouseFilter()
	if err != nil {
		s.log.Error().
			Err(err).
			Msg("GetWarehouse.WarehouseFilter")
		return transferobject.ResponseGetWarehouse{}, err
	}

	warehouses, _, err := s.repoReader.FindWarehouses(ctx, filter)
	if err != nil {
		s.log.Error().
			Err(err).
			Msg("GetWarehouse.FindWarehouses")
		return transferobject.ResponseGetWarehouse{}, err
	}

	warehouse, status := warehouses.One()
	if !status {
		s.log.Warn().
			Any("request", request).
			Msg("Warehouse.One")
		return transferobject.ResponseGetWarehouse{}, nil
	}

	return transferobject.NewResponseGetWarehouse(warehouse), nil
}

func (s *WarehouseService) SearchWarehouses(ctx context.Context, request transferobject.RequestSearchWarehouse) (transferobject.ResponseSearchWarehouse, error) {
	filter, err := request.WarehouseFilter()
	if err != nil {
		s.log.Error().
			Err(err).
			Msg("SearchWarehouses.WarehouseFilter")
		return transferobject.ResponseSearchWarehouse{}, err
	}
	warehouses, pagination, err := s.repoReader.FindWarehouses(ctx, filter)
	if err != nil {
		s.log.Error().
			Err(err).
			Msg("SearchWarehouses.FindWarehouses")
		return transferobject.ResponseSearchWarehouse{}, err
	}

	return transferobject.NewResponseSearchWarehouse(warehouses, pagination), nil
}

func (s *WarehouseService) OpenCloseWarehouse(ctx context.Context, request transferobject.RequestOpenCloseWarehouse) (transferobject.Warehouses, error) {
	filter := request.WarehouseFilter()

	warehouses, _, err := s.repoReader.FindWarehouses(ctx, filter)
	if err != nil {
		s.log.Error().
			Err(err).
			Msg("OpenCloseWarehouse.FindWarehouses")
		return transferobject.Warehouses{}, err
	}

	openCloseWarehouse := request.Warehouses()

	warehouses.OpenCloseWarehouse(openCloseWarehouse.MapWarehouseById())

	err = s.repoWriter.UpdateWarehousesActivationStatus(ctx, warehouses)
	if err != nil {
		s.log.Error().
			Err(err).
			Msg("OpenCloseWarehouse.UpdateWarehousesActivationStatus")
		return transferobject.Warehouses{}, err
	}

	return transferobject.NewWarehouses(warehouses), nil
}
