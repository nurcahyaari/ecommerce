package service

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/nurcahyaari/ecommerce/config"
	internalcontext "github.com/nurcahyaari/ecommerce/internal/x/context"
	internalerrors "github.com/nurcahyaari/ecommerce/internal/x/errors"
	"github.com/nurcahyaari/ecommerce/src/domain/entity"
	"github.com/nurcahyaari/ecommerce/src/domain/repository"
	"github.com/nurcahyaari/ecommerce/src/domain/service"
	"github.com/nurcahyaari/ecommerce/src/transferobject"
	"github.com/rs/zerolog"
	"golang.org/x/sync/errgroup"
	"gopkg.in/guregu/null.v4"
)

type CartService struct {
	cfg            *config.Config
	log            zerolog.Logger
	cartRepo       repository.CartRepositorier
	productSvc     service.ProductServicer
	userSvc        service.UserServicer
	userAddressSvc service.UserAddressServicer
}

func NewCartService(
	cfg *config.Config,
	log zerolog.Logger,
	cartRepo repository.CartRepositorier,
	productSvc service.ProductServicer,
	userSvc service.UserServicer,
	userAddressSvc service.UserAddressServicer,
) service.CartServicer {
	return &CartService{
		cfg:            cfg,
		log:            log,
		cartRepo:       cartRepo,
		productSvc:     productSvc,
		userSvc:        userSvc,
		userAddressSvc: userAddressSvc,
	}
}

func (s *CartService) GetCart(ctx context.Context, request transferobject.RequestGetCart) (transferobject.ResponseGetCart, error) {
	filter, err := request.CartFilter()
	if err != nil {
		s.log.Error().
			Err(err).
			Msg("GetCart.CartFilter")
		return transferobject.ResponseGetCart{}, err
	}

	carts, err := s.cartRepo.FindCart(ctx, filter)
	if err != nil {
		s.log.Error().
			Err(err).
			Msg("GetCart.FindCart")
		return transferobject.ResponseGetCart{}, err
	}

	resp, err := transferobject.NewResponseGetCart(carts)
	if err != nil {
		s.log.Error().
			Err(err).
			Msg("GetCart.NewResponseGetCart")
		return transferobject.ResponseGetCart{}, err
	}

	return resp, nil
}

func (s *CartService) AddItemToCart(ctx context.Context, request transferobject.RequestAddItemToCart) (transferobject.Cart, error) {
	routine, errGroupCtx := errgroup.WithContext(ctx)

	var (
		respUserAddress transferobject.ResponseGetUserAddress
		respProduct     transferobject.ResponseGetProduct
		userId, err     = internalcontext.GetUserIdInt64(ctx)
	)
	if err != nil {
		s.log.Error().
			Err(err).
			Msg("AddItemToCart.GetUserIdInt64")
		return transferobject.Cart{}, err
	}

	routine.Go(func() error {
		ua, err := s.userAddressSvc.GetUserAddress(errGroupCtx, transferobject.RequestSearchUserAddress{
			Ids: fmt.Sprintf("%d", request.UserAddressId),
		})
		if err != nil {
			s.log.Error().
				Err(err).
				Msg("AddItemToCart.GetUserAddress")
			return err
		}
		respUserAddress = ua
		return nil
	})

	routine.Go(func() error {
		p, err := s.productSvc.GetProduct(errGroupCtx, transferobject.RequestSearchProduct{
			Ids: fmt.Sprintf("%d", request.ProductId),
		})
		if err != nil {
			s.log.Error().
				Err(err).
				Msg("AddItemToCart.GetProduct")
			return err
		}

		respProduct = p
		return nil
	})

	if err := routine.Wait(); err != nil {
		s.log.Error().
			Err(err).
			Msg("AddItemToCart.ErrGroup")
		return transferobject.Cart{}, err
	}

	if respUserAddress.UserAddress.Id == 0 {
		s.log.Warn().
			Msg("AddItemToCart userAddress is not found")
		return transferobject.Cart{}, internalerrors.New(
			errors.New("err: user address is not found"),
			internalerrors.SetErrorCode(http.StatusNotFound))
	}

	if respProduct.Product.Id == 0 {
		s.log.Warn().
			Msg("AddItemToCart product is not found")
		return transferobject.Cart{}, internalerrors.New(
			errors.New("err: user address is not found"),
			internalerrors.SetErrorCode(http.StatusNotFound))
	}

	carts, err := s.cartRepo.FindCart(ctx, entity.CartFilter{
		UserId:        null.IntFrom(userId),
		UserAddressId: null.IntFrom(request.UserAddressId),
	})
	if err != nil {
		s.log.Error().Err(err).
			Msg("AddItemToCart")
		return transferobject.Cart{}, err
	}

	cart, ok := carts.One()
	product := respProduct.Product.Entity()
	productForCart := product.ProductForCart(request.Quantity)
	if !ok {
		// create a new cart
		cart, err := entity.NewCart(respUserAddress.UserAddress.Entity(), productForCart)
		if err != nil {
			s.log.Error().Err(err).
				Str("action", "insert").
				Msg("AddItemToCart.NewCart")
			return transferobject.Cart{}, err
		}
		err = s.cartRepo.UpsertCart(ctx, cart)
		if err != nil {
			s.log.Error().Err(err).
				Str("action", "insert").
				Msg("AddItemToCart failed to insert data")
			return transferobject.Cart{}, err
		}

		cartResp, err := transferobject.NewCart(cart)
		if err != nil {
			s.log.Error().Err(err).
				Str("action", "insert").
				Msg("AddItemToCart failed to insert data")
			return transferobject.Cart{}, err
		}
		return cartResp, nil
	}

	// update existing cart
	cart.UpdateCartItems(productForCart)

	err = s.cartRepo.UpsertCart(ctx, cart)
	if err != nil {
		s.log.Error().Err(err).
			Str("action", "create").
			Msg("AddItemToCart failed to upsert data")
		return transferobject.Cart{}, err
	}
	cartResp, err := transferobject.NewCart(cart)
	if err != nil {
		s.log.Error().Err(err).
			Str("action", "insert").
			Msg("AddItemToCart failed to upsert data")
		return transferobject.Cart{}, err
	}

	return cartResp, nil
}
