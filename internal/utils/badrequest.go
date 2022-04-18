package utils

import (
	"encoding/json"
	"net/http"
)

func BadRequest(w http.ResponseWriter, r *http.Request, err error) error {
	var payload Payload

	payload.Error = true
	payload.Message = err.Error()

	out, err := json.MarshalIndent(payload, "", "\t")
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(out)

	return nil
}
