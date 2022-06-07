package controller

import (
	"information-service/model"
	"information-service/storage"
	"log"

	pb "github.com/modular-project/protobuffers/order/order"
)

func checkProducts(ps []*pb.OrderProduct) error {
	size := len(ps)
	ids := make([]uint, size)
	for i, p := range ps {
		ids[i] = uint(p.ProductId)
	}
	log.Printf("%+v", ids)
	m := []model.Product{}
	res := storage.DB().Select("id").Where(ids).Find(&m)
	if size != int(res.RowsAffected) {
		return ErrProductNotFound
	}
	return res.Error
}

func checkLocalOrder(o *pb.Order) error {
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

func checkRemoteOrder(o *pb.Order) error {
	return nil
}

func ValidateOrder(o *pb.Order) error {
	rows := storage.DB().Where("id = ? ", o.EstablishmentId).Select("id").First(&model.Establishment{}).RowsAffected
	if rows == 0 {
		return ErrEstablishmentNotFound
	}
	switch o.Type.(type) {
	case *pb.Order_LocalOrder:
		err := checkLocalOrder(o)
		if err != nil {
			return err
		}
	case *pb.Order_RemoteOrder:
		err := checkRemoteOrder(o)
		if err != nil {
			return err
		}
	}
	return checkProducts(o.OrderProducts)
}
