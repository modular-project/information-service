package controller_test

import (
	"errors"
	"information-service/controller"
	"information-service/model"
	"information-service/storage"
	"testing"

	pb "github.com/modular-project/protobuffers/order/order"
)

func TestValidateOrder(t *testing.T) {
	storage.New(storage.TESTING)
	models := []interface{}{model.Product{}, model.Establishment{}, model.Table{}}
	err := storage.DB().AutoMigrate(models...)
	if err != nil {
		t.Fatalf("Failed to Create tables: %s", err)
	}
	t.Cleanup(func() { dropsTables(t, models...) })
	createEstablishment(2)
	createProducts(6)
	controller.AddTableToEstablishment(1)
	controller.AddTableToEstablishment(2)
	controller.AddTableToEstablishment(1)
	testCase := []struct {
		in  pb.Order
		err error
	}{
		{
			pb.Order{EstablishmentId: 1, Type: &pb.Order_LocalOrder{LocalOrder: &pb.LocalOrder{TableId: 1, EmployeeId: 1}},
				OrderProducts: []*pb.OrderProduct{
					{ProductId: 1, Quantity: 3},
					{ProductId: 2, Quantity: 1},
					{ProductId: 5, Quantity: 8},
				}},
			nil,
		}, {
			pb.Order{EstablishmentId: 1, Type: &pb.Order_RemoteOrder{RemoteOrder: &pb.RemoteOrder{AddressId: 2, UserId: 2}},
				OrderProducts: []*pb.OrderProduct{
					{ProductId: 1, Quantity: 3},
					{ProductId: 2, Quantity: 1},
				}},
			nil,
		}, {
			pb.Order{EstablishmentId: 1, Type: &pb.Order_LocalOrder{LocalOrder: &pb.LocalOrder{TableId: 1, EmployeeId: 2}},
				OrderProducts: []*pb.OrderProduct{
					{ProductId: 1, Quantity: 3},
				}},
			controller.ErrTableIsInUse,
		}, {
			pb.Order{EstablishmentId: 1, Type: &pb.Order_LocalOrder{LocalOrder: &pb.LocalOrder{TableId: 2, EmployeeId: 1}},
				OrderProducts: []*pb.OrderProduct{
					{ProductId: 1, Quantity: 3},
				}},
			controller.ErrTableNotFound,
		}, {
			pb.Order{EstablishmentId: 1, Type: &pb.Order_LocalOrder{LocalOrder: &pb.LocalOrder{TableId: 3, EmployeeId: 1}},
				OrderProducts: []*pb.OrderProduct{
					{ProductId: 8, Quantity: 3},
				}},
			controller.ErrProductNotFound,
		}, {
			pb.Order{EstablishmentId: 1, Type: &pb.Order_RemoteOrder{RemoteOrder: &pb.RemoteOrder{AddressId: 2, UserId: 2}},
				OrderProducts: []*pb.OrderProduct{
					{ProductId: 8, Quantity: 3},
				}},
			controller.ErrProductNotFound,
		}, {
			pb.Order{EstablishmentId: 3, Type: &pb.Order_RemoteOrder{RemoteOrder: &pb.RemoteOrder{AddressId: 2, UserId: 2}},
				OrderProducts: []*pb.OrderProduct{
					{ProductId: 1, Quantity: 3},
				}},
			controller.ErrEstablishmentNotFound,
		}, {pb.Order{EstablishmentId: 1, Type: &pb.Order_LocalOrder{LocalOrder: &pb.LocalOrder{TableId: 1, EmployeeId: 1}},
			OrderProducts: []*pb.OrderProduct{
				{ProductId: 1, Quantity: 3},
			}},
			nil,
		}, {
			pb.Order{EstablishmentId: 1, Type: &pb.Order_RemoteOrder{RemoteOrder: &pb.RemoteOrder{AddressId: 3, UserId: 2}},
				OrderProducts: []*pb.OrderProduct{
					{ProductId: 2, Quantity: 1},
				}},
			nil,
		}, {
			pb.Order{EstablishmentId: 1, Type: &pb.Order_RemoteOrder{RemoteOrder: &pb.RemoteOrder{AddressId: 2, UserId: 2}},
				OrderProducts: []*pb.OrderProduct{
					{ProductId: 2, Quantity: 1},
				}},
			nil,
		},
	}
	for _, tc := range testCase {
		err := controller.ValidateOrder(&tc.in)
		if !errors.Is(tc.err, err) {
			t.Logf("%+v", tc.in)
			t.Errorf("got error: %s, want error: %s", err, tc.err)
		}
	}
}
