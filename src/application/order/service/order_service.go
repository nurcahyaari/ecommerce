package service

import (
	"context"

	"github.com/nurcahyaari/ecommerce/config"
	"github.com/nurcahyaari/ecommerce/src/domain/entity"
	"github.com/nurcahyaari/ecommerce/src/domain/repository"
	"github.com/nurcahyaari/ecommerce/src/domain/service"
	"github.com/nurcahyaari/ecommerce/src/transferobject"
	"github.com/rs/zerolog"
	"gopkg.in/guregu/null.v4"
)

type OrderService struct {
	cfg             *config.Config
	log             zerolog.Logger
	orderRepoReader repository.OrderRepositoryReader
	orderRepoWriter repository.OrderRepositoryWriter
	orderAggregator repository.OrderAggregator
	userAddressSvc  service.UserAddressServicer
	cartSvc         service.CartServicer
	productSvc      service.ProductServicer
}

func NewOrderService(
	cfg *config.Config,
	log zerolog.Logger,
	orderAggregator repository.OrderAggregator,
	userAddressSvc service.UserAddressServicer,
	cartSvc service.CartServicer,
	productSvc service.ProductServicer,
	orderRepoReader repository.OrderRepositoryReader,
	orderRepoWriter repository.OrderRepositoryWriter,
) service.OrderServicer {
	return &OrderService{
		cfg:             cfg,
		log:             log,
		orderAggregator: orderAggregator,
		cartSvc:         cartSvc,
		userAddressSvc:  userAddressSvc,
		productSvc:      productSvc,
		orderRepoReader: orderRepoReader,
		orderRepoWriter: orderRepoWriter,
	}
}

func (s *OrderService) CreateOrder(ctx context.Context, request transferobject.RequestCreateOrder) (transferobject.Order, error) {
	respCart, err := s.cartSvc.GetCart(ctx, transferobject.RequestGetCart{
		UserId: request.UserId,
	})
	if err != nil {
		s.log.Error().
			Err(err).
			Msg("CreateOrder.GetCart")
		return transferobject.Order{}, err
	}

	cart, err := respCart.Carts.Entity()
	if err != nil {
		s.log.Error().
			Err(err).
			Msg("CreateOrder.Carts.Entity")
		return transferobject.Order{}, err
	}

	respUserAddress, err := s.userAddressSvc.GetUserAddresses(ctx, transferobject.RequestSearchUserAddress{
		Ids: cart.UserAddressStrs(),
	})
	if err != nil {
		s.log.Error().
			Err(err).
			Msg("CreateOrder.userAddressSvc.GetUserAddresses")
		return transferobject.Order{}, err
	}

	userAddress := respUserAddress.UserAddresses.Entity()
	mapUserAddressById := userAddress.MapById()

	userId, err := request.UserIdInt()
	if err != nil {
		s.log.Error().
			Err(err).
			Msg("CreateOrder.UserIdInt")
		return transferobject.Order{}, err
	}

	order, err := cart.Order(userId, mapUserAddressById)
	if err != nil {
		s.log.Error().
			Err(err).
			Msg("CreateOrder.cart.Order")
		return transferobject.Order{}, err
	}

	defer func() {
		// compensate the process
		if err != nil {
			// alerting failed to compensate request
		}
	}()

	order, err = s.orderAggregator.CreateOrder(ctx, order)
	if err != nil {
		s.log.Error().
			Err(err).
			Msg("CreateOrder.CreateOrder")
		return transferobject.Order{}, err
	}

	// Here the stock reserving
	// This code wasn't implemented using any distributed lock such as redis
	// because the locking is happened in the database level
	// mysql/ mariadb by nature has locking
	// so I just using the unsigned int to make sure the stock reduction won't make the minus
	reserveStocks := order.OrderReceipts.ReserveStocks()
	requestReserveStock := transferobject.NewReserveStocks(reserveStocks)
	respReserveStock, err := s.productSvc.AddReserveStock(ctx, transferobject.RequestReserveStoct{
		Data: requestReserveStock,
	})
	if err != nil {
		s.log.Error().
			Err(err).
			Interface("reserveStocks", respReserveStock).
			Msg("CreateOrder.productSvc.AddReserveStock")
		return transferobject.Order{}, err
	}

	err = s.cartSvc.DeleteCart(ctx, transferobject.RequestDeleteCart{
		UserId: request.UserId,
	})
	if err != nil {
		s.log.Error().
			Err(err).
			Msg("CreateOrder.DeleteCart")
		return transferobject.Order{}, err
	}

	return transferobject.NewOrder(order), nil
}

func (s *OrderService) ExpiredOrder(ctx context.Context) error {
	pagination := entity.Pagination{}
	pagination.Page = 1
	pagination.Size = 100

	s.log.Info().Msg("running cron job")

	for {
		expiredOrders, pagination, err := s.orderAggregator.GetOrders(ctx, entity.OrderFilter{
			TimeFrameInMinutes: 10,
			IsExpired:          null.BoolFrom(true),
			Pagination:         pagination,
		})
		if err != nil {
			s.log.Error().
				Err(err).
				Int("page", pagination.Page).
				Int("size", pagination.Size).
				Msg("ExpiredOrder.orderRepo.GetOrder")
			return err
		}
		if len(expiredOrders) == 0 {
			s.log.Debug().
				Int("page", pagination.Page).
				Int("size", pagination.Size).
				Msg("ExpiredOrder done")
			return nil
		}

		expiredOrders.SetAsExpired()

		err = s.orderRepoWriter.UpdateOrdersStatus(ctx, expiredOrders)
		if err != nil {
			return err
		}

		returnReserveStock := expiredOrders.ReturnReserveStocks()

		respAddReserveStock, err := s.productSvc.AddReserveStock(ctx, transferobject.RequestReserveStoct{
			Data: transferobject.NewReserveStocks(returnReserveStock),
		})
		if err != nil {
			s.log.Error().Err(err).
				Int("page", pagination.Page).
				Int("size", pagination.Size).
				Interface("resp reserve stock", respAddReserveStock).
				Msg("ExpiredOrder done")
			return err
		}

		pagination.Page++
	}

}
