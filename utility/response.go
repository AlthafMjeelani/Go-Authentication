// utils/response.go
package utils

import (
	"encoding/json"
	"net/http"
)

// Response defines the structure for both success and error messages
type Response struct {
	Status  bool        `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// WriteJSONResponse is a common utility function to send JSON responses
func WriteJSONResponse(w http.ResponseWriter, statusCode int, status bool, message string, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(Response{
		Status:  status,
		Message: message,
		Data:    data,
	})
}
