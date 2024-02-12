package main

import (
	"fmt"
	"net/http"
)

func (app *Application) logError (r *http.Request, err error){
    app.logger.Println(err)
}

func (app *Application) errorResponse (w http.ResponseWriter, r *http.Request, status int, message interface{}){
    err := app.writeJSON(w,status,envelope{"error": message},nil)
    if err != nil{
        app.logError(r, err)
        w.WriteHeader(500)
    }
}

func (app *Application) serverErrorResponse(w http.ResponseWriter,r *http.Request, err error){
    message := "Sorry, we cannot handle your request at the momment, internal server error"
    app.logError(r, err)
    app.errorResponse(w, r, http.StatusInternalServerError, message)
}

func (app *Application) notFoundResponse(w http.ResponseWriter, r *http.Request){
    message := "The requsted resource could not be found"
    app.errorResponse(w, r, http.StatusNotFound, message)
}

func (app *Application) methodNotAllowedResponse(w http.ResponseWriter, r * http.Request){
    message := fmt.Sprintf("Method %s is not allowed for this operation", r.Method)
    app.errorResponse(w, r, http.StatusMethodNotAllowed, message)
}
