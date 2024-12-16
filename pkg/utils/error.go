package utils

import (
	"encoding/json"
	"net/http"
	"tt/internal/models"
)

func SendError(w http.ResponseWriter, message string, statusCode int) {
	w.WriteHeader(statusCode)
	var err models.Error
	err.SetError(message)
	json.NewEncoder(w).Encode(err)
}
