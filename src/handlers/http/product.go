package http

import (
	"encoding/json"
	"net/http"

	"github.com/nurcahyaari/ecommerce/internal/protocols/http/response"
	"github.com/nurcahyaari/ecommerce/src/transferobject"
)

func (h *HttpHandle) SearchProducts(w http.ResponseWriter, r *http.Request) {
	request := transferobject.RequestSearchProduct{
		Ids:          r.FormValue("ids"),
		WarehouseIds: r.FormValue("warehouse_ids"),
		StoreIds:     r.FormValue("store_ids"),
	}

	resp, err := h.productService.SearchProducts(r.Context(), request)
	if err != nil {
		response.Err[any](w,
			response.SetErr[any](err),
			response.SetHttpCode[any](500))
		return
	}

	response.Json[transferobject.ResponseSearchProduct](w,
		response.SetMessage[transferobject.ResponseSearchProduct]("success"),
		response.SetData[transferobject.ResponseSearchProduct](resp))
}

func (h *HttpHandle) MoveWarehouse(w http.ResponseWriter, r *http.Request) {
	request := transferobject.RequestMoveWarehouse{}

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		response.Err[string](w,
			response.SetMessage[string](err.Error()))
		return
	}

	if err := request.ParseURLParams(r); err != nil {
		response.Err[string](w,
			response.SetMessage[string](err.Error()))
		return
	}

	resp, err := h.productService.MoveWarehouse(r.Context(), request)
	if err != nil {
		response.Err[any](w,
			response.SetErr[any](err),
			response.SetHttpCode[any](500))
		return
	}

	response.Json[transferobject.Product](w,
		response.SetMessage[transferobject.Product]("success"),
		response.SetData[transferobject.Product](resp))
}
