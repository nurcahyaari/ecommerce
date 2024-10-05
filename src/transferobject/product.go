package transferobject

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/nurcahyaari/ecommerce/src/domain/entity"
	"github.com/shopspring/decimal"
)

type RequestSearchProduct struct {
	Ids          string
	WarehouseIds string
	StoreIds     string
	Pagination
}

func (r RequestSearchProduct) composeUserFilter() (entity.ProductFilter, error) {
	productFilter := entity.ProductFilter{}
	if r.Ids != "" {
		ids := strings.Split(r.Ids, ",")
		for _, id := range ids {
			idInt, err := strconv.ParseInt(id, 10, 64)
			if err != nil {
				return productFilter, err
			}

			productFilter.Ids = append(productFilter.Ids, idInt)
		}
	}

	if r.WarehouseIds != "" {
		ids := strings.Split(r.WarehouseIds, ",")
		for _, id := range ids {
			idInt, err := strconv.ParseInt(id, 10, 64)
			if err != nil {
				return productFilter, err
			}

			productFilter.WarehouseIds = append(productFilter.WarehouseIds, idInt)
		}
	}

	if r.StoreIds != "" {
		ids := strings.Split(r.StoreIds, ",")
		for _, id := range ids {
			idInt, err := strconv.ParseInt(id, 10, 64)
			if err != nil {
				return productFilter, err
			}

			productFilter.StoreIds = append(productFilter.StoreIds, idInt)
		}
	}

	return productFilter, nil
}

func (r RequestSearchProduct) ProductFilter() (entity.ProductFilter, error) {
	userFilter, err := r.composeUserFilter()
	if err != nil {
		return userFilter, err
	}

	if r.Pagination.Page == 0 && r.Pagination.Size == 0 {
		r.Pagination.Default()
	}

	if r.Pagination.Page != 0 {
		userFilter.Pagination.Page = r.Pagination.Page
	}

	if r.Pagination.Size != 0 {
		userFilter.Pagination.Size = r.Pagination.Size
	}

	return userFilter, nil
}

type Product struct {
	Id            int64           `json:"id"`
	StoreId       int64           `json:"store_id"`
	WarehouseId   int64           `json:"warehouse_id"`
	Name          string          `json:"name"`
	Price         decimal.Decimal `json:"price"`
	StockReserved uint            `json:"stock_reserved"`
	StockOnHand   uint            `json:"stock_on_hand"`
}

func (p Product) Entity() entity.Product {
	return entity.Product{
		Id:          p.Id,
		StoreId:     p.StoreId,
		WarehouseId: p.WarehouseId,
		Name:        p.Name,
		Price:       p.Price,
		ProductStock: entity.ProductStock{
			ProductId:     p.Id,
			StockReserved: p.StockReserved,
			StockOnHand:   p.StockOnHand,
		},
	}
}

func NewProduct(product entity.Product) Product {
	return Product{
		Id:            product.Id,
		StoreId:       product.StoreId,
		WarehouseId:   product.WarehouseId,
		Name:          product.Name,
		Price:         product.Price,
		StockReserved: product.ProductStock.StockReserved,
		StockOnHand:   product.ProductStock.StockOnHand,
	}
}

type Products []Product

func NewProducts(products entity.Products) Products {
	respProducts := make(Products, 0)
	for _, product := range products {
		respProducts = append(respProducts, NewProduct(product))
	}

	return respProducts
}

type ResponseSearchProduct struct {
	Data       Products   `json:"products"`
	Pagination Pagination `json:"pagination"`
}

func NewResponseSearchProduct(products entity.Products, pagination entity.Pagination) ResponseSearchProduct {
	return ResponseSearchProduct{
		Data:       NewProducts(products),
		Pagination: NewPagination(pagination),
	}
}

type RequestMoveWarehouse struct {
	ProductId         int64 `json:"-"`
	WarehouseTargetId int64 `json:"warehouse_target_id"`
}

func (req RequestMoveWarehouse) WarehouseTargetIdStr() string {
	return fmt.Sprintf("%d", req.WarehouseTargetId)
}

func (req *RequestMoveWarehouse) ParseURLParams(r *http.Request) error {
	productId, err := strconv.ParseInt(chi.URLParam(r, "productId"), 10, 64)
	if err != nil {
		return err
	}

	req.ProductId = productId
	return nil
}

func (p RequestMoveWarehouse) ProductFilter() entity.ProductFilter {
	productFilter := entity.ProductFilter{
		Ids: []int64{p.ProductId},
	}

	productFilter.DefaultPagination()
	return productFilter
}

type ResponseGetProduct struct {
	Product Product
}

func NewResponseGetProduct(product entity.Product) ResponseGetProduct {
	return ResponseGetProduct{
		Product: NewProduct(product),
	}
}
