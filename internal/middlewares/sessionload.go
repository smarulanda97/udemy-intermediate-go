package middlewares

import (
	"net/http"

	"github.com/alexedwards/scs/v2"
)

func SessionLoadMiddleware(session *scs.SessionManager) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return session.LoadAndSave(next)
	}
}
