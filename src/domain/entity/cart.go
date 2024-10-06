package entity

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/shopspring/decimal"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"gopkg.in/guregu/null.v4"
)

type Cart struct {
	Id            primitive.ObjectID   `bson:"_id"`
	UserId        int64                `bson:"user_id"`
	UserAddressId int64                `bson:"user_address_id"`
	TotalQuantity int32                `bson:"total_quantity"`
	TotalPrice    primitive.Decimal128 `bson:"total_price"`
	Status        int                  `bson:"status"`
	CartItems     CartItems            `bson:"cart_items"`
}

func (c *Cart) SumTotalItemQuantity() {
	quantity := int32(0)
	for _, ci := range c.CartItems {
		quantity += ci.Quantity
	}

	c.TotalQuantity = quantity
}

func (c *Cart) SumTotalItemPrice() error {
	price := decimal.Decimal{}
	for _, ci := range c.CartItems {
		totalPrice, err := ci.DecimalTotalPrice()
		if err != nil {
			return err
		}
		price = price.Add(totalPrice)
	}

	val, err := primitive.ParseDecimal128(price.String())
	if err != nil {
		return err
	}
	c.TotalPrice = val
	return nil
}

func (c *Cart) UpdateCartItems(productForCart ProductForCart) error {
	mapByProductId := c.CartItems.MapByProductId()

	cartItem, ok := mapByProductId[productForCart.ProductId]
	if !ok {
		return nil
	}

	cartItem.Quantity += int32(productForCart.Quantity)
	totalPrice, err := cartItem.TotalProductPrice()
	if err != nil {
		return err
	}
	primitiveDecimalTotalPrice, err := primitive.ParseDecimal128(totalPrice.String())
	if err != nil {
		return err
	}
	cartItem.TotalPrice = primitiveDecimalTotalPrice
	mapByProductId[productForCart.ProductId] = cartItem

	c.CartItems = mapByProductId.CartItems()
	c.SumTotalItemPrice()
	c.SumTotalItemQuantity()

	return nil
}

func NewCart(
	userAddress UserAddress,
	productForCart ProductForCart,
) (Cart, error) {
	c := Cart{
		Id:            primitive.NewObjectID(),
		UserId:        userAddress.UserId,
		UserAddressId: userAddress.Id,
	}

	cartItem, err := productForCart.CartItem(c.Id.Hex())
	if err != nil {
		return c, err
	}

	c.CartItems = append(c.CartItems, cartItem)

	c.SumTotalItemPrice()
	c.SumTotalItemQuantity()

	return c, nil
}

type CartFilter struct {
	Id            primitive.ObjectID
	UserId        null.Int
	UserAddressId null.Int
}

func (cf CartFilter) Filter() bson.M {
	filter := make(bson.M)

	if !cf.Id.IsZero() {
		filter["_id"] = cf.Id
	}
	if cf.UserId.Valid {
		filter["user_id"] = cf.UserId.ValueOrZero()
	}

	if cf.UserAddressId.Valid {
		filter["user_address_id"] = cf.UserAddressId.ValueOrZero()
	}

	return filter
}

type MapCartById map[int64]Cart

type Carts []Cart

func (cs Carts) UserAddressStrs() string {
	userAddresses := []string{}
	for _, c := range cs {
		userAddresses = append(userAddresses, fmt.Sprintf("%d", c.UserAddressId))
	}
	return strings.Join(userAddresses, ",")
}

func (cs Carts) Order(userId int64, mapUserAddress MapUserAddress) (Order, error) {
	var (
		order         = Order{}
		orderReceipts = OrderReceipts{}
	)

	for _, c := range cs {
		totalPrice, err := decimal.NewFromString(c.TotalPrice.String())
		if err != nil {
			return Order{}, err
		}

		orderDetails := OrderDetails{}
		for _, ci := range c.CartItems {
			totalPrice, err := decimal.NewFromString(ci.TotalPrice.String())
			if err != nil {
				return Order{}, err
			}

			pricePerProduct, err := decimal.NewFromString(ci.PricePerProduct.String())
			if err != nil {
				return Order{}, err
			}
			orderDetails = append(orderDetails, OrderDetail{
				ProdutId:        ci.ProductId,
				Quantity:        ci.Quantity,
				PricePerProduct: pricePerProduct,
				TotalPrice:      totalPrice,
			})
		}

		userAddress, ok := mapUserAddress[c.UserAddressId]
		if !ok {
			return Order{}, errors.New("err: cart is not valid because user doesn't have the address")
		}

		orderAddress := OrderAddress{
			UserAddressId: userAddress.Id,
			UserId:        userAddress.UserId,
			FullAddress:   userAddress.FullAddress,
		}

		orderReceipts = append(orderReceipts, OrderReceipt{
			TotalQuantity: c.TotalQuantity,
			TotalPrice:    totalPrice,
			OrderDetails:  orderDetails,
			OrderAddress:  orderAddress,
		})

	}

	order.OrderReceipts = orderReceipts
	order.OrderStatus = Pending
	order.UserId = userId
	order.ExpiredOrder = null.TimeFrom(time.Now().Add(10 * time.Minute))
	order.SumTotalPrice()
	order.SumTotalQuantity()
	order.GenerateOrderCode()

	return order, nil
}

func (cs Carts) One() (Cart, bool) {
	if len(cs) == 0 {
		return Cart{}, false
	}

	return cs[0], true
}

type CartItem struct {
	Id              primitive.ObjectID   `bson:"_id"`
	CartId          string               `bson:"cart_id"`
	ProductId       int64                `bson:"product_id"`
	Quantity        int32                `bson:"quantity"`
	PricePerProduct primitive.Decimal128 `bson:"price_per_product"`
	TotalPrice      primitive.Decimal128 `bson:"total_price"`
}

func (ci CartItem) DecimalPricePerProduct() (decimal.Decimal, error) {
	return decimal.NewFromString(ci.PricePerProduct.String())
}

func (ci CartItem) DecimalTotalPrice() (decimal.Decimal, error) {
	return decimal.NewFromString(ci.TotalPrice.String())
}

func (ci CartItem) TotalProductPrice() (decimal.Decimal, error) {
	val, err := ci.DecimalPricePerProduct()
	if err != nil {
		return decimal.Decimal{}, err
	}

	return val.
		Mul(decimal.NewFromInt(int64(ci.Quantity))), nil
}

type CartItems []CartItem

func (cis CartItems) MapByProductId() MapCartItem {
	mapByProductId := make(MapCartItem)

	for _, ci := range cis {
		mapByProductId[ci.ProductId] = ci
	}

	return mapByProductId
}

type MapCartItem map[int64]CartItem

func (mci MapCartItem) CartItems() CartItems {
	cis := CartItems{}

	for _, ci := range mci {
		cis = append(cis, ci)
	}

	return cis
}
