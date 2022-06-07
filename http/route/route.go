package route

import (
	"information-service/http/handler"
	"information-service/http/middleware"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	pbe "github.com/modular-project/protobuffers/information/establishment"
	pbo "github.com/modular-project/protobuffers/information/order"
	pbp "github.com/modular-project/protobuffers/information/product"
	pbt "github.com/modular-project/protobuffers/information/table"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func Start() *grpc.Server {
	opts := []grpc_recovery.Option{
		grpc_recovery.WithRecoveryHandler(middleware.Recovery),
	}
	//d := grpc_middleware.WithUnaryServerChain(grpc_recovery.UnaryServerInterceptor(opts...))
	server := grpc.NewServer(
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_recovery.UnaryServerInterceptor(opts...),
		)),
		grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
			grpc_middleware.ChainStreamServer(),
		)),
	)
	pbp.RegisterProductServiceServer(server, &handler.ProductService{})
	pbt.RegisterTableServiceServer(server, &handler.TableService{})
	pbe.RegisterEstablishmentServiceServer(server, &handler.EstablishmentService{})
	pbo.RegisterValidateOrderServer(server, &handler.OrderService{})
	reflection.Register(server)
	return server
}
