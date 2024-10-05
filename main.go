package main

import (
	"context"

	"github.com/nurcahyaari/ecommerce/config"
	"github.com/nurcahyaari/ecommerce/infrastructure/database"
	"github.com/nurcahyaari/ecommerce/internal/graceful"
	"github.com/nurcahyaari/ecommerce/internal/logger"
	"github.com/nurcahyaari/ecommerce/internal/protocols/http"
	"github.com/nurcahyaari/ecommerce/internal/protocols/http/router"
	authsvc "github.com/nurcahyaari/ecommerce/src/application/auth/service"
	cartrepo "github.com/nurcahyaari/ecommerce/src/application/cart/repository"
	cartsvc "github.com/nurcahyaari/ecommerce/src/application/cart/service"
	productrepo "github.com/nurcahyaari/ecommerce/src/application/inventory/repository"
	productsvc "github.com/nurcahyaari/ecommerce/src/application/inventory/service"
	userrepo "github.com/nurcahyaari/ecommerce/src/application/user/repository"
	usersvc "github.com/nurcahyaari/ecommerce/src/application/user/service"
	httphandler "github.com/nurcahyaari/ecommerce/src/handlers/http"
)

func main() {
	// logger.InitLogger()
	logger := logger.NewLogger()
	config := config.Get()

	db := database.NewMysql(config)
	mongoDb, err := database.NewMongoDB(&config)
	if err != nil {
		logger.Fatal().Err(err)
	}

	userRepoReader := userrepo.NewUserRepositoryRead(db)
	userSvc := usersvc.NewUserService(&config, logger, userRepoReader)

	userAddressRepoReader := userrepo.NewUserAddressRepositoryRead(db)
	userAddressSvc := usersvc.NewUserAddressService(&config, logger, userAddressRepoReader)

	authSvc := authsvc.NewAuthService(&config, logger, nil, nil, userSvc)

	warehouseRepoReader := productrepo.NewWarehouseRepositoryRead(db)
	warehouseRepoWriter := productrepo.NewWarehouseRepositoryWrite(db)
	warehouseSvc := productsvc.NewWarehouseService(&config, logger, warehouseRepoReader, warehouseRepoWriter)

	productRepoReader := productrepo.NewProductRepositoryRead(db)
	productRepoWriter := productrepo.NewProductRepositoryWrite(db)
	productStockRepoReader := productrepo.NewProductStockRepositoryRead(db)
	productStockRepoWriter := productrepo.NewProductStockRepositoryWrite(db)
	productRepoAggregator := productrepo.NewProductAggregate(productRepoReader, productStockRepoReader)
	productSvc := productsvc.NewProductService(
		&config,
		logger,
		productRepoReader,
		productRepoWriter,
		productStockRepoWriter,
		productRepoAggregator,
		warehouseSvc)

	cartRepo := cartrepo.NewCartRepository(mongoDb)
	cartSvc := cartsvc.NewCartService(&config, logger, cartRepo, productSvc, userSvc, userAddressSvc)

	httpHandler := httphandler.NewHttpHandler(userSvc, authSvc, productSvc, warehouseSvc, cartSvc)

	httpRouter := router.NewHttpRouter(httpHandler)
	httpServer := http.New(config, httpRouter)
	// grpcHandler := grpchandler.NewGrpcHandler(&config, bookService)

	// grpcServer := grpc.New(&config, grpcHandler)

	// grpcServer.Listen()
	go httpServer.Listen()

	graceful.GracefulShutdown(context.Background(), graceful.RequestGraceful{
		WarnPeriod: 0,
		Operations: map[string]graceful.Operation{
			"httpServer": httpServer.Shutdown,
		},
	})
}
