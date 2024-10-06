package http

import (
	"encoding/json"
	"net/http"

	"github.com/nurcahyaari/ecommerce/internal/protocols/http/response"
	"github.com/nurcahyaari/ecommerce/src/transferobject"
)

func (h *HttpHandle) WarehouseActivationStatus(w http.ResponseWriter, r *http.Request) {
	request := transferobject.RequestOpenCloseWarehouse{}

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		h.log.Error().Err(err).Stack().Msg("WarehouseActivationStatus.Decode")
		response.Err[string](w,
			response.SetMessage[string](err.Error()))
		return
	}

	resp, err := h.warehouseService.OpenCloseWarehouse(r.Context(), request)
	if err != nil {
		h.log.Error().Err(err).Stack().Msg("WarehouseActivationStatus.OpenCloseWarehouse")
		response.Err[any](w,
			response.SetErr[any](err),
			response.SetHttpCode[any](500))
		return
	}

	response.Json[transferobject.Warehouses](w,
		response.SetMessage[transferobject.Warehouses]("success"),
		response.SetData[transferobject.Warehouses](resp))
}
