package grpc

import (
	"fmt"
	"log"
	"net"

	"github.com/nurcahyaari/ecommerce/config"
	grpchandler "github.com/nurcahyaari/ecommerce/src/handlers/grpc"
	"google.golang.org/grpc"
)

type Grpc struct {
	grpcServer *grpc.Server
	cfg        *config.Config
	handler    *grpchandler.GrpcHandler
}

func New(cfg *config.Config, handler *grpchandler.GrpcHandler) Grpc {
	return Grpc{
		cfg:     cfg,
		handler: handler,
	}
}

func (g *Grpc) Listen() {
	addr := fmt.Sprintf("127.0.0.1:%d", g.cfg.Application.Transport.Grpc.PORT)
	listen, err := net.Listen("tcp", addr)
	if err != nil {
		panic(err)
	}
	defer listen.Close()

	log.Printf("Listening on : %v \n", addr)

	grpcServer := grpc.NewServer()

	if err := grpcServer.Serve(listen); err != nil {
		panic(err)
	}

	g.grpcServer = grpcServer
}

func (g Grpc) Shutdown() {
	g.grpcServer.GracefulStop()
}
