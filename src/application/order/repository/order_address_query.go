package repository

var OrderAddressQuery = struct {
	Insert string
}{
	Insert: `
		INSERT INTO order_address (
			order_id,
			user_id,
			full_address
		) VALUES (
			:order_id,
			:user_id,
			:full_address
		)
	`,
}
