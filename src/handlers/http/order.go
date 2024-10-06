package http

import (
	"net/http"

	"github.com/nurcahyaari/ecommerce/internal/protocols/http/response"
	"github.com/nurcahyaari/ecommerce/src/transferobject"
)

func (h *HttpHandle) CreateOrder(w http.ResponseWriter, r *http.Request) {
	request := transferobject.RequestCreateOrder{}

	ctx := r.Context()
	request.PopulateContext(ctx)

	order, err := h.orderService.CreateOrder(ctx, request)
	if err != nil {
		h.log.Error().Err(err).Stack().Msg("CreateOrder.CreateOrder")
		response.Err[any](w,
			response.SetErr[any](err),
			response.SetHttpCode[any](500))
		return
	}

	response.Json[transferobject.Order](w,
		response.SetMessage[transferobject.Order]("success"),
		response.SetData[transferobject.Order](order))
}
