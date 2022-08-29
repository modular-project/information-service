package handler

import (
	"context"
	"fmt"
	"information-service/controller"

	pb "github.com/modular-project/protobuffers/information/order"
)

type OrderService struct {
	pb.UnimplementedValidateOrderServer
}

func (s *OrderService) ValidateOrder(c context.Context, in *pb.ValidateOrderRequest) (*pb.ValidateResponse, error) {
	t, err := controller.ValidateOrder(in.Order)
	if err != nil {
		return nil, fmt.Errorf("validateOrder: %w", err)
	}
	return &pb.ValidateResponse{Total: float32(t)}, err
}

func (s *OrderService) ValidateProducts(c context.Context, in *pb.ValidateProductsRequest) (*pb.ValidateResponse, error) {
	t, err := controller.CheckProducts(in.OrderProducts)
	if err != nil {
		return nil, fmt.Errorf("validateOrder: %w", err)
	}
	return &pb.ValidateResponse{Total: float32(t)}, err
}
