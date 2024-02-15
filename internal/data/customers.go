package data

import "gorm.io/gorm"

type Customer struct {
	gorm.Model
	FirstName   string `json:"first_name" gorm:"not null"`
	LastName    string `json:"last_name" gorm:"not null"`
	Email       string `json:"email" gorm:"unique,not null"`
	PhoneNumber string `json:"phone_number" gorm:"not null"`
	Password    string `json:"password" gorm:"not null"`
	Version     int    `json:"version" gorm:"version; default:1"`

	Orders []*Order `gorm:"foreignKey:CustomerID;constraint:OnDelete:CASCADE"`
}
