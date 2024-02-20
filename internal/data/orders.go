package data

import (
	"github.com/0xAckerMan/Savanah/internal/validator"
	"gorm.io/gorm"
)

type Order struct {
	gorm.Model
	CustomerID  uint      `json:"customer_id" gorm:"not null"`
	ProductID   uint      `json:"product_id" gorm:"not null"`
	Quantity    int       `json:"quantity" gorm:"not null"`
	OrderStatus string    `json:"order_status" gorm:"default:'placed'"`
	Version     int       `json:"version" gorm:"version; default:1"`

	Customer *Customer `gorm:"foreignKey:CustomerID"`
	Product  *Product  `gorm:"foreignKey:ProductID"`
}

func ValidateOrder(v *validator.Validator, order *Order) {
	v.Check(order.Quantity > 0, "quantity", "must be greater than zero")
	v.Check(order.OrderStatus == "placed" || order.OrderStatus == "purchased", "order_status", "must be either 'placed' or 'purchased'")
	v.Check(order.Version > 0, "version", "must be greater than zero")
	v.Check(order.CustomerID > 0, "customer_id", "must be a positive integer")
	v.Check(order.ProductID > 0, "product_id", "must be a positive integer")
}