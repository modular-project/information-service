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
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/xds"
)

func Start() *xds.GRPCServer {
	opts := []grpc_recovery.Option{
		grpc_recovery.WithRecoveryHandler(middleware.Recovery),
	}
	server := xds.NewGRPCServer(
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_recovery.UnaryServerInterceptor(opts...),
		)),
		grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
			grpc_middleware.ChainStreamServer(),
		)),
	)
	//d := grpc_middleware.WithUnaryServerChain(grpc_recovery.UnaryServerInterceptor(opts...))
	// server := grpc.NewServer(
	// 	grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
	// 		grpc_recovery.UnaryServerInterceptor(opts...),
	// 	)),
	// 	grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
	// 		grpc_middleware.ChainStreamServer(),
	// 	)),
	// )
	healthServer := health.NewServer()
	healthServer.SetServingStatus("", healthpb.HealthCheckResponse_SERVING)
	pbp.RegisterProductServiceServer(server, &handler.ProductService{})
	healthServer.SetServingStatus(pbp.ProductService_ServiceDesc.ServiceName, healthpb.HealthCheckResponse_SERVING)
	pbt.RegisterTableServiceServer(server, &handler.TableService{})
	healthServer.SetServingStatus(pbt.TableService_ServiceDesc.ServiceName, healthpb.HealthCheckResponse_SERVING)
	pbe.RegisterEstablishmentServiceServer(server, &handler.EstablishmentService{})
	healthServer.SetServingStatus(pbe.EstablishmentService_ServiceDesc.ServiceName, healthpb.HealthCheckResponse_SERVING)
	pbo.RegisterValidateOrderServer(server, &handler.OrderService{})
	healthServer.SetServingStatus(pbo.ValidateOrder_ServiceDesc.ServiceName, healthpb.HealthCheckResponse_SERVING)
	reflection.Register(server)
	return server
}
