package server

import (
	"net/http"

	"github.com/dipress/bgwebmon-go/internal/server/auth"
	"github.com/julienschmidt/httprouter"
)

const (
	internalServerErrorMessage = "internal server error"
)

// PostLogin func return login page
func postLogin(s auth.Authenticator) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		okResp := auth.Response{}

		if err := s.Authenticate(r, &okResp); err != nil {
			switch resp := err.(type) {
			case auth.ValidationErrorResponse:
				errorResponse(w, resp)
			case auth.ErrorResponse:
				errorResponse(w, resp)
			default:
				internalServerErrorResponse(w, internalServerErrorMessage)
			}
			return
		}
		okResponse(w, okResp)
	}
}
