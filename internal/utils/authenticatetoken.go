package utils

import (
	"errors"
	"net/http"
	"strings"

	"github.com/smarulanda97/app-stripe/internal/models"
)

func AuthenticateToken(r *http.Request, dbm models.DBModels) (*models.User, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return nil, errors.New("no authorization header received")
	}

	headerParts := strings.Split(authHeader, " ")
	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		return nil, errors.New("no authorization header received")
	}

	token := headerParts[1]
	if len(token) != 26 {
		return nil, errors.New("authentication token wrong size")
	}

	user, err := dbm.GetUserByToken(token)
	if err != nil {
		return nil, errors.New("no matching user found")
	}

	return user, nil
}
