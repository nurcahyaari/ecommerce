package repository

var OrderDetailQuery = struct {
	Select string
	Count  string
	Insert string
}{
	Select: `
		SELECT
			id,
			order_id,
			order_receipt_id,
			product_id,
			product_price,
			quantity,
			total_price,
			created_at,
			created_by,
			updated_at,
			updated_by,
			deleted_at,
			deleted_by
		FROM
			` + "`order_detail`" + `
	`,
	Count: `
		SELECT
			COUNT(*) c
		FROM
			` + "`order_detail`" + `
	`,
	Insert: `
		INSERT INTO order_detail (
			order_id,
			order_receipt_id,
			product_id,
			product_price,
			quantity,
			total_price
		) VALUES (
			:order_id,
			:order_receipt_id,
			:product_id,
			:product_price,
			:quantity,
			:total_price
		)
	`,
}
