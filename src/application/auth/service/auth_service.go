package service

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/nurcahyaari/ecommerce/config"
	internalerrors "github.com/nurcahyaari/ecommerce/internal/x/errors"
	"github.com/nurcahyaari/ecommerce/src/domain/entity"
	"github.com/nurcahyaari/ecommerce/src/domain/repository"
	"github.com/nurcahyaari/ecommerce/src/domain/service"
	"github.com/nurcahyaari/ecommerce/src/transferobject"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type AuthService struct {
	cfg     *config.Config
	log     zerolog.Logger
	userSvc service.UserServicer
	reader  repository.AuthRepositoryReader
	writer  repository.AuthRepositoryWriter
}

func NewAuthService(
	cfg *config.Config,
	log zerolog.Logger,
	reader repository.AuthRepositoryReader,
	writer repository.AuthRepositoryWriter,
	userSvc service.UserServicer,
) service.AuthServicer {
	return &AuthService{
		cfg:     cfg,
		log:     log,
		userSvc: userSvc,
		reader:  reader,
		writer:  writer,
	}
}

func (s *AuthService) ValidateToken(ctx context.Context, request transferobject.RequestValidateToken) (transferobject.ResponseValidateKey, error) {
	key, err := s.cfg.JwtPublicKey()
	if err != nil {
		log.Error().
			Err(err).
			Msg("[ValidateToken.JwtPrivateKey]")
		return transferobject.ResponseValidateKey{}, internalerrors.New(
			errors.New("err: internal service error"),
			internalerrors.SetErrorCode(http.StatusNotFound))
	}

	userToken, err := entity.NewAuth(request.Token, key)
	if err != nil {
		return transferobject.ResponseValidateKey{}, internalerrors.New(
			errors.New("err: token not found"),
			internalerrors.SetErrorCode(http.StatusUnauthorized))
	}

	now := time.Now()
	if userToken.AccessTokenExpired(now) {
		return transferobject.ResponseValidateKey{}, internalerrors.New(
			errors.New("err: token has expired"),
			internalerrors.SetErrorCode(http.StatusUnauthorized))
	}

	return transferobject.ResponseValidateKey{
		Valid: true,
	}, nil
}

func (s *AuthService) GenerateToken(ctx context.Context, request transferobject.RequestGenerateToken) (transferobject.ResponseGenerateToken, error) {
	respGetUser, err := s.userSvc.GetUser(ctx, transferobject.RequestSearchUser{
		Or: &transferobject.RequestSearchUser{
			Email: request.Key,
			Phone: request.Key,
		},
	})
	if err != nil {
		log.Error().
			Err(err).
			Msg("[ValidateToken.FindTokenByAccessToken] error find user access token")
		return transferobject.ResponseGenerateToken{}, internalerrors.New(
			err,
			internalerrors.SetErrorCode(http.StatusUnauthorized))
	}

	user := respGetUser.User.Entity()
	if err := user.ComparePassword(request.Password); err != nil {
		log.Warn().
			Err(err).
			Msg("[ValidateToken.ComparePassword] error password didn't match")
		return transferobject.ResponseGenerateToken{}, internalerrors.New(
			errors.New("err: password didn't match"),
			internalerrors.SetErrorCode(http.StatusNotFound))
	}

	userAuth := user.Auth(s.cfg.Auth.JwtToken.Duration)

	privateKey, err := s.cfg.JwtPrivateKey()
	if err != nil {
		log.Error().
			Err(err).
			Msg("[ValidateToken.JwtPrivateKey]")
		return transferobject.ResponseGenerateToken{}, internalerrors.New(
			errors.New("err: internal service error"),
			internalerrors.SetErrorCode(http.StatusNotFound))
	}

	accessToken, err := userAuth.GeneateJwt(privateKey)
	if err != nil {
		log.Error().
			Err(err).
			Msg("[ValidateToken.GenerateToken] error generate access token")
		return transferobject.ResponseGenerateToken{}, internalerrors.New(err)
	}

	return transferobject.ResponseGenerateToken{
		AccessToken: accessToken,
	}, nil
}

func (s *AuthService) ExtractToken(ctx context.Context, req transferobject.RequestValidateToken) (transferobject.ResponseExtractedToken, error) {
	key, err := s.cfg.JwtPublicKey()
	if err != nil {
		log.Error().
			Err(err).
			Msg("[ValidateToken.JwtPrivateKey]")
		return transferobject.ResponseExtractedToken{}, internalerrors.New(
			errors.New("err: internal service error"),
			internalerrors.SetErrorCode(http.StatusNotFound))
	}

	localUserAuthToken, err := entity.NewAuth(req.Token, key)
	if err != nil {
		log.Error().
			Err(err).
			Msg("[ExtractToken.NewLocalUserAuthTokenFromToken] error token is not valid")
		return transferobject.ResponseExtractedToken{}, internalerrors.New(err)
	}

	return transferobject.ResponseExtractedToken{
		UserId:    localUserAuthToken.Id,
		ExpiredAt: localUserAuthToken.ExpiresAt.Time,
		UserType:  localUserAuthToken.UserType,
	}, nil
}
