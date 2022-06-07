package controller_test

import (
	"errors"
	"information-service/controller"
	"information-service/model"
	"information-service/storage"
	"testing"
)

func dropsTables(t *testing.T, tables ...interface{}) {
	err := storage.DB().Migrator().DropTable(tables...)
	if err != nil {
		t.Fatalf("Failed to clean database: %s", err)
	}
}

func createEstablishment(size int) {
	for i := 0; i < size; i++ {
		m := model.Establishment{Name: "test"}
		controller.CreateEstablishment(&m)
	}

}

func TestIncreaseTablesQuantityInBatchInEstablishment(t *testing.T) {
	storage.New(storage.TESTING)
	models := []interface{}{model.Establishment{}, model.Table{}}
	err := storage.DB().AutoMigrate(models...)
	if err != nil {
		t.Fatalf("Failed to Create tables: %s", err)
	}
	t.Cleanup(func() { dropsTables(t, models...) })
	type In struct {
		id       uint
		quantity float64
	}
	testCase := []struct {
		in            In
		err           error
		finalQuantity uint
	}{
		{In{1, 3}, nil, 3},
		{In{2, 0}, controller.ErrQuantityMustBeGreaterThanZero, 0},
		{In{1, -1}, controller.ErrQuantityMustBeGreaterThanZero, 0},
		{In{3, 5}, controller.ErrEstablishmentNotFound, 0},
		{In{1, 4}, nil, 7},
	}

	createEstablishment(2)
	for _, tc := range testCase {
		err := controller.IncreaseQuantityTablesInEstablishment(tc.in.id, int(tc.in.quantity))
		if !errors.Is(tc.err, err) {
			t.Errorf("error at increase quantity, got: %s, want: %s", err, tc.err)
		}
		if err == nil {
			t.Logf("%+v", tc)
			quantiy, _ := controller.GetTablesQuantityInEstablishment(tc.in.id)
			if quantiy != tc.finalQuantity {
				t.Errorf("Got %d quantity, want %d quantity", quantiy, tc.finalQuantity)
			}
		}
	}
}

func TestDecreaseTableInEstablishment(t *testing.T) {
	storage.New(storage.TESTING)
	models := []interface{}{model.Establishment{}, model.Table{}}
	err := storage.DB().AutoMigrate(models...)
	if err != nil {
		t.Fatalf("Failed to Create tables: %s", err)
	}
	t.Cleanup(func() { dropsTables(t, models...) })
	type In struct {
		id       uint
		quantity uint
	}
	testCase := []struct {
		in            In
		err           error
		finalQuantity uint
		deleted       uint32
	}{
		{In{id: 1, quantity: 0}, nil, 5, 0},
		{In{id: 0, quantity: 10}, controller.ErrEstablishmentNotFound, 0, 0},
		{In{id: 3, quantity: 10}, controller.ErrEstablishmentNotFound, 0, 0},
		{In{id: 1, quantity: 3}, nil, 3, 2},
		{In{id: 2, quantity: 5}, nil, 0, 3},
		{In{id: 2, quantity: 5}, controller.ErrEstablishmentHasNoTables, 0, 0},
	}

	createEstablishment(2)
	controller.IncreaseQuantityTablesInEstablishment(1, 5)
	controller.IncreaseQuantityTablesInEstablishment(2, 3)
	for i := 1; i <= 3; i++ {
		storage.DB().Model(&model.Table{}).Where("id = ?", i).Update("user_id", 1)
	}
	ms, _ := controller.GetTablesInEstablishment(1)
	for _, m := range ms {
		t.Logf("%+v", m)
	}
	for _, tc := range testCase {
		t.Logf("%+v", tc)
		deleted, err := controller.RemoveTableFromEstablishment(tc.in.id, tc.in.quantity)
		if !errors.Is(tc.err, err) {
			t.Errorf("error at decrease quantity, got: %s, want: %s", err, tc.err)
		}
		if err == nil {
			quantiy, _ := controller.GetTablesQuantityInEstablishment(tc.in.id)
			if quantiy != tc.finalQuantity {
				t.Errorf("Got %d quantity, want %d quantity", quantiy, tc.finalQuantity)
			}
			if deleted != tc.deleted {
				t.Errorf("Got %d deleted rows, want %d deleted rows", deleted, tc.deleted)
			}
		}
	}
}
