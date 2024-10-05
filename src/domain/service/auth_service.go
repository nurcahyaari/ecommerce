package service

import (
	"context"

	"github.com/nurcahyaari/ecommerce/src/transferobject"
)

type AuthServicer interface {
	GenerateToken(context.Context, transferobject.RequestGenerateToken) (transferobject.ResponseGenerateToken, error)
	ValidateToken(context.Context, transferobject.RequestValidateToken) (transferobject.ResponseValidateKey, error)
	ExtractToken(context.Context, transferobject.RequestValidateToken) (transferobject.ResponseExtractedToken, error)
}
