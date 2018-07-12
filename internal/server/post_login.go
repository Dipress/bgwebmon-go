package server

import (
	"net/http"

	"github.com/dipress/bgwebmon-go/internal/server/user"
	"github.com/julienschmidt/httprouter"
)

// PostLogin func return login page
func postLogin(s user.Authenticator) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		okResp := user.Response{}

		if err := s.Authenticate(r, &okResp); err != nil {
			switch resp := err.(type) {
			case user.ErrorResponse:
				errorResponse(w, resp)
			default:
				internalServerResponse(w)
			}
			return
		}
		okResponse(w, okResp)
	}
}
