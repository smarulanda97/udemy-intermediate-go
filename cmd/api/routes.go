package main

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/smarulanda97/app-stripe/internal/middlewares"
)

func (app *application) routes() http.Handler {
	mux := chi.NewRouter()

	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Accept", "Authorization", "X-CSRF-Token"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	mux.Route("/api", func(r chi.Router) {
		r.Route("/payment", func(mux chi.Router) {
			mux.Get("/intent", app.GetPaymentIntentController)
		})

		r.Route("/product", func(mux chi.Router) {
			mux.Get("/{id}", app.GetProductController)
		})

		r.Route("/auth", func(mux chi.Router) {
			mux.Post("/login", app.CreateAuthTokenController)
			mux.Post("/check", app.CheckAuthTokenController)
		})

		r.Route("/admin", func(mux chi.Router) {
			mux.Use(middlewares.AuthMiddleware(app.DB))

			mux.Get("/test", func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte("got in"))
			})

			mux.Route("/terminal", func(mux chi.Router) {
				mux.Post("/receipt", app.ApiTerminalController)
			})
		})
	})

	return mux
}
