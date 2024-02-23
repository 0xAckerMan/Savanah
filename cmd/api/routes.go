package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (app *Application) routes() chi.Router {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.MethodNotAllowed(app.methodNotAllowedResponse)

	r.Get("/", app.Home)
	r.Get("/api/sessions/oauth/google", app.GoogleOauth)
	r.Route("/api/v1", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Use(app.DeserializeCustomer)

			r.Get("/customers/me", app.GetMe)
			r.Get("/customers/", app.handle_getCustomers)

			//logout/signout
			r.Get("/logout", app.SignOutCustomer)

			r.Route("/orders", func(r chi.Router) {
				r.Get("/{id}", app.handle_getSingleOrder)
				r.Post("/", app.handle_createOrder)
				r.Patch("/{id}", app.handle_updateOrder)
			})

			r.Route("/admin", func(r chi.Router) {
				r.Use(app.AdminAuthMiddleware)
				r.Get("/orders/", app.handle_getOrders)
				r.Get("/customers", app.handle_getCustomers)
				r.Post("/products", app.handle_createProduct)
				r.Patch("/products/{id}", app.handle_updateProduct)
				r.Delete("/products/{id}", app.handle_deleteProduct)
				r.Get("/customers/{id}", app.handle_getSingleCustomer)
			})
		})


		r.Get("/healthcheck", app.healthcheck)
		r.Route("/products", func(r chi.Router) {
			r.Get("/", app.handle_getProduct)
			r.Get("/{id}", app.handle_getSingleProduct)
		})


		
			r.Post("/register", app.SignUpCustomer)
			r.Post("/login", app.SignInCustomer)

	})

	return r
}
