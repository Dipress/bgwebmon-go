package server

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	internalServerErrorMessage = "internal server error"
)

func okResponse(w http.ResponseWriter, resp interface{}) {
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		internalServerErrorResponse(w)
	}
	w.WriteHeader(http.StatusOK)
}

func errorResponse(w http.ResponseWriter, resp interface{}) {
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		internalServerErrorResponse(w)
	}
	w.WriteHeader(http.StatusBadRequest)
}

func internalServerErrorResponse(w http.ResponseWriter) {
	w.WriteHeader(http.StatusInternalServerError)
	fmt.Fprint(w, internalServerErrorMessage)
}
