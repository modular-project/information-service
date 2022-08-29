package controller

import (
	"information-service/model"
	"information-service/storage"
)

func createTable(id uint) model.Table {
	return model.Table{EstablishmentID: id, UserID: 0}
}

func AddTableToEstablishment(id uint) ([]uint64, error) {
	return IncreaseQuantityTablesInEstablishment(id, 1)
}

func RemoveTableFromEstablishment(establishmentID, quantity uint) (uint32, error) {
	if quantity == 0 {
		return 0, nil
	}
	m := model.Establishment{}
	if establishmentID < 1 || storage.DB().Where("id = ?", establishmentID).First(&m).RowsAffected == 0 {
		return 0, ErrEstablishmentNotFound
	}
	tables := []model.Table{}
	res := storage.DB().Where("establishment_id = ? AND user_id = 0", establishmentID).Limit(int(quantity)).Find(&tables)
	if res.Error != nil {
		return 0, res.Error
	}
	if res.RowsAffected == 0 {
		return 0, ErrEstablishmentHasNoTables
	}
	res = storage.DB().Delete(tables)
	return uint32(res.RowsAffected), res.Error
}

func IncreaseQuantityTablesInEstablishment(id uint, quantity int) ([]uint64, error) {
	if quantity < 1 {
		return nil, ErrQuantityMustBeGreaterThanZero
	}
	if storage.DB().Where("id = ?", id).First(&model.Establishment{}).RowsAffected == 0 {
		return nil, ErrEstablishmentNotFound
	}
	m := make([]model.Table, quantity)
	for i := 0; i < quantity; i++ {
		m[i] = createTable(id)
	}
	res := storage.DB().CreateInBatches(m, int(quantity))
	if res.Error != nil {
		return nil, res.Error
	}
	if res.RowsAffected == 0 {
		return nil, ErrNoRowsAffected
	}
	ids := make([]uint64, res.RowsAffected)
	for i := range ids {
		ids[i] = uint64(m[i].ID)
	}
	return ids, nil
}

func GetTablesQuantityInEstablishment(id uint) (uint, error) {
	var count int64
	err := storage.DB().Model(&model.Table{}).Where("establishment_id = ?", id).Count(&count).Error
	return uint(count), err
}

func GetTablesInEstablishment(id uint) ([]model.Table, error) {
	m := []model.Table{}
	err := storage.DB().Where("establishment_id = ?", id).Find(&m).Error
	return m, err
}

func ChangeTableStatusById(userID uint, establishment_id uint, tableId uint) error {
	m := model.Table{}
	if storage.DB().Where("id = ?", tableId).First(&m).RowsAffected == 0 {
		return ErrTableNotFound
	}
	if m.EstablishmentID != establishment_id {
		return unauthorizedErr("the table is in another establishment")
	}
	if m.UserID != 0 && m.UserID != userID {
		return ErrTableNotAvailable
	}
	if m.UserID == 0 {
		return storage.DB().Model(&m).Update("user_id", userID).Error
	}
	return storage.DB().Model(&m).Update("user_id", userID).Error
}
