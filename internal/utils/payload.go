package utils

import "github.com/smarulanda97/app-stripe/internal/models"

type Payload struct {
	Error   bool          `json:"error"`
	Message string        `json:"message"`
	Token   *models.Token `json:"authentication_token"`
}
