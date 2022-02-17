package models

import (
	"gorm.io/gorm"
)

type OrderItem struct {
	gorm.Model
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price"`
	OrderID   uint    `json:"order_id"`
	ProductID int     `json:"product_id"`
	Product   Product
}

type Order struct {
	gorm.Model
	Quantity   int         `json:"quantity"`
	Price      float64     `json:"price"`
	UserID     int         `json:"user_id"`
	User       User        `gorm:"foreignkey:UserID"`
	OrderItems []OrderItem `gorm:"foreignkey:OrderID"`
}
