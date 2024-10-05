package entity

import "github.com/shopspring/decimal"

type OrderStatus int

func (w OrderStatus) String() string {
	return [...]string{"Pending", "Paid", "Processed", "Received", "Done", "Refund", "Expired"}[w-1]
}

func (w OrderStatus) EnumIndex() int {
	return int(w)
}

const (
	Pending OrderStatus = iota + 1
	Paid
	Processed
	Received
	Done
	Refund
	Expired
)

type Order struct {
	Id            int64           `db:"id"`
	UserId        int64           `db:"user_id"`
	OrderCode     string          `db:"order_code"`
	TotalQuantity int32           `db:"total_quantity"`
	TotalPrice    decimal.Decimal `db:"total_price"`
	Status        int             `db:"status"`
	Signature
}

type OrderDetail struct {
	Id              int64           `db:"id"`
	OrderId         int64           `db:"order_id"`
	Quantity        int32           `db:"quantity"`
	PricePerProduct decimal.Decimal `db:"price_per_product"`
	TotalPrice      decimal.Decimal `db:"total_price"`
	Signature
}
