package http

import (
	"github.com/go-chi/chi/v5"
	"github.com/nurcahyaari/ecommerce/src/domain/service"
)

type HttpHandle struct {
	userService      service.UserServicer
	authService      service.AuthServicer
	productService   service.ProductServicer
	warehouseService service.WarehouseServicer
	cartService      service.CartServicer
	orderService     service.OrderServicer
}

func (h HttpHandle) Router(r *chi.Mux) {
	r.Route("/v1", func(r chi.Router) {
		r.Route("/auth", func(r chi.Router) {
			r.Post("/login", h.Login)
		})
		r.With(h.AuthenticationMiddleware).Get("/users", h.SearchUsers)
		r.Group(func(r chi.Router) {
			r.Use(h.AuthenticationMiddleware)
			r.Route("/products", func(r chi.Router) {
				r.Get("/", h.SearchProducts)
				r.Put("/{productId}/warehouse", h.MoveWarehouse)
			})

			r.Route("/warehouses", func(r chi.Router) {
				r.Put("/activation-status", h.WarehouseActivationStatus)
			})

			r.Route("/carts", func(r chi.Router) {
				r.Get("/", h.GetCart)
				r.Post("/", h.AddItemToCart)
			})

			r.Route("/orders", func(r chi.Router) {
				r.Post("/", h.CreateOrder)
			})
		})
	})
}

func NewHttpHandler(
	userService service.UserServicer,
	authService service.AuthServicer,
	productService service.ProductServicer,
	warehouseService service.WarehouseServicer,
	cartService service.CartServicer,
	orderService service.OrderServicer,
) *HttpHandle {
	httpHandle := &HttpHandle{
		userService:      userService,
		authService:      authService,
		productService:   productService,
		warehouseService: warehouseService,
		cartService:      cartService,
		orderService:     orderService,
	}
	return httpHandle
}
