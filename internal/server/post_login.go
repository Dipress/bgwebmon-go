package server

import (
	"net/http"

	"github.com/dipress/bgwebmon-go/internal/server/user"
	"github.com/julienschmidt/httprouter"
)

const (
	internalServerErrorMessage = "internal server error"
)

// PostLogin func return login page
func postLogin(s user.Authenticator) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		okResp := user.Response{}

		if err := s.Authenticate(r, &okResp); err != nil {
			switch resp := err.(type) {
			case user.ValidationErrorResponse:
				errorResponse(w, resp)
			case user.ErrorResponse:
				errorResponse(w, resp)
			default:
				internalServerResponse(w, internalServerErrorMessage)
			}
			return
		}
		okResponse(w, okResp)
	}
}
