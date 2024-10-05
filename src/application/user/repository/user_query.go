package repository

var Query = struct {
	Select string
	Count  string
}{
	Select: `
		SELECT
			id,
			name,
			email,
			phone,
			password,
			is_actived,
			is_removed,
			created_at,
			created_by,
			updated_at,
			updated_by,
			deleted_at,
			deleted_by
		FROM
			user
	`,
	Count: `
		SELECT
			COUNT(*) c
		FROM
			user
	`,
}
