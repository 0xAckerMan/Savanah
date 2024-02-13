package data

import "gorm.io/gorm"

type Product struct {
	gorm.Model
	ProductName string `json:"product_name"`
	Description string `json:"description"`
	Price       int    `json:"price"`
	Version     int    `json:"version" gorm:"default:1"`
}