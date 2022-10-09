package handler

import (
	"context"
	"fmt"
	"information-service/controller"
	"information-service/model"

	pb "github.com/modular-project/protobuffers/information/establishment"
)

type EstablishmentService struct {
	pb.UnimplementedEstablishmentServiceServer
}

func (s *EstablishmentService) GetByAddress(ctx context.Context, in *pb.RequestGetByAddress) (*pb.ResponseAddress, error) {
	eID, q, err := controller.GetByAddress(in.Id)
	if err != nil {
		return &pb.ResponseAddress{}, fmt.Errorf("getByAddress: %w", err)
	}
	return &pb.ResponseAddress{Id: uint64(eID), Quantity: q}, nil
}

// Create(ctx context.Context, in *pb.Establishment) (*pb.Response, error)
func (s *EstablishmentService) GetAll(ctx context.Context, in *pb.RequestGetAll) (*pb.ResponseGetAll, error) {
	ms, err := controller.GetEstablishmentsAvailable()
	if err != nil {
		return &pb.ResponseGetAll{}, err
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
		return &pb.Response{}, err
	}
	return &pb.Response{}, nil
}

func (s *EstablishmentService) Delete(ctx context.Context, in *pb.RequestById) (*pb.ResponseDelete, error) {
	aID, err := controller.RemoveEstablishment(&model.Establishment{Model: model.Model{ID: uint(in.Id)}})
	if err != nil {
		return &pb.ResponseDelete{}, fmt.Errorf("RemoveEstablishment: %w", err)
	}
	return &pb.ResponseDelete{AddressId: aID}, nil
}

func (s *EstablishmentService) Create(ctx context.Context, in *pb.RequestCreate) (*pb.Response, error) {
	m := model.Establishment{AddressID: in.Establishment.AddressId}
	err := controller.CreateEstablishment(&m)
	if err != nil {
		return &pb.Response{}, err
	}
	return &pb.Response{Id: uint64(m.ID)}, nil
}

func (s *EstablishmentService) Get(ctx context.Context, in *pb.RequestById) (*pb.Establishment, error) {
	e, err := controller.GetEstablishmentByID(uint(in.Id))
	if err != nil {
		return &pb.Establishment{}, fmt.Errorf("GetEstablishmentById: %w", err)
	}
	q := len(e.Tables)
	return &pb.Establishment{Id: uint64(e.ID), AddressId: e.AddressID, Quantity: uint32(q)}, nil
}
