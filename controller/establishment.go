package controller

import (
	"information-service/model"
	"information-service/storage"
)

func CreateEstablishment(m *model.Establishment) error {
	return storage.DB().Create(m).Error
}

func UpdateEstablishmentData(m *model.Establishment) error {
	if m.ID == 0 {
		return ErrEstablishmentNotFound
	}
	return storage.DB().Model(&model.Establishment{}).Select("Name").Updates(m).Error
}

func RemoveEstablishment(m *model.Establishment) error {
	if m.ID == 0 {
		return ErrEstablishmentNotFound
	}
	// don't remove if it have pending orders
	return storage.DB().Delete(m).Error
}

func GetEstablishmentByID(id uint) (model.Establishment, error) {
	if id == 0 {
		return model.Establishment{}, ErrEstablishmentNotFound
	}
	m := model.Establishment{}
	res := storage.DB().Unscoped().Where("id = ?", id).First(&m)
	if res.RowsAffected == 0 {
		return model.Establishment{}, ErrEstablishmentNotFound
	}
	return m, res.Error
}

func GetEstablishmentsAvailable() ([]model.Establishment, error) {
	m := []model.Establishment{}
	err := storage.DB().Find(&m).Error
	return m, err
}
