package main

import (
    "errors"
	"fmt"
	"net/http"
)

var (
    ErrRecordNotFound = errors.New("record searched could not be found")
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

func (app *Application) errDuplicateUser(w http.ResponseWriter, r *http.Request){
    message := "Sorry their are some conflic and duplicated data"
    app.errorResponse(w,r,http.StatusConflict, message)
}
