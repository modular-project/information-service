package handler

import (
	"context"
	"information-service/controller"

	pb "github.com/modular-project/protobuffers/information/order"
)

type OrderService struct {
	pb.UnimplementedValidateOrderServer
}

func (s *OrderService) ValidateOrder(c context.Context, in *pb.ValidateOrderRequest) (*pb.ValidateOrderResponse, error) {
	err := controller.ValidateOrder(in.Order)
	return nil, err
}
