package entity

type Store struct {
	Id   int64  `db:"id"`
	Name string `db:"name"`
}
