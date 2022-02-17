package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	FirstName string `json:"first_name" gorm:"type:varchar(128);not null"`
	LastName  string `json:"last_name" gorm:"type:varchar(128);not null"`
	Email     string `json:"email" gorm:"type:varchar(128);not null;unique"`
	Password  string `json:"password"`
}
