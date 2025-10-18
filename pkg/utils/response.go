package utils

import (
	"encoding/json"
	"net/http"
)

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error      string      `json:"error"`
	StatusCode int         `json:"statusCode"`
	Details    interface{} `json:"details,omitempty"`
}

// RespondJSON sends a JSON response
func RespondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if data != nil {
		if err := json.NewEncoder(w).Encode(data); err != nil {
			Logger().WithError(err).Error("Failed to encode JSON response")
		}
	}
}

// RespondError sends an error response
func RespondError(w http.ResponseWriter, status int, message string) {
	RespondJSON(w, status, ErrorResponse{
		Error:      message,
		StatusCode: status,
	})
}

// RespondErrorWithDetails sends an error response with additional details
func RespondErrorWithDetails(w http.ResponseWriter, status int, message string, details interface{}) {
	RespondJSON(w, status, ErrorResponse{
		Error:      message,
		StatusCode: status,
		Details:    details,
	})
}
