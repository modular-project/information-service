package controller

import (
	"fmt"
	"information-service/model"
	"information-service/storage"
)

func GetByAddress(aID string) (uint, uint32, error) {
	e := model.Establishment{}
	res := storage.DB().Model(&model.Establishment{}).Preload("Tables").Where("address_id = ?", aID).First(&e)
	if res.Error != nil {
		return 0, 0, fmt.Errorf("first establishment by address: %w", res.Error)
	}
	if res.RowsAffected == 0 {
		return 0, 0, fmt.Errorf("no rows affected")
	}
	return e.ID, uint32(len(e.Tables)), nil
}

func CreateEstablishment(m *model.Establishment) error {
	return storage.DB().Create(m).Error
}

func UpdateEstablishmentData(m *model.Establishment) error {
	if m.ID == 0 {
		return ErrEstablishmentNotFound
	}
	return storage.DB().Model(&model.Establishment{}).Select("Name").Updates(m).Error
}

func RemoveEstablishment(m *model.Establishment) (string, error) {
	if m.ID == 0 {
		return "", ErrEstablishmentNotFound
	}
	// don't remove if it have pending orders
	if err := storage.DB().Select("AddressID").First(m).Error; err != nil {
		return "", fmt.Errorf("find establishment: %w", err)
	}
	if err := storage.DB().Delete(m).Error; err != nil {
		return "", fmt.Errorf("delete establishment: %w", err)
	}
	return m.AddressID, nil
}

func GetEstablishmentByID(id uint) (model.Establishment, error) {
	if id == 0 {
		return model.Establishment{}, ErrEstablishmentNotFound
	}
	m := model.Establishment{}

	res := storage.DB().Model(&model.Establishment{}).Preload("Tables").Unscoped().Where("id = ?", id).First(&m)
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
