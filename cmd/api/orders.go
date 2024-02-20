package main

import (
	"net/http"

	"github.com/0xAckerMan/Savanah/internal/data"
	"github.com/0xAckerMan/Savanah/internal/validator"
)

func (app *Application) handle_getOrders(w http.ResponseWriter, r *http.Request) {
	var orders []*data.Order
	order := app.DB.Preload("Product").First(&orders)
	if order.Error != nil {
        app.writeJSON(w,http.StatusOK,envelope{"orders": "No orders available at the moment"}, nil)
		return
	}

    if order.RowsAffected == 0{
        app.errorResponse(w,r,http.StatusOK, "no orders at the moment")
        return
    }

    err := app.writeJSON(w,http.StatusOK, envelope{"orders": orders},nil)
    if err != nil{
        app.serverErrorResponse(w,r,err)
        return
    }
}

func (app *Application) handle_getSingleOrder(w http.ResponseWriter, r *http.Request) {
    id, err := app.readIDParam(r)
    if err != nil {
        app.notFoundResponse(w, r)
        return
    }
    var order data.Order
    err = app.DB.First(&order, id).Error
    if err != nil {
        app.notFoundResponse(w, r)
        return
    }
    err = app.writeJSON(w, http.StatusOK, envelope{"order": &order}, nil)
    if err != nil {
        app.serverErrorResponse(w, r, err)
        return
    }
}

//customer create order
func (app *Application) handle_createOrder(w http.ResponseWriter, r *http.Request) {
    var input struct {
        CustomerID uint `json:"customer_id"`
        ProductID  uint `json:"product_id"`
        Quantity   int `json:"quantity"`
    }

    err := app.readJSON(w, r, &input)
    if err != nil {
        app.badRequestResponse(w, r, err)
        return
    }

    order := &data.Order{
        CustomerID: input.CustomerID,
        ProductID:  input.ProductID,
        Quantity:   input.Quantity,
    }

    v := validator.New()

    if data.ValidateOrder(v, order); !v.Valid() {
        app.failedValidationResponse(w, r, v.Errors)
        return
    }

    err = app.DB.Create(&order).Error
    if err != nil {
        app.serverErrorResponse(w, r, err)
        return
    }

    err = app.writeJSON(w, http.StatusCreated, envelope{"order": order}, nil)
    if err != nil {
        app.serverErrorResponse(w, r, err)
        return
    }
}

func (app *Application) handle_updateOrder(w http.ResponseWriter, r *http.Request) {
    id, err := app.readIDParam(r)
    if err != nil {
        app.notFoundResponse(w, r)
        return
    }

    var order data.Order

    if app.DB.First(&order, id).Error != nil{
        app.notFoundResponse(w,r)
        return
    }
    var input struct {
        CustomerID *uint `json:"customer_id"`
        ProductID  *uint `json:"product_id"`
        Quantity   *int `json:"quantity"`
    }

    err = app.readJSON(w, r, &input)
    if err != nil {
        app.badRequestResponse(w, r, err)
        return
    }

    if input.CustomerID != nil {
        order.CustomerID = *input.CustomerID
    }

    if input.ProductID != nil {
        order.ProductID = *input.ProductID
    }

    if input.Quantity != nil {
        order.Quantity = *input.Quantity
    }

    v := validator.New()

    if data.ValidateOrder(v, &order); !v.Valid() {
        app.failedValidationResponse(w, r, v.Errors)
        return
    }

    err = app.DB.Save(&order).Error
    if err != nil {
        app.serverErrorResponse(w, r, err)
        return
    }

    err = app.writeJSON(w, http.StatusOK, envelope{"order": &order}, nil)
    if err != nil {
        app.serverErrorResponse(w, r, err)
        return
    }
}

func (app *Application) handle_deleteOrder(w http.ResponseWriter, r *http.Request) {
    id, err := app.readIDParam(r)
    if err != nil {
        app.notFoundResponse(w, r)
        return
    }

    var order data.Order
    result := app.DB.Delete(&order, id)
    if result.Error != nil{
        app.serverErrorResponse(w,r,result.Error)
        return
    }

    if result.RowsAffected == 0{
        app.notFoundResponse(w,r)
        return
    }

    err = app.writeJSON(w,http.StatusOK,envelope{"message": "Deleted successfuly"},nil)
    if err != nil{
        app.serverErrorResponse(w,r,err)
        return
    }
}