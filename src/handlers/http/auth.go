package http

import (
	"encoding/json"
	"net/http"

	"github.com/nurcahyaari/ecommerce/internal/protocols/http/response"
	"github.com/nurcahyaari/ecommerce/src/transferobject"
)

func (h *HttpHandle) Login(w http.ResponseWriter, r *http.Request) {
	request := transferobject.RequestGenerateToken{}

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		response.Err[string](w,
			response.SetMessage[string](err.Error()))
		return
	}

	resp, err := h.authService.GenerateToken(r.Context(), request)
	if err != nil {
		response.Err[any](w,
			response.SetErr[any](err),
			response.SetHttpCode[any](500))
		return
	}

	response.Json[transferobject.ResponseGenerateToken](w,
		response.SetMessage[transferobject.ResponseGenerateToken]("success"),
		response.SetData[transferobject.ResponseGenerateToken](resp))
}
