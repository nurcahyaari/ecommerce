package transferobject

import (
	"strconv"
	"strings"
	"time"

	"github.com/nurcahyaari/ecommerce/src/domain/entity"
	"gopkg.in/guregu/null.v4"
)

type Warehouse struct {
	Id        int64       `json:"id"`
	StoreId   int64       `json:"store_id"`
	Name      string      `json:"name"`
	Location  string      `json:"location"`
	IsActived bool        `json:"is_actived"`
	IsRemoved bool        `json:"is_removed"`
	CreatedAt time.Time   `json:"created_at"`
	CreatedBy null.String `json:"created_by"`
	UpdatedAt time.Time   `json:"updated_at"`
	UpdatedBy null.String `json:"updated_by"`
	DeletedAt null.Time   `json:"deleted_at"`
	DeletedBy null.String `json:"deleted_by"`
}

func (w Warehouse) Entity() entity.Warehouse {
	return entity.Warehouse{
		Id:       w.Id,
		StoreId:  w.StoreId,
		Name:     w.Name,
		Location: w.Location,
		Signature: entity.Signature{
			IsActived: w.IsActived,
			IsRemoved: w.IsRemoved,
			CreatedAt: w.CreatedAt,
			CreatedBy: w.CreatedBy,
			UpdatedAt: w.UpdatedAt,
			UpdatedBy: w.UpdatedBy,
			DeletedAt: w.DeletedAt,
			DeletedBy: w.DeletedBy,
		},
	}
}

func NewWarehouse(warehouse entity.Warehouse) Warehouse {
	return Warehouse{
		Id:        warehouse.Id,
		StoreId:   warehouse.StoreId,
		Name:      warehouse.Name,
		Location:  warehouse.Location,
		IsActived: warehouse.IsActived,
		IsRemoved: warehouse.IsRemoved,
		CreatedAt: warehouse.CreatedAt,
		CreatedBy: warehouse.CreatedBy,
		UpdatedAt: warehouse.UpdatedAt,
		UpdatedBy: warehouse.UpdatedBy,
		DeletedAt: warehouse.DeletedAt,
		DeletedBy: warehouse.DeletedBy,
	}
}

type Warehouses []Warehouse

func NewWarehouses(warehouses entity.Warehouses) Warehouses {
	respWarehouses := make(Warehouses, 0)
	for _, warehouse := range warehouses {
		respWarehouses = append(respWarehouses, NewWarehouse(warehouse))
	}

	return respWarehouses
}

type RequestSearchWarehouse struct {
	Ids      string
	StoreIds string
	Pagination
}

func (r RequestSearchWarehouse) composeWarehouseFilter() (entity.WarehouseFilter, error) {
	WarehouseFilter := entity.WarehouseFilter{}
	if r.Ids != "" {
		ids := strings.Split(r.Ids, ",")
		for _, id := range ids {
			idInt, err := strconv.ParseInt(id, 10, 64)
			if err != nil {
				return WarehouseFilter, err
			}

			WarehouseFilter.Ids = append(WarehouseFilter.Ids, idInt)
		}
	}

	if r.StoreIds != "" {
		ids := strings.Split(r.StoreIds, ",")
		for _, id := range ids {
			idInt, err := strconv.ParseInt(id, 10, 64)
			if err != nil {
				return WarehouseFilter, err
			}

			WarehouseFilter.StoreIds = append(WarehouseFilter.StoreIds, idInt)
		}
	}

	return WarehouseFilter, nil
}

func (r RequestSearchWarehouse) WarehouseFilter() (entity.WarehouseFilter, error) {
	WarehouseFilter, err := r.composeWarehouseFilter()
	if err != nil {
		return WarehouseFilter, err
	}

	if r.Pagination.Page == 0 && r.Pagination.Size == 0 {
		r.Pagination.Default()
	}

	if r.Pagination.Page != 0 {
		WarehouseFilter.Pagination.Page = r.Pagination.Page
	}

	if r.Pagination.Size != 0 {
		WarehouseFilter.Pagination.Size = r.Pagination.Size
	}

	return WarehouseFilter, nil
}

type ResponseSearchWarehouse struct {
	Data       Warehouses `json:"warehouses"`
	Pagination Pagination `json:"pagination"`
}

func NewResponseSearchWarehouse(warehouses entity.Warehouses, pagination entity.Pagination) ResponseSearchWarehouse {
	return ResponseSearchWarehouse{
		Data:       NewWarehouses(warehouses),
		Pagination: NewPagination(pagination),
	}
}

type ResponseGetWarehouse struct {
	Warehouse Warehouse `json:"warehouse"`
}

func NewResponseGetWarehouse(warehouse entity.Warehouse) ResponseGetWarehouse {
	return ResponseGetWarehouse{
		Warehouse: NewWarehouse(warehouse),
	}
}

type RequestOpenCloseWarehouse struct {
	Open  OpenCloseWarehouses `json:"open_warehouse"`
	Close OpenCloseWarehouses `json:"close_warehouse"`
}

func (rocw RequestOpenCloseWarehouse) WarehouseFilter() entity.WarehouseFilter {
	filter := entity.WarehouseFilter{}
	openWarehouseFilter := rocw.Open.WarehouseFilter()
	closeWarehouseFilter := rocw.Close.WarehouseFilter()

	filter.Ids = append(filter.Ids, openWarehouseFilter.Ids...)
	filter.Ids = append(filter.Ids, closeWarehouseFilter.Ids...)
	filter.StoreIds = append(filter.StoreIds, openWarehouseFilter.StoreIds...)
	filter.StoreIds = append(filter.StoreIds, closeWarehouseFilter.StoreIds...)

	filter.Page = 1
	filter.Size = len(filter.Ids)

	return filter
}

func (rocw RequestOpenCloseWarehouse) Warehouses() entity.Warehouses {
	open := rocw.Open.Warehouses(true)
	close := rocw.Close.Warehouses(false)

	ws := entity.Warehouses{}
	ws = append(ws, open...)
	ws = append(ws, close...)
	return ws
}

type OpenCloseWarehouse struct {
	WarehouseId int64 `json:"warehouse_id"`
	StoreId     int64 `json:"store_id"`
}

type OpenCloseWarehouses []OpenCloseWarehouse

func (ocws OpenCloseWarehouses) Warehouses(open bool) entity.Warehouses {
	ws := entity.Warehouses{}
	for _, ocw := range ocws {
		ws = append(ws, entity.Warehouse{
			Id:      ocw.WarehouseId,
			StoreId: ocw.StoreId,
			Signature: entity.Signature{
				IsActived: open,
			},
		})
	}

	return ws
}

func (ocws OpenCloseWarehouses) WarehouseFilter() entity.WarehouseFilter {
	filter := entity.WarehouseFilter{}

	for _, ocw := range ocws {
		filter.Ids = append(filter.Ids, ocw.StoreId)
		filter.StoreIds = append(filter.StoreIds, ocw.StoreId)
	}

	return filter
}
