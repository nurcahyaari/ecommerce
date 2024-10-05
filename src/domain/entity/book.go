package entity

import (
	"strings"

	"gopkg.in/guregu/null.v4"
)

type Book struct {
	Id            int64     `db:"id"`
	AuthorId      int64     `db:"author_id"`
	CategoryId    int64     `db:"category_id"`
	Name          string    `db:"name"`
	PublishedYear int16     `db:"published_year"`
	Stock         BookStock `db:"-"`
}

type Books []Book

type BookFilter struct {
	Id            null.Int
	AuthorId      null.Int
	CategoryId    null.Int
	PublishedYear null.Int
	Page          int
	Size          int
}

func (f BookFilter) Pagination() (string, []interface{}) {
	page := (f.Page - 1) * f.Size
	args := make([]interface{}, 0)
	args = append(args, f.Size)
	args = append(args, page)
	return "LIMIT ? OFFSET ?", args
}

func (f BookFilter) ComposeFilter() (string, []interface{}, error) {
	var (
		query []string
		args  = make([]interface{}, 0)
	)

	if f.Id.Valid {
		query = append(query, "id = ?")
		args = append(args, f.Id.ValueOrZero())
	}

	if f.AuthorId.Valid {
		query = append(query, "author_id = ?")
		args = append(args, f.AuthorId.ValueOrZero())
	}

	if f.CategoryId.Valid {
		query = append(query, "category_id = ?")
		args = append(args, f.CategoryId.ValueOrZero())
	}

	// Combine query parts into a single WHERE clause
	whereClause := strings.Join(query, " AND ")
	if whereClause != "" {
		whereClause = "WHERE " + whereClause
	}

	return whereClause, args, nil
}

type BookStock struct {
	BookId         int64  `db:"bookId"`
	StockAvailable uint64 `db:"stock_available"`
	StockBorrowing uint64 `db:"stock_borrowing"`
}

type BookStocks []BookStock

type BookBorrowing struct {
	Id       int64 `db:"id"`
	BookId   int64 `db:"book_id"`
	UserId   int64 `db:"user_id"`
	IsReturn Book  `db:"is_return"`
}

type BookBorrowings []BookBorrowing
