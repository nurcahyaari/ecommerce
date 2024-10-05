package router

import (
	"github.com/go-chi/chi/v5"
	httphandler "github.com/nurcahyaari/ecommerce/src/handlers/http"
)

type HttpRoute struct {
	handler *httphandler.HttpHandle
}

func (h *HttpRoute) Router(r *chi.Mux) {
	h.handler.Router(r)
}

func NewHttpRouter(
	handler *httphandler.HttpHandle,
) *HttpRoute {
	return &HttpRoute{
		handler: handler,
	}
}
