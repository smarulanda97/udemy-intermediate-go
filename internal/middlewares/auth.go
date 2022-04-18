package middlewares

import (
	"net/http"

	"github.com/smarulanda97/app-stripe/internal/models"
	"github.com/smarulanda97/app-stripe/internal/utils"
)

func AuthMiddleware(dbm models.DBModels) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, err := utils.AuthenticateToken(r, dbm)
			if err != nil {
				utils.InvalidCredentials(w, r)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
