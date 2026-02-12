package utils

import (
	"encoding/json"
	"net/http"

	"github.com/loop/backend/rider-auth/rest/internals/models"
)

func respondWithJSON(w http.ResponseWriter, statusCode int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(payload)
}

func respondWithError(w http.ResponseWriter, statusCode int, message string, errorDetail string) {
	errResp := models.ErrorResponse{
		Success: false,
		Message: message,
		Status:  int64(statusCode),
		Error:   errorDetail,
	}
	respondWithJSON(w, statusCode, errResp)
}
