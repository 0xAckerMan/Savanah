package main

import (
	"net/http"
)

func (app *Application) healthcheck(w http.ResponseWriter, r *http.Request){
	status := map[string]interface{}{
		"status": "active",
		"environment": app.config.env,
		"version": "1.0.0",
	}

	err := app.writeJSON(w, http.StatusOK, envelope{"health":status},nil)
	if err != nil {
		app.serverErrorResponse(w,r,err)
		return
	}
}

func (app *Application) Home(w http.ResponseWriter, r *http.Request){
	message := map[string]interface{}{
		"message": "Welcome to Savanah API",
	}

	app.writeJSON(w, http.StatusOK, message, nil)
}