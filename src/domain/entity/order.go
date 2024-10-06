package entity

import (
	"fmt"
	"strings"
	"time"

	"github.com/shopspring/decimal"
	"gopkg.in/guregu/null.v4"
)

type OrderStatus int

func (w OrderStatus) String() string {
	return [...]string{"Pending", "Paid", "Processed", "Received", "Done", "Refund", "Expired"}[w-1]
}

func (w OrderStatus) EnumIndex() int {
	return int(w)
}

const (
	Pending OrderStatus = iota + 1
	Paid
	Processed
	Received
	Done
	Refund
	Expired
)

type Order struct {
	Id            int64           `db:"id"`
	UserId        int64           `db:"user_id"`
	OrderCode     string          `db:"order_code"`
	TotalQuantity int32           `db:"total_quantity"`
	TotalPrice    decimal.Decimal `db:"total_price"`
	OrderStatus   OrderStatus     `db:"order_status"`
	ExpiredOrder  null.Time       `db:"expired_order"`
	OrderReceipts OrderReceipts   `db:"-"`
	Signature
}

func (o *Order) SetOrderId(orderId int64) {
	o.Id = orderId
	for i, or := range o.OrderReceipts {
		o.OrderReceipts[i].OrderId = orderId
		o.OrderReceipts[i].OrderAddress.OrderId = orderId
		for j, _ := range or.OrderDetails {
			o.OrderReceipts[i].OrderDetails[j].OrderId = orderId
		}
	}
}

func (o Order) OrderAddresses() OrderAddresses {
	orderAddresses := OrderAddresses{}

	for _, or := range o.OrderReceipts {
		orderAddresses = append(orderAddresses, or.OrderAddress)
	}

	return orderAddresses
}

func (o *Order) SumTotalQuantity() int32 {
	totalQuantity := int32(0)
	for _, or := range o.OrderReceipts {
		totalQuantity += or.TotalQuantity
	}
	o.TotalQuantity = totalQuantity
	return totalQuantity
}

func (o *Order) SumTotalPrice() decimal.Decimal {
	totalPrice := decimal.Decimal{}
	for _, or := range o.OrderReceipts {
		totalPrice = totalPrice.Add(or.TotalPrice)
	}
	o.TotalPrice = totalPrice
	return totalPrice
}

func (o *Order) GenerateOrderCode() {
	o.OrderCode = fmt.Sprintf("%d-%d", o.UserId, time.Now().Nanosecond())
}

type Orders []Order

func (os Orders) SetOrderReceipts(mapOrderReceiptByOrderId MapOrderReceipts) {
	for i, o := range os {
		orderReceipts, ok := mapOrderReceiptByOrderId[o.Id]
		if !ok {
			continue
		}

		os[i].OrderReceipts = orderReceipts
	}
}

func (os Orders) Ids() []int64 {
	ids := make([]int64, 0)
	for _, o := range os {
		ids = append(ids, o.Id)
	}
	return ids
}

func (os Orders) SetAsExpired() {
	for i, _ := range os {
		os[i].OrderStatus = Expired
	}
}

func (os Orders) ReturnReserveStocks() ReserveStocks {
	rrs := ReserveStocks{}

	for _, o := range os {
		for _, or := range o.OrderReceipts {
			for _, od := range or.OrderDetails {
				rrs = append(rrs, ReserveStock{
					ProductStock: ProductStock{
						ProductId:     od.ProdutId,
						StockReserved: uint(od.Quantity),
						IsRevert:      true,
					},
				})
			}
		}
	}

	return rrs
}

type OrderFilter struct {
	TimeFrameInMinutes int64
	IsExpired          null.Bool
	Pagination
}

func (f OrderFilter) composeFilter() ([]string, []interface{}) {
	var (
		query []string
		args  = make([]interface{}, 0)
	)

	if f.IsExpired.Valid {
		query = append(query, "order_status = ?")
		args = append(args, Pending)

		query = append(query, "order_status != ?")
		args = append(args, Expired)

		query = append(query, fmt.Sprintf("NOW() - INTERVAL %d MINUTE >= expired_order", f.TimeFrameInMinutes))
	}

	return query, args
}

func (f OrderFilter) ComposeFilter() (string, []interface{}, error) {
	var (
		query []string
		args  = make([]interface{}, 0)
	)

	query, args = f.composeFilter()

	// Combine query parts into a single WHERE clause
	whereClause := strings.Join(query, " AND ")
	if whereClause != "" {
		whereClause = "WHERE " + whereClause
	}

	return whereClause, args, nil
}

type OrderReceipt struct {
	Id             int64           `db:"id"`
	OrderId        int64           `db:"order_id"`
	OrderAddressId int64           `db:"order_address_id"`
	TotalQuantity  int32           `db:"total_quantity"`
	TotalPrice     decimal.Decimal `db:"total_price"`
	OrderDetails   OrderDetails    `db:"-"`
	OrderAddress   OrderAddress    `db:"-"`
	Signature
}

type OrderReceipts []OrderReceipt

func (ors OrderReceipts) ReserveStocks() ReserveStocks {
	reserveStocks := ReserveStocks{}

	for _, or := range ors {
		reserveStocks = append(reserveStocks, or.OrderDetails.ReserveStocks()...)
	}

	return reserveStocks
}

func (ors OrderReceipts) SetOrderAddressId(mapOrderAddressByUserAddressId MapOrderAddress) {
	for i, or := range ors {
		orderAddress, ok := mapOrderAddressByUserAddressId[or.OrderAddress.UserAddressId]
		if !ok {
			continue
		}

		ors[i].OrderAddressId = orderAddress.Id
	}
}

func (ors OrderReceipts) SetOrderDetail(mapOrderDetailByOrderReceiptId MapOrderDetails) {
	for i, or := range ors {
		orderDetails, ok := mapOrderDetailByOrderReceiptId[or.Id]
		if !ok {
			continue
		}

		ors[i].OrderDetails = orderDetails
	}
}

func (ors OrderReceipts) SetOrderReceiptIdToDetail() {
	for i, or := range ors {
		for j, _ := range or.OrderDetails {
			ors[i].OrderDetails[j].OrderReceiptId = or.Id
		}
	}
}

func (ors OrderReceipts) OrderDetails() OrderDetails {
	orderDetails := OrderDetails{}

	for _, or := range ors {
		orderDetails = append(orderDetails, or.OrderDetails...)
	}

	return orderDetails
}

func (ors OrderReceipts) MapOrderReceiptsByOrderId() MapOrderReceipts {
	mapByOrderId := make(MapOrderReceipts)

	for _, or := range ors {
		mapByOrderId[or.OrderId] = append(mapByOrderId[or.OrderId], or)
	}

	return mapByOrderId
}

type MapOrderReceipt map[int64]OrderReceipt

type MapOrderReceipts map[int64]OrderReceipts

type OrderReceiptFilter struct {
	OrderIds   []int64
	Pagination *Pagination
}

func (f OrderReceiptFilter) composeFilter() ([]string, []interface{}) {
	var (
		query []string
		args  = make([]interface{}, 0)
	)

	if len(f.OrderIds) > 0 {
		query = append(query, "order_id IN (?)")
		args = append(args, f.OrderIds)
	}

	return query, args
}

func (f OrderReceiptFilter) ComposeFilter() (string, []interface{}, error) {
	var (
		query []string
		args  = make([]interface{}, 0)
	)

	query, args = f.composeFilter()

	// Combine query parts into a single WHERE clause
	whereClause := strings.Join(query, " AND ")
	if whereClause != "" {
		whereClause = "WHERE " + whereClause
	}

	return whereClause, args, nil
}

type OrderDetail struct {
	Id              int64           `db:"id"`
	OrderId         int64           `db:"order_id"`
	OrderReceiptId  int64           `db:"order_receipt_id"`
	ProdutId        int64           `db:"product_id"`
	Quantity        int32           `db:"quantity"`
	PricePerProduct decimal.Decimal `db:"product_price"`
	TotalPrice      decimal.Decimal `db:"total_price"`
	Signature
}

func (od OrderDetail) ReserveStock() ReserveStock {
	return ReserveStock{
		ProductStock: ProductStock{
			ProductId:     od.ProdutId,
			StockReserved: uint(od.Quantity),
		},
	}
}

type OrderDetails []OrderDetail

func (ods OrderDetails) ReserveStocks() ReserveStocks {
	reserveStocks := ReserveStocks{}
	for _, od := range ods {
		reserveStocks = append(reserveStocks, od.ReserveStock())
	}

	return reserveStocks
}

func (ods OrderDetails) MapOrderDetailsByOrderReceiptId() MapOrderDetails {
	mapByOrderReceiptId := make(MapOrderDetails)

	for _, od := range ods {
		mapByOrderReceiptId[od.OrderReceiptId] = append(mapByOrderReceiptId[od.OrderReceiptId], od)
	}

	return mapByOrderReceiptId
}

type MapOrderDetails map[int64]OrderDetails

type OrderDetailFilter struct {
	OrderIds   []int64
	Pagination *Pagination
}

func (f OrderDetailFilter) composeFilter() ([]string, []interface{}) {
	var (
		query []string
		args  = make([]interface{}, 0)
	)

	if len(f.OrderIds) > 0 {
		query = append(query, "order_id IN (?)")
		args = append(args, f.OrderIds)
	}

	return query, args
}

func (f OrderDetailFilter) ComposeFilter() (string, []interface{}, error) {
	var (
		query []string
		args  = make([]interface{}, 0)
	)

	query, args = f.composeFilter()

	// Combine query parts into a single WHERE clause
	whereClause := strings.Join(query, " AND ")
	if whereClause != "" {
		whereClause = "WHERE " + whereClause
	}

	return whereClause, args, nil
}

type OrderAddress struct {
	Id            int64  `db:"id"`
	OrderId       int64  `db:"order_id"`
	UserId        int64  `db:"user_id"`
	UserAddressId int64  `db:"-"`
	FullAddress   string `db:"full_address"`
	Signature
}

type OrderAddresses []OrderAddress

func (oas OrderAddresses) MapById() MapOrderAddress {
	mapById := make(MapOrderAddress)

	for _, oa := range oas {
		mapById[oa.Id] = oa
	}

	return mapById
}

func (oas OrderAddresses) MapByUserAddressId() MapOrderAddress {
	mapById := make(MapOrderAddress)

	for _, oa := range oas {
		mapById[oa.UserAddressId] = oa
	}

	return mapById
}

type MapOrderAddress map[int64]OrderAddress
