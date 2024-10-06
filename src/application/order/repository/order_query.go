package repository

var OrderQuery = struct {
	Select            string
	Count             string
	Create            string
	UpdateOrderStatus string
}{
	Select: `
		SELECT
			id,
			user_id,
			total_price,
			total_quantity,
			order_code,
			order_status,
			expired_order,
			is_removed,
			created_at,
			created_by,
			updated_at,
			updated_by,
			deleted_at,
			deleted_by
		FROM
			` + "`order`" + `
	`,
	Count: `
		SELECT
			COUNT(*) c
		FROM
			` + "`order`" + `
	`,
	Create: `
		INSERT INTO ` + "`order`" + ` (
			user_id,
			total_price,
			total_quantity,
			order_code,
			order_status,
			expired_order
		) VALUE (
			:user_id,
			:total_price,
			:total_quantity,
			:order_code,
			:order_status,
			:expired_order
		)
	`,
	UpdateOrderStatus: `
		UPDATE ` + "`order`" + `
		SET
			order_status = :order_status
		WHERE
			id = :id
	`,
}
