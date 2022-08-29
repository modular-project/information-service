package controller_test

import (
	"errors"
	"information-service/controller"
	"information-service/model"
	"information-service/storage"
	"testing"
	"time"
)

func createProducts(size int) {
	for i := 0; i < size; i++ {
		m := model.Product{Name: "not updated"}
		controller.CreateProduct(&m)
	}
}

func equalsProducts(p1, p2 *model.Product) bool {
	return p1.Name == p2.Name && p1.Price == p2.Price && p1.Url == p2.Url && p1.Description == p2.Description
}

func TestUpdateProduct(t *testing.T) {
	storage.NewDB(&TestConfigDB)
	models := []interface{}{model.Product{}}
	err := storage.DB().AutoMigrate(models...)
	if err != nil {
		t.Fatalf("Failed to Create tables: %s", err)
	}
	t.Cleanup(func() { dropsTables(t, models...) })
	testCase := []struct {
		id  uint
		in  model.Product
		err error
	}{
		{1, model.Product{Name: "updated", Price: 10.45, Description: "updated", Url: "updated"}, nil},
		{1, model.Product{Name: "updated", Price: 10.45, Description: "updated", Url: "updated"}, controller.ErrProductNotFound},
		{2, model.Product{Model: model.Model{ID: 1}, Name: "updated 2", Price: 20.45, Description: "updated 2", Url: "updated 2"}, nil},
		{6, model.Product{Name: "updated 3", Price: 10.45, Description: "updated 3", Url: "updated 3"}, controller.ErrProductNotFound},
	}
	c := time.Now().Add(time.Hour)
	u := time.Now().Add(time.Hour)
	testCase[2].in.CreatedAt = &c
	testCase[2].in.UpdatedAt = &u
	size := 3
	now := time.Now().Add(time.Minute)
	createProducts(size)
	for _, tc := range testCase {
		err := controller.UpdateProductById(tc.id, &tc.in)
		t.Logf("%+v", tc.in)
		if !errors.Is(tc.err, err) {
			t.Errorf("got error: %s, want error: %s", err, tc.err)
		}
		if err == nil {
			m := model.Product{}
			storage.DB().Last(&m)
			if !equalsProducts(&m, &tc.in) {
				t.Errorf("got product: %+v, want %+v", m, tc.in)
			}
			if now.Before(*m.CreatedAt) {
				t.Errorf("got created at: %s", m.CreatedAt)
			}
			if now.Before(m.DeletedAt.Time) {
				t.Errorf("got deleted at: %s", m.DeletedAt.Time)
			}
			if !m.UpdatedAt.Equal(*m.CreatedAt) {
				t.Errorf("time updated and created are diferent")
			}
		}
	}
	var count int64
	storage.DB().Model(&model.Product{}).Count(&count)
	if int(count) != size {
		t.Errorf("got size: %d, want size: %d", count, size)
	}
}
