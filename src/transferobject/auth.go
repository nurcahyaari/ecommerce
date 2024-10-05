package transferobject

import (
	"net/http"
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
	internalerror "github.com/nurcahyaari/ecommerce/internal/x/errors"
)

type RequestGenerateToken struct {
	// key can be phone or email
	Key string `json:"key"`
	// secret can be user password
	Password string `json:"password"`
}

func (r RequestGenerateToken) Validate() error {
	err := validation.ValidateStruct(&r,
		validation.Field(&r.Key, validation.Required),
		validation.Field(&r.Password, validation.Required),
	)

	if err == nil {
		return nil
	}

	return internalerror.New(err,
		internalerror.SetErrorCode(http.StatusBadRequest))
}

type ResponseGenerateToken struct {
	AccessToken     string `json:"accessToken"`
	AccessTokenExp  string `json:"accessTokenExp"`
	RefreshToken    string `json:"refreshToken"`
	RefreshTokenExp string `json:"refreshTokenExp"`
}

type RequestValidateToken struct {
	Token string
}

type ResponseValidateKey struct {
	Valid bool `json:"valid"`
}

type ResponseExtractedToken struct {
	UserId    string    `json:"userId"`
	ExpiredAt time.Time `json:"expiredAt"`
	UserType  string    `json:"userType"`
}
