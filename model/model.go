package model

import (
	"time"

	"gorm.io/gorm"
)

type Model struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt *time.Time     `json:"created_at"`
	UpdatedAt *time.Time     `json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

type Product struct {
	Model
	Name        string  `json:"name,omitempty"`
	Price       float64 `json:"price,omitempty"`
	Description string  `json:"description,omitempty"`
	Url         string  `json:"url,omitempty"`
	BaseID      uint
}

type Establishment struct {
	Model
	AddressID string
	Tables    []Table
}

type Table struct {
	Model
	EstablishmentID uint
	UserID          uint
}
