package models

import (
	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	Name     string  `json:"name" gorm:"unique"`
	Price    float64 `json:"price"`
	Quantity int     `json:"quantity"`
}
