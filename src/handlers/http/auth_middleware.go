package http

import (
	"net/http"
	"strings"

	"github.com/nurcahyaari/ecommerce/internal/protocols/http/response"
	internalcontext "github.com/nurcahyaari/ecommerce/internal/x/context"
	"github.com/nurcahyaari/ecommerce/src/domain/entity"
	"github.com/nurcahyaari/ecommerce/src/transferobject"
	"github.com/rs/zerolog/log"
)

func (h *HttpHandle) AuthenticationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authorization := r.Header.Get("authorization")

		token := strings.Split(authorization, " ")

		ctx := r.Context()

		if len(token) < 2 && token[0] != string(entity.BearerAuthenticationToken) {
			log.Error().
				Msg("[MiddlewareLocalAuthAccessTokenValidate.ValidateToken] token mismatch")
			response.Err[string](
				w,
				response.SetErr[string]("err: token mismatch"),
				response.SetHttpCode[string](http.StatusUnauthorized),
			)
			return
		}

		resp, err := h.authService.ValidateToken(ctx, transferobject.RequestValidateToken{
			Token: token[1],
		})
		if err != nil {
			log.Error().
				Err(err).
				Msg("[MiddlewareLocalAuthAccessTokenValidate.ValidateToken]")
			response.Err[string](
				w,
				response.SetErr[string]("err: unauthorized"),
				response.SetHttpCode[string](http.StatusUnauthorized),
			)
			return
		}
		if !resp.Valid {
			log.Error().
				Msg("[MiddlewareLocalAuthAccessTokenValidate.ValidateToken] token isn't valid")
			response.Err[string](
				w,
				response.SetErr[string]("err: unauthorized"),
				response.SetHttpCode[string](http.StatusUnauthorized),
			)
			return
		}

		extracted, err := h.authService.ExtractToken(ctx, transferobject.RequestValidateToken{
			Token: token[1],
		})
		if err != nil {
			log.Error().
				Msg("[MiddlewareLocalAuthAccessTokenValidate.ExtractToken] token cannot be extracted")
			response.Err[string](
				w,
				response.SetErr[string]("err: unauthorized"),
				response.SetHttpCode[string](http.StatusUnauthorized),
			)
			return
		}

		r.Header.Set("x-user-id", extracted.UserId)
		ctx = internalcontext.SetUserId(ctx, extracted.UserId)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
