package handler

import (
	"context"
	"fmt"
	"information-service/controller"

	pb "github.com/modular-project/protobuffers/information/table"
)

type TableService struct {
	pb.UnimplementedTableServiceServer
}

func (service *TableService) AddTable(ctx context.Context, in *pb.RequestById) (*pb.ResponseAdd, error) {
	ids, err := controller.AddTableToEstablishment(uint(in.Id))
	if err != nil {
		return &pb.ResponseAdd{}, err
	}
	return &pb.ResponseAdd{Ids: ids}, nil
}
func (service *TableService) AddTables(ctx context.Context, in *pb.RequestAdd) (*pb.ResponseAdd, error) {
	ids, err := controller.IncreaseQuantityTablesInEstablishment(uint(in.Id), int(in.Quantity))
	if err != nil {
		return &pb.ResponseAdd{}, err
	}
	return &pb.ResponseAdd{Ids: ids}, nil
}
func (service *TableService) GetFromEstablishment(ctx context.Context, in *pb.RequestById) (*pb.ResponseGetAll, error) {
	ms, err := controller.GetTablesInEstablishment(uint(in.Id))
	if err != nil {
		return &pb.ResponseGetAll{}, err
	}
	res := make([]*pb.Table, len(ms))
	for i, m := range ms {
		res[i] = &pb.Table{Id: uint64(m.ID), EstablishmenId: uint64(m.EstablishmentID), UserId: uint64(m.UserID)}
	}
	return &pb.ResponseGetAll{Tables: res}, nil
}
func (service *TableService) ChangeStatus(ctx context.Context, in *pb.Table) (*pb.ResponseStatus, error) {
	err := controller.ChangeTableStatusById(uint(in.UserId), uint(in.EstablishmenId), uint(in.Id))
	if err != nil {
		return &pb.ResponseStatus{}, err
	}
	return &pb.ResponseStatus{}, nil
}

func (service *TableService) RemoveFromEstablishment(ctx context.Context, in *pb.RequestDelete) (*pb.ResponseDelete, error) {
	deleted, err := controller.RemoveTableFromEstablishment(uint(in.EstablishmenId), uint(in.Quantity))
	if err != nil {
		return &pb.ResponseDelete{}, fmt.Errorf("removeFromEstablishment: %w", err)
	}
	return &pb.ResponseDelete{Deleted: deleted}, nil
}
