package repository

var Query = struct {
	Select            string
	Count             string
	UpdateWarehouse   string
	ActivatingProduct string
}{
	Select: `
		SELECT
			id,
			store_id,
			warehouse_id,
			name,
			price,
			is_actived,
			is_removed,
			created_at,
			created_by,
			updated_at,
			updated_by,
			deleted_at,
			deleted_by
		FROM
			product
	`,
	Count: `
		SELECT
			COUNT(*) c
		FROM
			product
	`,
	UpdateWarehouse: `
		UPDATE product
		SET
			warehouse_id = :warehouse_id
		WHERE
			id = :id
	`,
	ActivatingProduct: `
		UPDATE product
		SET
			is_active = ?
		WHERE
			id = ?
	`,
}

var ProductStockQuery = struct {
	Select       string
	Count        string
	ReserveStock string
	RevertStock  string
	UpdateStock  string
}{
	Select: `
		SELECT
			product_id,
			stock_reserved,
			stock_on_hand
		FROM
			product_stock
	`,
	ReserveStock: `
		UPDATE product_stock 
		SET 
			stock_on_hand = stock_on_hand - ?,
			stock_reserved = stock_reserved + ?
		WHERE
			product_id = ?
	`,
	RevertStock: `
		UPDATE product_stock 
		SET 
			stock_on_hand = stock_on_hand + ?,
			stock_reserved = stock_reserved - ?
		WHERE
			product_id = ?
	`,
	UpdateStock: `
		UPDATE product_stock 
		SET 
			stock_reserved = stock_reserved - ?
		WHERE
			product_id = ?
	`,
}
