package repository

var UserAddressQuery = struct {
	Select string
	Count  string
}{
	Select: `
		SELECT
			id,
			user_id,
			full_address,
			is_actived,
			is_removed,
			created_at,
			created_by,
			updated_at,
			updated_by,
			deleted_at,
			deleted_by
		FROM
			user_address
	`,
	Count: `
		SELECT
			COUNT(*) c
		FROM
			user_address
	`,
}
