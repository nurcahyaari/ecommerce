package http

import (
	"net/http"

	"github.com/nurcahyaari/ecommerce/internal/protocols/http/response"
	"github.com/nurcahyaari/ecommerce/src/transferobject"
)

func (h *HttpHandle) SearchUsers(w http.ResponseWriter, r *http.Request) {
	request := transferobject.RequestSearchUser{
		Ids:   r.FormValue("ids"),
		Email: r.FormValue("email"),
	}

	resp, err := h.userService.SearchUsers(r.Context(), request)
	if err != nil {
		response.Err[any](w,
			response.SetErr[any](err),
			response.SetHttpCode[any](500))
		return
	}

	response.Json[transferobject.ResponseSearchUser](w,
		response.SetMessage[transferobject.ResponseSearchUser]("success"),
		response.SetData[transferobject.ResponseSearchUser](resp))
}
