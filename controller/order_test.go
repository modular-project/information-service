package controller

import (
	"information-service/model"
	"information-service/storage"
	"testing"

	pb "github.com/modular-project/protobuffers/order/order"
)

var TestConfigDB storage.DBConnection = storage.DBConnection{
	TypeDB:   storage.POSTGRESQL,
	User:     "admin_restaurant",
	Password: "RestAuraNt_pgsql.561965697",
	Host:     "localhost",
	Port:     "5433",
	NameDB:   "testing",
}

func createProducts(size int) {
	for i := 0; i < size; i++ {
		m := model.Product{Name: "not updated"}
		CreateProduct(&m)
	}
}

func createEstablishment(size int) {
	for i := 0; i < size; i++ {
		m := model.Establishment{}
		CreateEstablishment(&m)
	}

}

func TestValidateOrder(t *testing.T) {
	storage.NewDB(&TestConfigDB)
	models := []interface{}{model.Product{}, model.Establishment{}, model.Table{}}
	err := storage.DB().AutoMigrate(models...)
	if err != nil {
		t.Fatalf("Failed to Create tables: %s", err)
	}
	t.Cleanup(func() { dropsTables(t, models...) })
	createEstablishment(2)
	createProducts(6)
	AddTableToEstablishment(1)
	AddTableToEstablishment(2)
	AddTableToEstablishment(1)
	testCase := []struct {
		in      pb.Order
		wantErr bool
	}{
		{
			in: pb.Order{EstablishmentId: 1, Type: &pb.Order_LocalOrder{LocalOrder: &pb.LocalOrder{TableId: 1, EmployeeId: 1}},
				OrderProducts: []*pb.OrderProduct{
					{ProductId: 1, Quantity: 3},
					{ProductId: 2, Quantity: 1},
					{ProductId: 5, Quantity: 8},
				}},
		}, {
			in: pb.Order{EstablishmentId: 1, Type: &pb.Order_RemoteOrder{RemoteOrder: &pb.RemoteOrder{AddressId: "", UserId: 2}},
				OrderProducts: []*pb.OrderProduct{
					{ProductId: 1, Quantity: 3},
					{ProductId: 2, Quantity: 1},
				}},
		}, {
			in: pb.Order{EstablishmentId: 1, Type: &pb.Order_LocalOrder{LocalOrder: &pb.LocalOrder{TableId: 1, EmployeeId: 2}},
				OrderProducts: []*pb.OrderProduct{
					{ProductId: 1, Quantity: 3},
				}},
			wantErr: true,
		}, {
			in: pb.Order{EstablishmentId: 1, Type: &pb.Order_LocalOrder{LocalOrder: &pb.LocalOrder{TableId: 2, EmployeeId: 1}},
				OrderProducts: []*pb.OrderProduct{
					{ProductId: 1, Quantity: 3},
				}},
			wantErr: true,
		}, {
			in: pb.Order{EstablishmentId: 1, Type: &pb.Order_LocalOrder{LocalOrder: &pb.LocalOrder{TableId: 3, EmployeeId: 1}},
				OrderProducts: []*pb.OrderProduct{
					{ProductId: 8, Quantity: 3},
				}},
			wantErr: true,
		}, {
			in: pb.Order{EstablishmentId: 1, Type: &pb.Order_RemoteOrder{RemoteOrder: &pb.RemoteOrder{AddressId: "", UserId: 2}},
				OrderProducts: []*pb.OrderProduct{
					{ProductId: 8, Quantity: 3},
				}},
			wantErr: true,
		}, {
			in: pb.Order{EstablishmentId: 3, Type: &pb.Order_RemoteOrder{RemoteOrder: &pb.RemoteOrder{AddressId: "", UserId: 2}},
				OrderProducts: []*pb.OrderProduct{
					{ProductId: 1, Quantity: 3},
				}},
		}, {in: pb.Order{EstablishmentId: 1, Type: &pb.Order_LocalOrder{LocalOrder: &pb.LocalOrder{TableId: 1, EmployeeId: 1}},
			OrderProducts: []*pb.OrderProduct{
				{ProductId: 1, Quantity: 3},
			}},
		}, {
			in: pb.Order{EstablishmentId: 1, Type: &pb.Order_RemoteOrder{RemoteOrder: &pb.RemoteOrder{AddressId: "", UserId: 2}},
				OrderProducts: []*pb.OrderProduct{
					{ProductId: 2, Quantity: 1},
				}},
		}, {
			in: pb.Order{EstablishmentId: 1, Type: &pb.Order_RemoteOrder{RemoteOrder: &pb.RemoteOrder{AddressId: "", UserId: 2}},
				OrderProducts: []*pb.OrderProduct{
					{ProductId: 2, Quantity: 1},
				}},
		},
	}
	for _, tc := range testCase {
		_, err := ValidateOrder(&tc.in)
		// if !errors.Is(tc.err, err) {
		// 	t.Logf("%+v", tc.in)
		// 	t.Errorf("got error: %s, want error: %s", err, tc.err)
		// }
		if (err != nil) != tc.wantErr {
			t.Errorf("checkProducts() error = %v, wantErr %v", err, tc.wantErr)
			t.Logf("%+v", tc.in)
			return
		}
	}
}

func insertProducts(t *testing.T) {
	products := []model.Product{
		{
			Price: 10.50,
		}, {
			Price: 21,
		}, {
			Price: 15,
		}, {
			Price: 33.33,
		},
	}
	r := storage.DB().Create(&products)
	if r.RowsAffected != 4 {
		t.Fatalf("no rows afected %s", r.Error)
	}
	if r.Error != nil {
		t.Fatalf("create: %s", r.Error)
	}
}

func dropsTables(t *testing.T, tables ...interface{}) {
	err := storage.DB().Migrator().DropTable(tables...)
	if err != nil {
		t.Fatalf("Failed to clean database: %s", err)
	}
}

func Test_checkProducts(t *testing.T) {
	if err := storage.NewDB(&TestConfigDB); err != nil {
		t.Fatalf("start db: %s", err)
	}
	models := []interface{}{model.Product{}, model.Establishment{}, model.Table{}}
	err := storage.DB().AutoMigrate(models...)
	if err != nil {
		t.Fatalf("Failed to Create tables: %s", err)
	}
	t.Cleanup(func() { dropsTables(t, models...) })
	insertProducts(t)
	type args struct {
		ps []*pb.OrderProduct
	}
	tests := []struct {
		name    string
		args    args
		want    float64
		wantErr bool
	}{
		{
			name: "Product repeated OK",
			args: args{ps: []*pb.OrderProduct{
				{
					ProductId: 1,
					Quantity:  2,
				}, {
					ProductId: 4,
					Quantity:  3,
				}, {
					ProductId: 3,
					Quantity:  5,
				}, {
					ProductId: 1,
					Quantity:  3,
				}, {
					ProductId: 3,
					Quantity:  1,
				},
			},
			},
			want: 242.49,
		}, {
			name: "not found",
			args: args{ps: []*pb.OrderProduct{
				{
					ProductId: 1,
					Quantity:  2,
				}, {
					ProductId: 4,
					Quantity:  3,
				}, {
					ProductId: 3,
					Quantity:  5,
				}, {
					ProductId: 6,
					Quantity:  3,
				}, {
					ProductId: 3,
					Quantity:  1,
				},
			},
			},
			want:    0,
			wantErr: true,
		},
	}
	storage.NewDB(&TestConfigDB)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CheckProducts(tt.args.ps)
			if (err != nil) != tt.wantErr {
				t.Errorf("checkProducts() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("checkProducts() = %v, want %v", got, tt.want)
			}
		})
	}
}
