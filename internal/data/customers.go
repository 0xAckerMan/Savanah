package data

import (
	"gorm.io/gorm"
	"github.com/0xAckerMan/Savanah/internal/validator"
)

type Customer struct {
	gorm.Model
	FirstName   string `json:"first_name" gorm:"not null"`
	LastName    string `json:"last_name" gorm:"not null"`
	Email       string `json:"email" gorm:"uniqueIndex;not null"`
	PhoneNumber string `json:"phone_number" gorm:"uniqueIndex;not null"`
	Password    string `json:"-" gorm:"not null"`
	Version     int    `json:"-" gorm:"version; default:1"`
    IsAdmin     bool   `json:"-" gorm:"default:false"`
	Verified  bool      `json:"-" gorm:"default:false;"`
	Provider  string    `json:"-" gorm:"default:'local';"`

	Orders []*Order `gorm:"foreignKey:CustomerID;constraint:OnDelete:CASCADE"`
}

type LoginCustomer struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func ValidateEmail(v *validator.Validator, email string) {
	v.Check(email != "", "email", "must be provided")
	v.Check(validator.Matches(email, *validator.EmailRX), "email", "must be a valid email address")
}

func ValidatePassword(v *validator.Validator, password string) {
	v.Check(password != "", "password", "must be provided")
	v.Check(len(password) >= 8, "password", "must be at least 8 bytes long")
}

func ValidateUser(v *validator.Validator, user *Customer) {
	v.Check(user.FirstName != "", "first_name", "must be provided")
	v.Check(user.LastName != "", "last_name", "must be provided")

	ValidateEmail(v, user.Email)
	if user.Password != "" {
		ValidatePassword(v, user.Password)
	}

	v.Check(user.PhoneNumber != "", "phone_number", "must be provided")
	v.Check(len(user.PhoneNumber) == 13, "phone_number", "must be 10 digits long")

	// v.Check(user.Provider == "local" || user.Provider == "google", "provider", "invalid provider")
}
