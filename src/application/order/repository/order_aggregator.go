package repository

import (
	"context"

	"github.com/nurcahyaari/ecommerce/infrastructure/database"
	"github.com/nurcahyaari/ecommerce/src/domain/entity"
	"github.com/nurcahyaari/ecommerce/src/domain/repository"
)

type OrderAggregate struct {
	db                           *database.SQLDatabase
	orderRepositoryWriter        repository.OrderRepositoryWriter
	orderReceiptRepositoryWriter repository.OrderReceiptRepositoryWriter
	orderDetailRepositoryWriter  repository.OrderDetailRepositoryWriter
	orderAddressRpositoryWriter  repository.OrderAddressRepositoryWriter
	orderRepositoryReader        repository.OrderRepositoryReader
	orderReceiptRepositoryReader repository.OrderReceiptRepositoryReader
	orderDetailRepositoryReader  repository.OrderDetailRepositoryReader
}

func NewOrderAggregate(
	db *database.SQLDatabase,
	orderRepositoryWriter repository.OrderRepositoryWriter,
	orderReceiptRepositoryWriter repository.OrderReceiptRepositoryWriter,
	orderDetailRepositoryWriter repository.OrderDetailRepositoryWriter,
	orderAddressRpositoryWriter repository.OrderAddressRepositoryWriter,
	orderRepositoryReader repository.OrderRepositoryReader,
	orderReceiptRepositoryReader repository.OrderReceiptRepositoryReader,
	orderDetailRepositoryReader repository.OrderDetailRepositoryReader,
) repository.OrderAggregator {
	return &OrderAggregate{
		db:                           db,
		orderRepositoryWriter:        orderRepositoryWriter,
		orderReceiptRepositoryWriter: orderReceiptRepositoryWriter,
		orderDetailRepositoryWriter:  orderDetailRepositoryWriter,
		orderAddressRpositoryWriter:  orderAddressRpositoryWriter,
		orderRepositoryReader:        orderRepositoryReader,
		orderReceiptRepositoryReader: orderReceiptRepositoryReader,
		orderDetailRepositoryReader:  orderDetailRepositoryReader,
	}
}

func (r *OrderAggregate) BeginTx(ctx context.Context) (*OrderAggregate, error) {
	tx, err := r.db.DB.BeginTxx(ctx, nil)
	if err != nil {
		return nil, err
	}

	db := &database.SQLDatabase{
		DB: r.db.DB,
		Tx: tx,
	}

	return &OrderAggregate{
		db:                           db,
		orderRepositoryWriter:        NewOrderRepositoryWrite(db),
		orderReceiptRepositoryWriter: NewOrderReceiptRepositoryWrite(db),
		orderDetailRepositoryWriter:  NewOrderDetailRepositoryWrite(db),
		orderAddressRpositoryWriter:  NewOrderAddressRepositoryWrite(db),
	}, nil
}

func (r *OrderAggregate) Commit(ctx context.Context) error {
	if r.db.Tx == nil {
		return nil
	}

	if err := r.db.Tx.Commit(); err != nil {
		return err
	}

	r.db.Tx = nil

	return nil
}

func (r *OrderAggregate) Rollback(ctx context.Context) error {
	if r.db.Tx == nil {
		return nil
	}

	if err := r.db.Tx.Rollback(); err != nil {
		return err
	}

	r.db.Tx = nil

	return nil
}

func (r *OrderAggregate) GetOrders(ctx context.Context, filter entity.OrderFilter) (entity.Orders, entity.Pagination, error) {
	orders, pagination, err := r.orderRepositoryReader.GetOrder(ctx, filter)
	if err != nil {
		return orders, pagination, err
	}

	orderIds := orders.Ids()
	orderReceipts, _, err := r.orderReceiptRepositoryReader.GetOrderReceipts(ctx, entity.OrderReceiptFilter{
		OrderIds: orderIds,
	})
	if err != nil {
		return orders, pagination, err
	}

	orderDetails, _, err := r.orderDetailRepositoryReader.GetOrderDetails(ctx, entity.OrderDetailFilter{
		OrderIds: orderIds,
	})
	if err != nil {
		return orders, pagination, err
	}

	mapByOrderReceipts := orderDetails.MapOrderDetailsByOrderReceiptId()
	orderReceipts.SetOrderDetail(mapByOrderReceipts)

	mapByOrderId := orderReceipts.MapOrderReceiptsByOrderId()
	orders.SetOrderReceipts(mapByOrderId)

	return orders, pagination, nil
}

func (r *OrderAggregate) CreateOrder(ctx context.Context, order entity.Order) (entity.Order, error) {

	var (
		orderAggregator *OrderAggregate
		respOrder       entity.Order
		orderReceipts   entity.OrderReceipts
		orderDetails    entity.OrderDetails
		err             error
	)
	orderAggregator, err = r.BeginTx(ctx)
	if err != nil {
		return entity.Order{}, err
	}

	defer func() {
		if p := recover(); p != nil {
			orderAggregator.Rollback(ctx)
			panic(p)
		} else if err != nil {
			orderAggregator.Rollback(ctx)
		}
	}()

	respOrder, err = orderAggregator.orderRepositoryWriter.CreateOrder(ctx, order)
	if err != nil {
		return entity.Order{}, err
	}

	order.SetOrderId(respOrder.Id)

	orderAddresses, err := orderAggregator.orderAddressRpositoryWriter.CreateOrderAddresses(ctx, order.OrderAddresses())
	if err != nil {
		return entity.Order{}, err
	}

	order.OrderReceipts.SetOrderAddressId(orderAddresses.MapByUserAddressId())

	orderReceipts, err = orderAggregator.orderReceiptRepositoryWriter.CreateOrderReceipts(ctx, order.OrderReceipts)
	if err != nil {
		return entity.Order{}, err
	}
	order.OrderReceipts = orderReceipts
	order.OrderReceipts.SetOrderReceiptIdToDetail()
	orderDetails = orderReceipts.OrderDetails()
	orderDetails, err = orderAggregator.orderDetailRepositoryWriter.CreateOrderDetails(ctx, orderDetails)
	if err != nil {
		return entity.Order{}, err
	}

	order.OrderReceipts.SetOrderDetail(orderDetails.MapOrderDetailsByOrderReceiptId())

	err = orderAggregator.Commit(ctx)
	if err != nil {
		return entity.Order{}, err
	}

	return respOrder, nil
}
