package data

import "gorm.io/gorm"

type Product struct {
    gorm.Model
    ProductName string `json:"product_name" gorm:"unique,not null"`
    Description string `json:"description"`
    Price       int    `json:"price" gorm:"not null"`
    Version     int    `json:"version" gorm:"version; default:1"`
}
