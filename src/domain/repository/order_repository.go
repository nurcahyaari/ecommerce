package repository

import (
	"context"

	"github.com/nurcahyaari/ecommerce/infrastructure/database"
	"github.com/nurcahyaari/ecommerce/src/domain/entity"
)

type OrderRepositoryReader interface {
	GetOrder(ctx context.Context, filter entity.OrderFilter) (entity.Orders, entity.Pagination, error)
}

type OrderRepositoryWriter interface {
	database.SQLDatabaseTrx[OrderRepositoryWriter]
	CreateOrder(ctx context.Context, data entity.Order) (entity.Order, error)
	UpdateOrderStatus(ctx context.Context, order entity.Order) error
	UpdateOrdersStatus(ctx context.Context, orders entity.Orders) error
}

type OrderReceiptRepositoryReader interface {
	GetOrderReceipts(ctx context.Context, flter entity.OrderReceiptFilter) (entity.OrderReceipts, entity.Pagination, error)
}

type OrderReceiptRepositoryWriter interface {
	database.SQLDatabaseTrx[OrderReceiptRepositoryWriter]
	CreateOrderReceipts(ctx context.Context, data entity.OrderReceipts) (entity.OrderReceipts, error)
}

type OrderDetailRepositoryReader interface {
	GetOrderDetails(ctx context.Context, flter entity.OrderDetailFilter) (entity.OrderDetails, entity.Pagination, error)
}

type OrderDetailRepositoryWriter interface {
	database.SQLDatabaseTrx[OrderDetailRepositoryWriter]
	CreateOrderDetails(ctx context.Context, data entity.OrderDetails) (entity.OrderDetails, error)
}

type OrderAddressRepositoryWriter interface {
	database.SQLDatabaseTrx[OrderAddressRepositoryWriter]
	CreateOrderAddresses(ctx context.Context, data entity.OrderAddresses) (entity.OrderAddresses, error)
}

type OrderAggregator interface {
	GetOrders(ctx context.Context, filter entity.OrderFilter) (entity.Orders, entity.Pagination, error)
	CreateOrder(ctx context.Context, order entity.Order) (entity.Order, error)
}
