package entity

import (
	"crypto/rsa"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

type AuthenticationTokenType string

var (
	BearerAuthenticationToken AuthenticationTokenType = "Bearer"
)

type Auth struct {
	jwt.RegisteredClaims
	Id       string `json:"id"`
	UserType string `json:"userType"`
}

func (a Auth) GeneateJwt(key *rsa.PrivateKey) (string, error) {
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodRS256, a)
	token, err := jwtToken.SignedString(key)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (a Auth) AccessTokenExpired(now time.Time) bool {
	return a.ExpiresAt.Before(now)
}

func NewAuth(token string, key *rsa.PublicKey) (Auth, error) {
	localUserAuthToken := Auth{}

	jwtToken, err := jwt.ParseWithClaims(token, &localUserAuthToken, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}

		return key, nil
	})
	if err != nil {
		return localUserAuthToken, err
	}

	if !jwtToken.Valid {
		return localUserAuthToken, errors.New("err: token is not valid")
	}

	return localUserAuthToken, nil
}

func (a Auth) GenerateJWTToken(key *rsa.PrivateKey) (string, error) {
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodRS256, a)
	token, err := jwtToken.SignedString(key)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (a Auth) GenerateUuidToken() string {
	return uuid.New().String()
}

type AuthRefreshToken struct {
	Id        int64     `db:"id"`
	UserId    int64     `db:"user_id"`
	Token     string    `db:"token"`
	IsActive  bool      `db:"is_active"`
	ExpiredAt time.Time `db:"expired_at"`
}

type AuthRefreshTokenFilter struct {
	UserId int64
}
