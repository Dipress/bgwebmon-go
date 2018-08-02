package auth

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	requestExample = `{
		"login": "login",
		"password": "password"
	}`
)

func Test_jwtToken_Token(t *testing.T) {
	tests := []struct {
		name      string
		userID    int
		secretKey string
		wantErr   bool
	}{
		{
			name:      "create token",
			userID:    27,
			secretKey: "secret",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			jwtToken := jwtToken{}
			_, err := jwtToken.Token(tt.userID, tt.secretKey)

			if !tt.wantErr {
				assert.NoError(t, err)
			}
		})
	}
}

func Test_passwordChecker_Check(t *testing.T) {
	tests := []struct {
		name            string
		requestPassword string
		modelPassword   string
		wantErr         bool
		errorResponse   ErrorResponse
	}{
		{
			name:            "password valid",
			requestPassword: "password",
			modelPassword:   "5F4DCC3B5AA765D61D8327DEB882CF99",
		},
		{
			name:            "password invalid",
			requestPassword: "password123",
			modelPassword:   "5F4DCC3B5AA765D61D8327DEB882CF99",
			wantErr:         true,
			errorResponse: ErrorResponse{
				Status:  "error",
				Message: invalidLoginOrPasswordMessage,
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			passwordChecker := passwordChecker{}
			err := passwordChecker.Check(tt.requestPassword, tt.modelPassword)

			if tt.wantErr {
				assert.NotNil(t, err)
				assert.Equal(t, tt.errorResponse, err.(ErrorResponse))
			}

			if !tt.wantErr {
				assert.NoError(t, err)
			}
		})
	}
}

func Test_bodyDecoder_Decode(t *testing.T) {
	tests := []struct {
		name    string
		request http.Request
		model   Model
		wantErr bool
	}{
		{
			name:    "decode request body",
			request: http.Request{},
			model:   Model{},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			bodyDecoder := bodyDecoder{}
			tt.request.Body = ioutil.NopCloser(bytes.NewReader([]byte(requestExample)))
			err := bodyDecoder.Decode(&tt.request, &tt.model)

			if !tt.wantErr {
				assert.NoError(t, err)
			}
		})
	}
}
