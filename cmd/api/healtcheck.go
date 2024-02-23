package main

import (
	"encoding/json"
	"net/http"
)

func (app *Application) healthcheck(w http.ResponseWriter, r *http.Request){
	status := map[string]interface{}{
		"status": "active",
		"environment": app.config.env,
		"version": "1.0.0",
	}

	js, err := json.Marshal(status)

	if err != nil {
		app.logger.Println(err)
		http.Error(w, "Internal Server Error", 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(js)
}

func (app *Application) Home(w http.ResponseWriter, r *http.Request){
	message := map[string]interface{}{
		"message": "Welcome to Savanah API",
	}

	app.writeJSON(w, http.StatusOK, message, nil)
}