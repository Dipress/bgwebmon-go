package server

import (
	"encoding/json"
	"net/http"
)

const (
	internalServerErrorMessage = "Internal server error"
)

func okResponse(w http.ResponseWriter, resp interface{}) {
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

func errorResponse(w http.ResponseWriter, resp interface{}) {
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)
	json.NewEncoder(w).Encode(resp)
}

func internalServerResponse(w http.ResponseWriter) {
	http.Error(w, internalServerErrorMessage, http.StatusInternalServerError)
}
