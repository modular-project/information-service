package model

import (
	"time"

	"gorm.io/gorm"
)

type Model struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
type Product struct {
	Model
	Name        string
	Price       float64
	Description string
	Url         string
}

type Establishment struct {
	Model
	Tables []Table
	Name   string
}

type Table struct {
	Model
	EstablishmentID uint
	UserID          uint
}
