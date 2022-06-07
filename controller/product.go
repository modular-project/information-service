package controller

import (
	"information-service/model"
	"information-service/storage"
	"time"
)

func CreateProduct(m *model.Product) error {
	m.CreatedAt = time.Time{}
	m.UpdatedAt = time.Time{}
	m.ID = 0
	return storage.DB().Omit("deleted_at").Create(m).Error
}

func DeleteProductById(id uint) error {
	return storage.DB().Where("id = ?", id).Delete(&model.Product{}).Error
}

func UpdateProductById(id uint, m *model.Product) error {
	if storage.DB().Where("id = ?", id).First(&model.Product{}).RowsAffected == 0 {
		return ErrProductNotFound
	}
	err := CreateProduct(m)
	if err != nil {
		return err
	}
	return DeleteProductById(id)
}

func GetProductById(id uint) (model.Product, error) {
	m := model.Product{}
	err := storage.DB().Unscoped().Where("id = ?", id).First(&m).Error
	return m, err
}

func GetAllProducts() ([]model.Product, error) {
	m := []model.Product{}
	err := storage.DB().Find(&m).Error
	return m, err
}
