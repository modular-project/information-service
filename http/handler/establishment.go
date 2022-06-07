package handler

import (
	"context"
	"information-service/controller"
	"information-service/model"

	pb "github.com/modular-project/protobuffers/information/establishment"
)

type EstablishmentService struct {
	pb.UnimplementedEstablishmentServiceServer
}

// Create(ctx context.Context, in *pb.Establishment) (*pb.Response, error)
func (s *EstablishmentService) GetAll(ctx context.Context, in *pb.RequestGetAll) (*pb.ResponseGetAll, error) {
	ms, err := controller.GetEstablishmentsAvailable()
	if err != nil {
		return nil, err
	}
	res := make([]*pb.Establishment, len(ms))
	for i, m := range ms {
		res[i] = &pb.Establishment{Id: uint64(m.ID)}
	}
	return &pb.ResponseGetAll{Establishments: res}, nil
}

func (s *EstablishmentService) Update(ctx context.Context, in *pb.RequestUpdate) (*pb.Response, error) {
	m := model.Establishment{Model: model.Model{ID: uint(in.Establishment.Id)}}
	err := controller.UpdateEstablishmentData(&m)
	if err != nil {
		return nil, err
	}
	return &pb.Response{}, nil
}

func (s *EstablishmentService) Delete(ctx context.Context, in *pb.RequestById) (*pb.ResponseDelete, error) {
	err := controller.RemoveEstablishment(&model.Establishment{Model: model.Model{ID: uint(in.Id)}})
	if err != nil {
		return nil, err
	}
	return &pb.ResponseDelete{}, nil
}

func (s *EstablishmentService) Create(ctx context.Context, in *pb.RequestCreate) (*pb.Response, error) {
	m := model.Establishment{}
	err := controller.CreateEstablishment(&m)
	if err != nil {
		return nil, err
	}
	return &pb.Response{Id: uint64(m.ID)}, nil
}
