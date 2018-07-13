package server

import (
	"encoding/json"
	"net/http"
)

func okResponse(w http.ResponseWriter, resp interface{}) {
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

func errorResponse(w http.ResponseWriter, resp interface{}) {
	w.WriteHeader(http.StatusUnauthorized)
	json.NewEncoder(w).Encode(resp)
}

func internalServerResponse(w http.ResponseWriter, resp interface{}) {
	w.WriteHeader(http.StatusInternalServerError)
	json.NewEncoder(w).Encode(resp)
}
