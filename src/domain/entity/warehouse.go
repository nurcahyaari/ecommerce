package entity

import "strings"

type Warehouse struct {
	Id       int64  `db:"id"`
	StoreId  int64  `db:"store_id"`
	Name     string `db:"name"`
	Location string `db:"location"`
	Signature
}

type MapWarehouseByInt64 map[int64]Warehouse

type Warehouses []Warehouse

func (ws Warehouses) MapWarehouseById() MapWarehouseByInt64 {
	mapWarehouseById := make(MapWarehouseByInt64)
	for _, w := range ws {
		mapWarehouseById[w.Id] = w
	}
	return mapWarehouseById
}

// One returns the first warehouse on the list and the status (found or not found)
func (ws Warehouses) One() (Warehouse, bool) {
	if len(ws) == 0 {
		return Warehouse{}, false
	}

	return ws[0], true
}

func (ws Warehouses) OpenCloseWarehouse(mapWarehouseById MapWarehouseByInt64) {
	for i, w := range ws {
		warehouseById, ok := mapWarehouseById[w.Id]
		if !ok {
			continue
		}

		ws[i].Signature.IsActived = warehouseById.IsActived
	}
}

type WarehouseFilter struct {
	Or       *WarehouseFilter
	Ids      []int64
	StoreIds []int64
	Pagination
}

func (f WarehouseFilter) composeFilter() ([]string, []interface{}) {
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

	return query, args
}

func (f WarehouseFilter) ComposeFilter() (string, []interface{}, error) {
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
