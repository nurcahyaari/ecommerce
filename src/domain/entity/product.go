package entity

import (
	"fmt"
	"strings"

	"github.com/shopspring/decimal"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Product struct {
	Id           int64           `db:"id"`
	StoreId      int64           `db:"store_id"`
	WarehouseId  int64           `db:"warehouse_id"`
	Name         string          `db:"name"`
	Price        decimal.Decimal `db:"price"`
	ProductStock ProductStock    `db:"-"`
	Signature
}

func (p Product) ProductForCart(quantity uint) ProductForCart {
	return ProductForCart{
		ProductId: p.Id,
		Quantity:  quantity,
		Price:     p.Price,
	}
}

func (p *Product) MoveWarehouse(to Warehouse) {
	p.WarehouseId = to.Id
}

func (p Product) WarehouseIdStr() string {
	return fmt.Sprintf("%d", p.WarehouseId)
}

func (p Product) StoreIdStr() string {
	return fmt.Sprintf("%d", p.StoreId)
}

type Products []Product

func (ps Products) WarehouseIdsStr() string {
	warehouseIdsStr := []string{}
	for _, p := range ps {
		warehouseIdsStr = append(warehouseIdsStr, fmt.Sprintf("%d", p.WarehouseId))
	}
	return strings.Join(warehouseIdsStr, ",")
}

func (ps Products) AggregateStock(mapProductStockByProductId MapProductStock) {
	for i, p := range ps {
		productStock, ok := mapProductStockByProductId[p.Id]
		if !ok {
			continue
		}
		ps[i].ProductStock = productStock
	}
}

func (ps Products) ProductStockFilter() ProductStockFilter {
	psf := ProductStockFilter{}

	for _, p := range ps {
		psf.ProductIds = append(psf.ProductIds, p.Id)
	}

	return psf
}

// One returns the first warehouse on the list and the status (found or not found)
func (ps Products) One() (Product, bool) {
	if len(ps) == 0 {
		return Product{}, false
	}
	return ps[0], true
}

type ProductFilter struct {
	Or           *ProductFilter
	Ids          []int64
	StoreIds     []int64
	WarehouseIds []int64
	Pagination
}

func (f ProductFilter) composeFilter() ([]string, []interface{}) {
	var (
		query []string
		args  = make([]interface{}, 0)
	)

	if len(f.Ids) > 0 {
		query = append(query, "id IN (?)")
		args = append(args, f.Ids)
	}

	if len(f.StoreIds) > 0 {
		query = append(query, "store_id IN (?)")
		args = append(args, f.StoreIds)
	}

	if len(f.WarehouseIds) > 0 {
		query = append(query, "warehouse_id IN (?)")
		args = append(args, f.WarehouseIds)
	}

	return query, args
}

func (f ProductFilter) ComposeFilter() (string, []interface{}, error) {
	var (
		query []string
		args  = make([]interface{}, 0)
	)

	query, args = f.composeFilter()
	if f.Or != nil {
		queryOr, argsOr := f.Or.composeFilter()

		args = append(args, argsOr...)
		whereClauseOr := strings.Join(queryOr, " OR ")
		query = append(query, "("+whereClauseOr+")")
	}

	// Combine query parts into a single WHERE clause
	whereClause := strings.Join(query, " AND ")
	if whereClause != "" {
		whereClause = "WHERE " + whereClause
	}

	return whereClause, args, nil
}

type ReserveStock struct {
	Success bool   `db:"-"`
	Message string `db:"-"`
	ProductStock
}

type ReserveStocks []ReserveStock

type ProductStock struct {
	ProductId     int64 `db:"product_id"`
	StockReserved uint  `db:"stock_reserved"`
	StockOnHand   uint  `db:"stock_on_hand"`
	IsRevert      bool  `db:"-"`
}

// returns with this statement
// stockOnHand, stockReserved, id
func (ps ProductStock) ReserveStockArgs() []any {
	return []any{
		ps.StockReserved,
		ps.StockReserved,
		ps.ProductId,
	}
}

// it will reduce the reserve stock after the order was processed
func (ps ProductStock) UpdateStockArgs() []any {
	return []any{
		ps.StockReserved,
		ps.ProductId,
	}
}

type MapProductStock map[int64]ProductStock

type ProductStockFilter struct {
	ProductIds []int64
}

func (f ProductStockFilter) composeFilter() ([]string, []interface{}) {
	var (
		query []string
		args  = make([]interface{}, 0)
	)

	if len(f.ProductIds) > 0 {
		query = append(query, "product_id IN (?)")
		args = append(args, f.ProductIds)
	}

	return query, args
}

func (f ProductStockFilter) ComposeFilter() (string, []interface{}, error) {
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

type ProductStocks []ProductStock

func (pss ProductStocks) MapByProductId() MapProductStock {
	mapByProductId := make(MapProductStock)
	for _, ps := range pss {
		mapByProductId[ps.ProductId] = ps
	}
	return mapByProductId
}

type ProductForCart struct {
	ProductId int64
	Quantity  uint
	Price     decimal.Decimal
}

func (pfc ProductForCart) TotalPrice() decimal.Decimal {
	return pfc.Price.Mul(decimal.NewFromInt(int64(pfc.Quantity)))
}

func (pfc ProductForCart) CartItem(cartId string) (CartItem, error) {
	price, err := primitive.ParseDecimal128(pfc.Price.String())
	if err != nil {
		return CartItem{}, err
	}

	totalPrice, err := primitive.ParseDecimal128(pfc.TotalPrice().String())
	if err != nil {
		return CartItem{}, err
	}

	return CartItem{
		Id:              primitive.NewObjectID(),
		CartId:          cartId,
		ProductId:       pfc.ProductId,
		Quantity:        int32(pfc.Quantity),
		PricePerProduct: price,
		TotalPrice:      totalPrice,
	}, nil
}

type ProductsForCart []ProductForCart
