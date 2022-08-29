package controller

import (
	"fmt"
	"information-service/model"
	"information-service/storage"
)

func CreateProduct(m *model.Product) error {
	m.CreatedAt = nil
	m.UpdatedAt = nil
	m.ID = 0
	return storage.DB().Omit("deleted_at").Create(m).Error
}

func DeleteProductById(id uint) error {
	res := storage.DB().Where("id = ? AND deleted_at IS NULL", id).Delete(&model.Product{})
	if res.Error != nil {
		return fmt.Errorf("%w, %v", unauthorizedErr("error asd"), res.Error)
	}
	if res.RowsAffected == 0 {
		return unauthorizedErr("no rows")
	}
	return nil
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

func GetProductsInBatch(ids []uint64) ([]model.Product, error) {
	m := []model.Product{}
	err := storage.DB().Unscoped().Find(&m, ids).Error
	return m, err
}
