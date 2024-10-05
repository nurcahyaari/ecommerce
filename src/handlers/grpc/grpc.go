package grpc

import (
	"github.com/nurcahyaari/ecommerce/config"
)

type GrpcHandler struct {
}

func NewGrpcHandler(cfg *config.Config) *GrpcHandler {
	return &GrpcHandler{}
}
