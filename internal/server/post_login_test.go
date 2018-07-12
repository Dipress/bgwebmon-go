package server

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dipress/bgwebmon-go/internal/server/user"
	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/assert"
)

func Test_postLogin(t *testing.T) {
	tests := []struct {
		name     string
		authFunc authenticatorFunc
		code     int
	}{
		{
			name: "ok",
			authFunc: func(r *http.Request, m *user.Response) error {
				return nil
			},
			code: http.StatusOK,
		},
		{
			name: "login or password mismatch",
			authFunc: func(r *http.Request, m *user.Response) error {
				return user.ErrorResponse{Status: "error", Message: "login or password a wrong"}
			},
			code: http.StatusUnauthorized,
		},
		{
			name: "internal server error",
			authFunc: func(r *http.Request, m *user.Response) error {
				return errors.New("mock error")
			},
			code: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			req := httptest.NewRequest("POST", "http://any-host/auth", nil)
			res := httptest.NewRecorder()

			f := postLogin(authenticatorFunc(tt.authFunc))
			f(res, req, make(httprouter.Params, 0))

			assert.Equal(t, tt.code, res.Code)
		})
	}
}

type authenticatorFunc func(*http.Request, *user.Response) error

func (f authenticatorFunc) Authenticate(r *http.Request, m *user.Response) error {
	return f(r, m)
}
