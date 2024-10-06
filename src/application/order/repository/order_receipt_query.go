package repository

var OrderReceiptQuery = struct {
	Select string
	Count  string
	Insert string
}{
	Select: `
		SELECT
			id,
			order_id,
			order_address_id,
			total_price,
			total_quantity,
			created_at,
			created_by,
			updated_at,
			updated_by,
			deleted_at,
			deleted_by
		FROM
			` + "`order_receipt`" + `
	`,
	Count: `
		SELECT
			COUNT(*) c
		FROM
			` + "`order_receipt`" + `
	`,
	Insert: `
		INSERT INTO order_receipt (
			order_id,
			order_address_id,
			total_price,
			total_quantity
		) VALUES (
			:order_id,
			:order_address_id,
			:total_price,
			:total_quantity
		)
	`,
}
