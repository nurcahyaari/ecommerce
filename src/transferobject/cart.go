package transferobject

import (
	"context"
	"strconv"

	internalcontext "github.com/nurcahyaari/ecommerce/internal/x/context"
	"github.com/nurcahyaari/ecommerce/src/domain/entity"
	"github.com/shopspring/decimal"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"gopkg.in/guregu/null.v4"
)

type Cart struct {
	Id            string          `json:"id"`
	UserId        int64           `json:"user_id"`
	UserAddressId int64           `json:"user_address_id"`
	TotalQuantity int32           `json:"total_quantity"`
	TotalPrice    decimal.Decimal `json:"total_price"`
	Status        int             `json:"status"`
	CartItems     CartItems       `json:"cart_items"`
}

func (c Cart) Entity() (entity.Cart, error) {
	id, err := primitive.ObjectIDFromHex(c.Id)
	if err != nil {
		return entity.Cart{}, err
	}

	totalPrice, err := primitive.ParseDecimal128(c.TotalPrice.String())
	if err != nil {
		return entity.Cart{}, err
	}

	cartItems, err := c.CartItems.Entity()
	if err != nil {
		return entity.Cart{}, err
	}

	return entity.Cart{
		Id:            id,
		UserId:        c.UserId,
		UserAddressId: c.UserAddressId,
		TotalQuantity: c.TotalQuantity,
		TotalPrice:    totalPrice,
		Status:        c.Status,
		CartItems:     cartItems,
	}, nil
}

func NewCart(cart entity.Cart) (Cart, error) {
	totalPrice, err := decimal.NewFromString(cart.TotalPrice.String())
	if err != nil {
		return Cart{}, err
	}

	ci, err := NewCartItems(cart.CartItems)
	if err != nil {
		return Cart{}, err
	}
	c := Cart{
		Id:            cart.Id.Hex(),
		UserId:        cart.UserId,
		UserAddressId: cart.UserAddressId,
		TotalQuantity: cart.TotalQuantity,
		TotalPrice:    totalPrice,
		Status:        cart.Status,
		CartItems:     ci,
	}

	return c, nil
}

type Carts []Cart

func (cs Carts) Entity() (entity.Carts, error) {
	carts := entity.Carts{}
	for _, c := range cs {
		cart, err := c.Entity()
		if err != nil {
			return entity.Carts{}, err
		}

		carts = append(carts, cart)
	}
	return carts, nil
}

func NewCarts(carts entity.Carts) (Carts, error) {
	respCarts := Carts{}
	for _, cart := range carts {
		c, err := NewCart(cart)
		if err != nil {
			return respCarts, err
		}
		respCarts = append(respCarts, c)
	}
	return respCarts, nil
}

type CartItem struct {
	Id              string          `json:"id"`
	CartId          string          `json:"cart_id"`
	ProductId       int64           `json:"product_id"`
	Quantity        int32           `json:"quantity"`
	PricePerProduct decimal.Decimal `json:"price_per_product"`
	TotalPrice      decimal.Decimal `json:"total_price"`
}

func (ci CartItem) Entity() (entity.CartItem, error) {
	id, err := primitive.ObjectIDFromHex(ci.Id)
	if err != nil {
		return entity.CartItem{}, err
	}

	totalPrice, err := primitive.ParseDecimal128(ci.TotalPrice.String())
	if err != nil {
		return entity.CartItem{}, err
	}

	pricePerProduct, err := primitive.ParseDecimal128(ci.PricePerProduct.String())
	if err != nil {
		return entity.CartItem{}, err
	}

	return entity.CartItem{
		Id:              id,
		CartId:          ci.CartId,
		ProductId:       ci.ProductId,
		Quantity:        ci.Quantity,
		PricePerProduct: pricePerProduct,
		TotalPrice:      totalPrice,
	}, nil
}

func NewCartItem(cartItem entity.CartItem) (CartItem, error) {
	pricePerProduct, err := decimal.NewFromString(cartItem.PricePerProduct.String())
	if err != nil {
		return CartItem{}, err
	}

	totalPrice, err := decimal.NewFromString(cartItem.TotalPrice.String())
	if err != nil {
		return CartItem{}, err
	}

	return CartItem{
		Id:              cartItem.Id.Hex(),
		CartId:          cartItem.CartId,
		ProductId:       cartItem.ProductId,
		Quantity:        cartItem.Quantity,
		PricePerProduct: pricePerProduct,
		TotalPrice:      totalPrice,
	}, nil
}

type CartItems []CartItem

func (cis CartItems) Entity() (entity.CartItems, error) {
	cartItems := entity.CartItems{}

	for _, ci := range cis {
		cartItem, err := ci.Entity()
		if err != nil {
			return entity.CartItems{}, err
		}
		cartItems = append(cartItems, cartItem)
	}

	return cartItems, nil
}

func NewCartItems(cartItems entity.CartItems) (CartItems, error) {
	respCartItems := CartItems{}

	for _, cartItem := range cartItems {
		ci, err := NewCartItem(cartItem)
		if err != nil {
			return respCartItems, err
		}
		respCartItems = append(respCartItems, ci)
	}

	return respCartItems, nil
}

type RequestGetCart struct {
	UserId string
}

func (r *RequestGetCart) PopulateContext(ctx context.Context) {
	userId := internalcontext.GetUserId(ctx)
	r.UserId = userId
}

func (r RequestGetCart) CartFilter() (entity.CartFilter, error) {
	filter := entity.CartFilter{
		Id: primitive.NilObjectID,
	}
	userId, err := strconv.ParseInt(r.UserId, 10, 64)
	if err != nil {
		return filter, err
	}
	filter.UserId = null.IntFrom(userId)

	return filter, nil
}

type ResponseGetCart struct {
	Carts Carts `json:"cart"`
}

func NewResponseGetCart(carts entity.Carts) (ResponseGetCart, error) {
	cart, err := NewCarts(carts)
	if err != nil {
		return ResponseGetCart{}, err
	}

	return ResponseGetCart{
		Carts: cart,
	}, nil
}

type RequestAddItemToCart struct {
	UserAddressId int64 `json:"user_address_id"`
	ProductId     int64 `json:"product_id"`
	Quantity      uint  `json:"quantity"`
}

type RequestDeleteCart struct {
	UserId string
}

func (r RequestDeleteCart) UserIdInt() (int64, error) {
	return strconv.ParseInt(r.UserId, 10, 64)
}
