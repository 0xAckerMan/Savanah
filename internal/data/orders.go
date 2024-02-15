package data

import (
	"time"

	"gorm.io/gorm"
)

type Order struct {
	gorm.Model
	CustomerID  uint      `json:"customer_id" gorm:"not null"`
	ProductID   uint      `json:"product_id" gorm:"not null"`
	Quantity    int       `json:"quantity" gorm:"not null"`
	OrderStatus string    `json:"order_status" gorm:"default:'placed'"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Version     int       `json:"version" gorm:"version; default:1"`

	Customer *Customer `gorm:"foreignKey:CustomerID"`
	Product  *Product  `gorm:"foreignKey:ProductID"`
}
