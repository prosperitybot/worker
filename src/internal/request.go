package internal

import (
	"encoding/json"
	"net/http"
)

func RespondToRequest(w http.ResponseWriter, status int, body interface{}) error {
	// Writes the status code passed through as the header for the response.
	w.WriteHeader(status)

	// Generates a json object based on the interface that is passed in and also specifies the status code in the response.
	return json.NewEncoder(w).Encode(body)
}
