package transferobject

import (
	"context"
	"strconv"

	internalcontext "github.com/nurcahyaari/ecommerce/internal/x/context"
	"github.com/nurcahyaari/ecommerce/src/domain/entity"
	"github.com/shopspring/decimal"
)

type Order struct {
	OrderCode     string          `json:"order_code"`
	TotalQuantity int32           `json:"total_quantity"`
	TotalPrice    decimal.Decimal `json:"total_price"`
}

func NewOrder(order entity.Order) Order {
	return Order{
		OrderCode:     order.OrderCode,
		TotalQuantity: order.TotalQuantity,
		TotalPrice:    order.TotalPrice,
	}
}

type Orders []Order

type RequestCreateOrder struct {
	UserId string
}

func (r RequestCreateOrder) UserIdInt() (int64, error) {
	return strconv.ParseInt(r.UserId, 10, 64)
}

func (r *RequestCreateOrder) PopulateContext(ctx context.Context) {
	userId := internalcontext.GetUserId(ctx)
	r.UserId = userId
}
