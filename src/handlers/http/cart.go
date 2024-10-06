package http

import (
	"encoding/json"
	"net/http"

	"github.com/nurcahyaari/ecommerce/internal/protocols/http/response"
	"github.com/nurcahyaari/ecommerce/src/transferobject"
)

func (h *HttpHandle) GetCart(w http.ResponseWriter, r *http.Request) {
	request := transferobject.RequestGetCart{}

	ctx := r.Context()
	request.PopulateContext(ctx)

	resp, err := h.cartService.GetCart(ctx, request)
	if err != nil {
		response.Err[any](w,
			response.SetErr[any](err),
			response.SetHttpCode[any](500))
		return
	}

	response.Json[transferobject.ResponseGetCart](w,
		response.SetMessage[transferobject.ResponseGetCart]("success"),
		response.SetData[transferobject.ResponseGetCart](resp))
}

func (h *HttpHandle) AddItemToCart(w http.ResponseWriter, r *http.Request) {
	request := transferobject.RequestAddItemToCart{}

	ctx := r.Context()
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		h.log.Error().Err(err).Stack().Msg("AddItemToCart.Decode")
		response.Err[string](w,
			response.SetMessage[string](err.Error()))
		return
	}

	cart, err := h.cartService.AddItemToCart(ctx, request)
	if err != nil {
		h.log.Error().Err(err).Stack().Msg("AddItemToCart.AddItemToCart")
		response.Err[any](w,
			response.SetErr[any](err),
			response.SetHttpCode[any](500))
		return
	}

	response.Json[transferobject.Cart](w,
		response.SetMessage[transferobject.Cart]("success"),
		response.SetData[transferobject.Cart](cart))
}
