package repository

var WarehouseQuery = struct {
	Select           string
	Count            string
	ActivationStatus string
}{
	Select: `
		SELECT
			id,
			store_id,
			name,
			location,
			is_actived,
			is_removed,
			created_at,
			created_by,
			updated_at,
			updated_by,
			deleted_at,
			deleted_by
		FROM
			warehouse
	`,
	Count: `
		SELECT
			COUNT(*) c
		FROM
			warehouse
	`,
	ActivationStatus: `
		UPDATE
			warehouse
		SET
			is_actived = :is_actived
		WHERE id = :id
	`,
}
