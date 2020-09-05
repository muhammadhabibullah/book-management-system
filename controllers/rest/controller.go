package rest

import (
	"encoding/json"
	"net/http"

	"book-management-system/controllers/rest/responses"
)

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, responses.ErrorResponse{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_, _ = w.Write(response)
}
