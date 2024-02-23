package main

import (
	"errors"
	"net/http"

	"github.com/0xAckerMan/Savanah/internal/data"
	"github.com/0xAckerMan/Savanah/internal/validator"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func (app *Application) handle_getCustomers(w http.ResponseWriter, r *http.Request) {
	var customers []data.Customer
	result := app.DB.Preload("Orders").Preload("Orders.Product").Find(&customers)
	if result.Error != nil {
		app.serverErrorResponse(w, r, result.Error)
		return
	}
	if result.RowsAffected == 0 {
		app.writeJSON(w, http.StatusOK, envelope{"response": "No customers available"}, nil)
		return
	}

	// filter out admin users
	var filteredCustomers []data.Customer
	for _, customer := range customers {
		if !customer.IsAdmin {
			filteredCustomers = append(filteredCustomers, customer)
		}
	}

	err := app.writeJSON(w, http.StatusOK, envelope{"customers": filteredCustomers}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

func (app *Application) handle_getSingleCustomer(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}
	var customer data.Customer
	err = app.DB.Preload("Orders").Preload("Orders.Product").First(&customer, id).Error
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}
	err = app.writeJSON(w, http.StatusOK, envelope{"customer": &customer}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

func (app *Application) handle_updateCustomer(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	var input struct {
		FirstName   *string `json:"first_name"`
		LastName    *string `json:"last_name"`
		Email       *string `json:"email"`
		PhoneNumber *string `json:"phone_number"`
		Password    *string `json:"password"`
		IsAdmin     *bool   `json:"is_admin"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	customer := &data.Customer{Model: gorm.Model{ID: uint(id)}}
	result := app.DB.First(&customer)
	if result.Error != nil {
		app.notFoundResponse(w, r)
		return
	}

	if input.FirstName != nil {
		customer.FirstName = *input.FirstName
	}

	if input.LastName != nil {
		customer.LastName = *input.LastName
	}

	if input.Email != nil {
		customer.Email = *input.Email
	}

	if input.PhoneNumber != nil {
		customer.PhoneNumber = *input.PhoneNumber
	}

	if input.Password != nil {
		hash, err := bcrypt.GenerateFromPassword([]byte(*input.Password), 12)
		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}
		customer.Password = string(hash)
	}

	if input.IsAdmin != nil {
		customer.IsAdmin = *input.IsAdmin
	}

	v := validator.New()

	if data.ValidateUser(v, customer); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	result = app.DB.Save(&customer)
	if result.Error != nil {
		switch {
		case errors.Is(result.Error, gorm.ErrDuplicatedKey):
			app.errDuplicateUser(w,r)
		default:
			app.serverErrorResponse(w, r, err)
			return
		}
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"customer": customer}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}


func (app *Application) handle_deleteCustomer(w http.ResponseWriter, r *http.Request) {
	id,err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}
	customer := &data.Customer{Model: gorm.Model{ID: uint(id)}}
	result := app.DB.Delete(&customer)
	if result.Error != nil {
		app.serverErrorResponse(w, r, result.Error)
		return
	}
	err = app.writeJSON(w, http.StatusOK, envelope{"customer": customer}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

func (app *Application)handle_getMyOrders(w http.ResponseWriter, r *http.Request){
	var orders []*data.Order
	err := app.writeJSON(w, http.StatusOK, envelope{"orders": orders}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}
