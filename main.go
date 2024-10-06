package main

import (
	"context"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/nurcahyaari/ecommerce/config"
	"github.com/nurcahyaari/ecommerce/infrastructure/database"
	"github.com/nurcahyaari/ecommerce/internal/graceful"
	"github.com/nurcahyaari/ecommerce/internal/logger"
	"github.com/nurcahyaari/ecommerce/internal/protocols/cron"
	"github.com/nurcahyaari/ecommerce/internal/protocols/http"
	"github.com/nurcahyaari/ecommerce/internal/protocols/http/router"
	authsvc "github.com/nurcahyaari/ecommerce/src/application/auth/service"
	cartrepo "github.com/nurcahyaari/ecommerce/src/application/cart/repository"
	cartsvc "github.com/nurcahyaari/ecommerce/src/application/cart/service"
	productrepo "github.com/nurcahyaari/ecommerce/src/application/inventory/repository"
	productsvc "github.com/nurcahyaari/ecommerce/src/application/inventory/service"
	orderrepo "github.com/nurcahyaari/ecommerce/src/application/order/repository"
	ordersvc "github.com/nurcahyaari/ecommerce/src/application/order/service"
	userrepo "github.com/nurcahyaari/ecommerce/src/application/user/repository"
	usersvc "github.com/nurcahyaari/ecommerce/src/application/user/service"
	cronhandler "github.com/nurcahyaari/ecommerce/src/handlers/cron"
	httphandler "github.com/nurcahyaari/ecommerce/src/handlers/http"
)

func main() {
	logger := logger.NewLogger()
	cfg := config.Get()

	db := database.NewMysql(cfg)
	if cfg.DB.MySQL.WithMigration {
		err := db.RunMigrations(cfg)
		if err != nil {
			logger.Fatal().Err(err).
				Msg("failure to migrate")
		}
	}

	mongoDb, err := database.NewMongoDB(&cfg)
	if err != nil {
		logger.Fatal().Err(err)
	}

	userRepoReader := userrepo.NewUserRepositoryRead(db)
	userSvc := usersvc.NewUserService(&cfg, logger, userRepoReader)

	userAddressRepoReader := userrepo.NewUserAddressRepositoryRead(db)
	userAddressSvc := usersvc.NewUserAddressService(&cfg, logger, userAddressRepoReader)

	authSvc := authsvc.NewAuthService(&cfg, logger, nil, nil, userSvc)

	warehouseRepoReader := productrepo.NewWarehouseRepositoryRead(db)
	warehouseRepoWriter := productrepo.NewWarehouseRepositoryWrite(db)
	warehouseSvc := productsvc.NewWarehouseService(&cfg, logger, warehouseRepoReader, warehouseRepoWriter)

	productRepoReader := productrepo.NewProductRepositoryRead(db)
	productRepoWriter := productrepo.NewProductRepositoryWrite(db)
	productStockRepoReader := productrepo.NewProductStockRepositoryRead(db)
	productStockRepoWriter := productrepo.NewProductStockRepositoryWrite(db)
	productRepoAggregator := productrepo.NewProductAggregate(productRepoReader, productStockRepoReader)
	productSvc := productsvc.NewProductService(
		&cfg,
		logger,
		productRepoReader,
		productRepoWriter,
		productStockRepoWriter,
		productRepoAggregator,
		warehouseSvc)

	cartRepo := cartrepo.NewCartRepository(mongoDb)
	cartSvc := cartsvc.NewCartService(&cfg, logger, cartRepo, productSvc, userSvc, userAddressSvc)

	orderAddressRepoWrite := orderrepo.NewOrderAddressRepositoryWrite(db)
	orderDetailRepoWrite := orderrepo.NewOrderDetailRepositoryWrite(db)
	orderDetailRepoRead := orderrepo.NewOrderDetailRepositoryRead(db)
	orderReeiptRepoWrite := orderrepo.NewOrderReceiptRepositoryWrite(db)
	orderReceiptRepoRead := orderrepo.NewOrderReceiptRepositoryRead(db)
	orderRepoWrite := orderrepo.NewOrderRepositoryWrite(db)
	orderRepoReader := orderrepo.NewOrderRepositoryRead(db)
	orderAggregatorRepo := orderrepo.NewOrderAggregate(
		db,
		orderRepoWrite,
		orderReeiptRepoWrite,
		orderDetailRepoWrite,
		orderAddressRepoWrite,
		orderRepoReader,
		orderReceiptRepoRead,
		orderDetailRepoRead,
	)
	orderService := ordersvc.NewOrderService(&cfg, logger, orderAggregatorRepo, userAddressSvc, cartSvc, productSvc, orderRepoReader, orderRepoWrite)

	httpHandler := httphandler.NewHttpHandler(logger, userSvc, authSvc, productSvc, warehouseSvc, cartSvc, orderService)

	httpRouter := router.NewHttpRouter(httpHandler)
	httpServer := http.New(cfg, httpRouter)
	// grpcHandler := grpchandler.NewGrpcHandler(&config, bookService)

	// grpcServer := grpc.New(&config, grpcHandler)

	cronnHandler := cronhandler.NewCronhandler(orderService)
	cron := cron.New(cfg, cronnHandler, logger)

	// grpcServer.Listen()
	go httpServer.Listen()
	go cron.Listen()

	graceful.GracefulShutdown(context.Background(), graceful.RequestGraceful{
		WarnPeriod: 0,
		Operations: map[string]graceful.Operation{
			"httpServer": httpServer.Shutdown,
			"cron":       cron.Shutdown,
		},
	})
}
