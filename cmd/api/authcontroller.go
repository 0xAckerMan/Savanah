package main

import (
	"net/http"
	"os"
	"time"

	"github.com/0xAckerMan/Savanah/cmd/utils"
	"github.com/0xAckerMan/Savanah/internal/data"
	"github.com/0xAckerMan/Savanah/internal/validator"
	"golang.org/x/crypto/bcrypt"
)

func (app *Application) SignUpCustomer(w http.ResponseWriter, r *http.Request) {
	var input struct {
		FirstName   string `json:"first_name"`
		LastName    string `json:"last_name"`
		Email       string `json:"email"`
		PhoneNumber string `json:"phone_number"`
		Password    string `json:"password"`
		IsAdmin     bool   `json:"is_admin"`
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
		IsAdmin:     input.IsAdmin,
		Verified:    true,
	}

	v := validator.New()
	if data.ValidateUser(v, customer); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), 12)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	customer.Password = string(hash)

	result := app.DB.Create(&customer)
	if result.Error != nil && result.Error.Error() == "Error 1062: Duplicate entry" {
		app.failedValidationResponse(w, r, map[string]string{"email": "email already exists"})
		return
	} else if result.Error != nil {
		app.serverErrorResponse(w, r, result.Error)
		return
	}

	err = app.writeJSON(w, http.StatusCreated, envelope{"customer": customer}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

func (app *Application) SignInCustomer(w http.ResponseWriter, r *http.Request) {
	var input data.LoginCustomer

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	var customer data.Customer

	result := app.DB.First(&customer, "email = ?", input.Email)
	if result.Error != nil {
		app.serverErrorResponse(w, r, result.Error)
		return
	}

	if customer.Provider == "Google" {
		app.failedValidationResponse(w, r, map[string]string{"email": "email already exists"})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(customer.Password), []byte(input.Password))
	if err != nil {
		app.failedValidationResponse(w, r, map[string]string{"password": "invalid password"})
		return
	}

	var ttl = time.Hour * 24 * 14
	token, err := utils.GenerateToken(ttl, customer.ID, os.Getenv("JWT_SECRET"))
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	//set cookie with the token
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    token,
		Expires:  time.Now().Add(ttl),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteNoneMode,
		Path:     "/",
		MaxAge:   int(ttl.Seconds()),
		Domain:   "localhost",
	})

	err = app.writeJSON(w, http.StatusOK, envelope{"customer": customer}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	//redirect to the home page
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *Application) SignOutCustomer(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    "",
		Expires:  time.Now(),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteNoneMode,
		Path:     "/",
		MaxAge:   -1,
		Domain:   "localhost",
	})

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *Application) GoogleOauth(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")

	var pathUrl string = "/"
	if r.URL.Query().Get("state") != "" {
		pathUrl = r.URL.Query().Get("state")
	}

	if code == "" {
		app.failedValidationResponse(w, r, map[string]string{"code": "invalid code"})
		return
	}

	tokenRes, err := utils.GetGoogleOAuthToken(code)

	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	google_user, err := utils.GetGoogleUserInfo(tokenRes.AccessToken, tokenRes.IdToken)

	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	customer_data := data.Customer{
		FirstName: google_user.Given_name,
		LastName:  google_user.Family_name,
		Email:     google_user.Email,
		Provider:  "Google",
		Verified:  true,
	}

	if app.DB.Model(&customer_data).Where("email = ?", customer_data.Email).Updates(&customer_data).RowsAffected == 0 {
		result := app.DB.Create(&customer_data)
		if result.Error != nil {
			app.serverErrorResponse(w, r, result.Error)
			return
		}
	}

	var customer data.Customer
	result := app.DB.First(&customer, "email = ?", customer_data.Email)
	if result.Error != nil {
		app.serverErrorResponse(w, r, result.Error)
		return
	}

	var ttl = time.Hour * 24 * 14
	token, err := utils.GenerateToken(ttl, customer.ID, os.Getenv("JWT_SECRET"))
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	//set cookie with the token
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    token,
		Expires:  time.Now().Add(ttl),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteNoneMode,
		Path:     "/",
		MaxAge:   int(ttl.Seconds()),
		Domain:   "localhost",
	})

	http.Redirect(w, r, pathUrl, http.StatusSeeOther)

}

func (app *Application) GetMe(w http.ResponseWriter, r *http.Request) {
	currentCustomer := r.Context().Value("customer").(data.Customer)

	err := app.writeJSON(w, http.StatusOK, envelope{"customer": currentCustomer}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}
