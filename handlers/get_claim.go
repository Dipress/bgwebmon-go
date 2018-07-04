package handlers

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// GetClaim func return login page
func GetClaim(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "This is a Claoim Page!\n")
}
