package utils

import "net/http"

func InvalidCredentials(w http.ResponseWriter, r *http.Request) error {
	var payload Payload

	payload.Error = true
	payload.Message = "invalid authentication credentials"

	err := WriteJson(w, r, http.StatusUnauthorized, payload)
	if err != nil {
		return err
	}

	return nil
}
