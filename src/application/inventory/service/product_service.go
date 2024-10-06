package service

import (
	"context"
	"errors"
	"net/http"

	"github.com/nurcahyaari/ecommerce/config"
	internalerrors "github.com/nurcahyaari/ecommerce/internal/x/errors"
	"github.com/nurcahyaari/ecommerce/src/domain/repository"
	"github.com/nurcahyaari/ecommerce/src/domain/service"
	"github.com/nurcahyaari/ecommerce/src/transferobject"
	"github.com/rs/zerolog"
)

type ProductService struct {
	cfg                   *config.Config
	log                   zerolog.Logger
	repositoryReader      repository.ProductRepositoryReader
	repositoryWritter     repository.ProductRepositoryWriter
	stockRepositoryWriter repository.ProductStockRepositoryWriter
	repositoryAggregator  repository.ProductAggregator
	warehouseServicer     service.WarehouseServicer
}

func NewProductService(
	cfg *config.Config,
	log zerolog.Logger,
	repositoryReader repository.ProductRepositoryReader,
	repositoryWritter repository.ProductRepositoryWriter,
	stockRepositoryWriter repository.ProductStockRepositoryWriter,
	repositoryAggregator repository.ProductAggregator,
	warehouseServicer service.WarehouseServicer,
) service.ProductServicer {
	return &ProductService{
		cfg:                   cfg,
		log:                   log,
		repositoryReader:      repositoryReader,
		repositoryWritter:     repositoryWritter,
		stockRepositoryWriter: stockRepositoryWriter,
		repositoryAggregator:  repositoryAggregator,
		warehouseServicer:     warehouseServicer,
	}
}

func (s *ProductService) GetProduct(ctx context.Context, request transferobject.RequestSearchProduct) (transferobject.ResponseGetProduct, error) {
	filter, err := request.ProductFilter()
	if err != nil {
		s.log.Error().
			Err(err).
			Msg("GetProduct.UserFilter")
		return transferobject.ResponseGetProduct{}, err
	}

	products, _, err := s.repositoryReader.FindProduct(ctx, filter)
	if err != nil {
		s.log.Error().
			Err(err).
			Msg("GetProduct.FindProduct")
		return transferobject.ResponseGetProduct{}, err
	}

	product, status := products.One()
	if !status {
		s.log.Warn().
			Any("request", request).
			Msg("GetProduct.One")
		return transferobject.ResponseGetProduct{}, nil
	}

	return transferobject.NewResponseGetProduct(product), nil
}

func (s *ProductService) SearchProducts(ctx context.Context, request transferobject.RequestSearchProduct) (transferobject.ResponseSearchProduct, error) {
	filter, err := request.ProductFilter()
	if err != nil {
		s.log.Error().
			Err(err).
			Msg("SearchUser.UserFilter")
		return transferobject.ResponseSearchProduct{}, err
	}
	products, pagination, err := s.repositoryAggregator.FindProduct(ctx, filter)
	if err != nil {
		s.log.Error().
			Err(err).
			Msg("SearchProducts.FindProduct")
		return transferobject.ResponseSearchProduct{}, err
	}

	return transferobject.NewResponseSearchProduct(products, pagination), nil
}

func (s *ProductService) MoveWarehouse(ctx context.Context, request transferobject.RequestMoveWarehouse) (transferobject.Product, error) {
	products, _, err := s.repositoryReader.FindProduct(ctx, request.ProductFilter())
	if err != nil {
		s.log.Error().
			Err(err).
			Interface("request.body", request).
			Msg("MoveWarehouse.FindProduct")
		return transferobject.Product{}, err
	}
	product, exist := products.One()
	if !exist {
		return transferobject.Product{}, nil
	}

	respGetWarehouse, err := s.warehouseServicer.GetWarehouse(ctx, transferobject.RequestSearchWarehouse{
		Ids:      request.WarehouseTargetIdStr(),
		StoreIds: product.StoreIdStr(),
	})
	if err != nil {
		s.log.Error().
			Err(err).
			Interface("request.body", request).
			Msg("MoveWarehouse.GetWarehouse")
		return transferobject.Product{}, err
	}
	if respGetWarehouse.Warehouse.Id == 0 {
		s.log.Warn().
			Interface("request.body", request).
			Msg("Warehouse is not found")
		return transferobject.Product{}, internalerrors.New(
			errors.New("err: warehouse is not found"),
			internalerrors.SetErrorCode(http.StatusNotFound))
	}

	product.MoveWarehouse(respGetWarehouse.Warehouse.Entity())
	if err := s.repositoryWritter.MoveWarehouse(ctx, product); err != nil {
		s.log.Error().
			Err(err).
			Interface("request.body", request).
			Msg("MoveWarehouse.MoveWarehouse")
		return transferobject.Product{}, err
	}

	return transferobject.NewProduct(product), nil
}

func (s *ProductService) AddReserveStock(ctx context.Context, request transferobject.RequestReserveStoct) (transferobject.ResponseReserveStock, error) {
	resp, err := s.stockRepositoryWriter.ReserveStocks(ctx, request.Data.ProductStocks())
	if err != nil {
		s.log.Error().
			Err(err).
			Interface("request.body", request).
			Msg("AddReserveStock.ReserveStocks")
	}
	return transferobject.NewResponseReserveStock(resp), nil
}
