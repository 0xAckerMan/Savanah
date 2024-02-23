package main

import (
	"context"
	"net/http"
	"os"
	"strings"

	"github.com/0xAckerMan/Savanah/cmd/utils"
	"github.com/0xAckerMan/Savanah/internal/data"
)
func (app *Application) DeserializeCustomer(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        var token string
        cookie, err := r.Cookie("token")

        authorizationHeader := r.Header.Get("Authorization")
        fields := strings.Fields(authorizationHeader)

        if len(fields) != 0 && fields[0] == "Bearer" {
            token = fields[1]
        } else if err == nil {
            token = cookie.Value
        }

        if token == "" {
            app.UnAuthorizedResponse(w, r)
            return
        }

        sub, err := utils.ValidateToken(token, os.Getenv("JWT_SECRET"))
        if err != nil {
            app.UnAuthorizedResponse(w, r)
            return
        }

        var customer data.Customer
        result := app.DB.First(&customer, "id = ?", sub)
        if result.Error != nil {
            app.serverErrorResponse(w, r, result.Error)
            return
        }

        // Set current user to user context
        ctx := r.Context()
        ctx = context.WithValue(ctx, "Customer", customer)
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}

func (app *Application) AdminAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		customer := r.Context().Value("customer").(data.Customer)
		if !customer.IsAdmin {
			app.UnAuthorizedResponse(w, r)
			return
		}
		next.ServeHTTP(w, r)
	})
}