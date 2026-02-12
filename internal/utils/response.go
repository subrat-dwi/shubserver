package utils

import (
	"encoding/json"
	"net/http"
	"strconv"
)

// JSON sends a JSON response with the given status code and data
func JSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(status)

	jsonData, _ := json.Marshal(data)
	w.Write(jsonData)
}

// Error sends a JSON error response with the given status code and message
func Error(w http.ResponseWriter, status int, message string) {
	JSON(w, status, map[string]string{
		"status":  strconv.Itoa(status),
		"message": message,
	})
}
