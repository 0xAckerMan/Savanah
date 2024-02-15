package main

import (
	"net/http"

	"github.com/0xAckerMan/Savanah/internal/data"
)

func (app *Application) handle_getProduct(w http.ResponseWriter, r *http.Request) {
	var Products []data.Product
	product := app.DB.Find(&Products)
    if product.RowsAffected == 0 {
        app.writeJSON(w,http.StatusOK,envelope{"response": "No products available"},nil)
        return
    }
	if product.Error != nil {
		app.serverErrorResponse(w, r, ErrRecordNotFound)
		return
	}

	err := app.writeJSON(w, http.StatusOK, envelope{"Products": Products}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

func (app *Application) handle_getSingleProduct(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}
	var product data.Product
	err = app.DB.First(&product, id).Error
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}
	err = app.writeJSON(w, http.StatusOK, envelope{"product": &product}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

func (app *Application) handle_createProduct(w http.ResponseWriter, r *http.Request) {
	var input struct {
		ProductName string `json:"product_name"`
		Description string `json:"description"`
		Price       int    `json:"price"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	Product := &data.Product{
		ProductName: input.ProductName,
		Description: input.Description,
		Price:       input.Price,
	}

	result := app.DB.Create(&Product)
	if result.Error != nil {
		app.serverErrorResponse(w, r, result.Error)
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"product": Product}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

func (app *Application) handle_updateProduct(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	var product data.Product

	if app.DB.First(&product, id).Error != nil{
        app.notFoundResponse(w,r)
        return
    }
	var input struct {
		ProductName *string `json:"product_name"`
		Description *string `json:"description"`
		Price       *int    `json:"price"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if input.ProductName != nil {
		product.ProductName = *input.ProductName
	}

	if input.Description != nil {
		product.Description = *input.Description
	}

	if input.Price != nil {
		product.Price = *input.Price
	}

    product.Version ++

    result := app.DB.Save(&product)
    if result.Error != nil{
        app.serverErrorResponse(w,r,result.Error)
        return
    }

    err = app.writeJSON(w, http.StatusOK, envelope{"product": &product}, nil)
    if err != nil{
        app.serverErrorResponse(w,r,err)
        return
    }
}

func (app *Application) handle_deleteProduct(w http.ResponseWriter, r *http.Request){
    id, err := app.readIDParam(r)
    if err != nil{
        app.notFoundResponse(w,r)
        return
    }

    var product *data.Product

    result := app.DB.Delete(&product, id)
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
