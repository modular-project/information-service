package controller

import (
	"fmt"
	"information-service/model"
	"information-service/storage"

	pb "github.com/modular-project/protobuffers/order/order"
)

func CheckProducts(ps []*pb.OrderProduct) (float64, error) {
	var total float64
	size := len(ps)
	ids := make([]uint, size)
	for i, p := range ps {
		ids[i] = uint(p.ProductId)
	}
	m := []model.Product{}
	res := storage.DB().Select("id", "price").Where(ids).Find(&m)
	if res.Error != nil {
		return 0, fmt.Errorf("find products in batch: %w", res.Error)
	}
	mp := make(map[uint]float64)
	for _, p := range m {
		mp[p.ID] = p.Price
	}
	for i := range ps {
		p, ok := mp[uint(ps[i].ProductId)]
		if !ok {
			return 0, fmt.Errorf("product %d not found", ps[i].ProductId)
		}
		total += p * float64(ps[i].Quantity)
	}
	return total, nil
}

func checkLocalOrder(o *pb.Order) error {
	rows := storage.DB().Where("id = ? ", o.EstablishmentId).Select("id").First(&model.Establishment{}).RowsAffected
	if rows == 0 {
		return ErrEstablishmentNotFound
	}
	table := model.Table{}
	res := storage.DB().Where("id = ? AND establishment_id = ?", o.GetLocalOrder().TableId, o.EstablishmentId).First(&table)
	if res.RowsAffected == 0 {
		return ErrTableNotFound
	}
	if res.Error != nil {
		return res.Error
	}
	userID := uint(o.GetLocalOrder().EmployeeId)
	if table.UserID != userID {
		if table.UserID != 0 {
			return ErrTableIsInUse
		}
		storage.DB().Model(&model.Table{}).Where("id = ?", table.ID).Update("user_id", userID)
	}
	return nil
}

func ValidateOrder(o *pb.Order) (float64, error) {
	// TODO MOVE
	switch o.Type.(type) {
	case *pb.Order_LocalOrder:
		err := checkLocalOrder(o)
		if err != nil {
			return 0, err
		}
	}
	t, err := CheckProducts(o.OrderProducts)
	if err != nil {
		return 0, fmt.Errorf("checkProducts: %w", err)
	}
	return t, nil
}
