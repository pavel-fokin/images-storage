package httputil

import (
	"fmt"
	"encoding/json"
	"net/http"

)

type Error struct {
	Message string `json:"message"`
}

// ErrorResponse
type ErrorResponse struct {
	Data struct {
		Errors []Error `json:"errors"`
	} `json:"data"`
}

func AsErrorResponse(
	w http.ResponseWriter, err error, statusCode int,
) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	// make error response
	payload := ErrorResponse{}
	payload.Data.Errors = []Error{{Message: fmt.Sprint(err)}}

	// encode json
	json.NewEncoder(w).Encode(payload)
}
