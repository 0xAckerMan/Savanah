package main

import (
	"errors"
	"net/http"

	"github.com/0xAckerMan/Savanah/internal/data"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func (app *Application) handle_getCustomers(w http.ResponseWriter, r *http.Request) {
	var customers []*data.Customer
	err := app.writeJSON(w, http.StatusOK, envelope{"customers": customers}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

func (app *Application) handle_getSingleCustomer(w http.ResponseWriter, r *http.Request) {
	var customer *data.Customer

	err := app.writeJSON(w, http.StatusOK, envelope{"customer": customer}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

func (app *Application) handle_createCustomer(w http.ResponseWriter, r *http.Request) {
	var input struct {
		FirstName   string `json:"first_name"`
		LastName    string `json:"last_name"`
		Email       string `json:"email"`
		PhoneNumber string `json:"phone_number"`
		Password    string `json:"password"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	customer := &data.Customer{
		FirstName:   input.FirstName,
		LastName:    input.LastName,
		Email:       input.Email,
		PhoneNumber: input.PhoneNumber,
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), 12)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	customer.Password = string(hash)

	result := app.DB.Create(&customer)
	if result.Error != nil {
		switch {
		case errors.Is(result.Error, gorm.ErrDuplicatedKey):
			app.errDuplicateUser(w,r)
		default:
			app.serverErrorResponse(w, r, err)
			return
		}
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"Customer": customer}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}
