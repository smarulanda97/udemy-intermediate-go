package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/smarulanda97/app-stripe/internal/middlewares"
)

// Create application routes
func (app *application) routes() http.Handler {
	mux := chi.NewRouter()

	mux.Use(middlewares.SessionLoadMiddleware(session))

	mux.Handle("/static/*", app.kernel.FileServer())

	mux.Get("/", app.HomeController)

	mux.Route("/cart", func(r chi.Router) {
		r.Post("/payment", app.PaymentController)
		r.Get("/receipt", app.ReceiptController)
	})

	mux.Route("/terminal", func(r chi.Router) {
		r.Get("/payment", app.TerminalController)
		r.Get("/receipt", app.ReceiptController)

	})

	mux.Route("/products", func(r chi.Router) {
		r.Get("/{id}", app.ProductController)
	})

	mux.Route("/plan", func(r chi.Router) {
		r.Get("/bronze", app.BronzePlanController)
	})

	mux.Route("/user", func(r chi.Router) {
		r.Get("/login", app.LoginController)
	})

	return mux
}
