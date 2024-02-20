package data

import (
	"github.com/0xAckerMan/Savanah/internal/validator"
	"gorm.io/gorm"
)

type Product struct {
    gorm.Model
    ProductName string `json:"product_name" gorm:"unique,not null"`
    Description string `json:"description"`
    Price       int    `json:"price" gorm:"not null"`
    Version     int    `json:"version" gorm:"version; default:1"`
}

func ValidateProduct(v *validator.Validator, product *Product) {
    v.Check(product.ProductName != "", "product_name", "must be provided")
    v.Check(len(product.ProductName) <= 50, "product_name", "must not be more than 500 bytes")
    v.Check(product.Description != "", "description", "must be provided")
    v.Check(len(product.Description) <= 600, "description", "must not be more than 1000 bytes")
    v.Check(product.Price > 0, "price", "must be greater than zero")
    v.Check(product.Version > 0, "version", "must be greater than zero")
}