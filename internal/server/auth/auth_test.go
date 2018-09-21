package auth

import (
	"bytes"
	"errors"
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

func Test_authenticator_Authenticate(t *testing.T) {
	tests := []struct {
		name         string
		decodeFunc   func(r *http.Request, model *Model) error
		validateFunc func(model Model) error
		findFunc     func(login string, model *Model) error
		checkFunc    func(requestPassword, modelPassword string) error
		tokenFunc    func(userID int, secretKey string) (string, error)
		expect       Response
		wantErr      bool
	}{
		{
			name: "decoder error",
			decodeFunc: func(r *http.Request, model *Model) error {
				return errors.New("mock error")
			},
			wantErr: true,
		},
		{
			name: "validator error",
			decodeFunc: func(r *http.Request, model *Model) error {
				return nil
			},
			validateFunc: func(model Model) error {
				return errors.New("mock error")
			},
			wantErr: true,
		},
		{
			name: "finder error",
			decodeFunc: func(r *http.Request, model *Model) error {
				return nil
			},
			validateFunc: func(model Model) error {
				return nil
			},
			findFunc: func(login string, model *Model) error {
				return errors.New("mock error")
			},
			wantErr: true,
		},
		{
			name: "checker error",
			decodeFunc: func(r *http.Request, model *Model) error {
				return nil
			},
			validateFunc: func(model Model) error {
				return nil
			},
			findFunc: func(login string, model *Model) error {
				return nil
			},
			checkFunc: func(requestPassword, modelPassword string) error {
				return errors.New("mock error")
			},
			wantErr: true,
		},
		{
			name: "tokener error",
			decodeFunc: func(r *http.Request, model *Model) error {
				return nil
			},
			validateFunc: func(model Model) error {
				return nil
			},
			findFunc: func(login string, model *Model) error {
				return nil
			},
			checkFunc: func(requestPassword, modelPassword string) error {
				return nil
			},
			tokenFunc: func(userID int, secretKey string) (string, error) {
				return "", errors.New("mock error")
			},
			wantErr: true,
		},
		{
			name: "ok",
			decodeFunc: func(r *http.Request, model *Model) error {
				return nil
			},
			validateFunc: func(model Model) error {
				return nil
			},
			findFunc: func(login string, model *Model) error {
				return nil
			},
			checkFunc: func(requestPassword, modelPassword string) error {
				return nil
			},
			tokenFunc: func(userID int, secretKey string) (string, error) {
				return "mytoken", nil
			},
			expect: Response{
				Status: "ok",
				Data: dataField{
					Token: "mytoken",
				},
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			auth := authenticator{
				secretKey: "secret",
				decoder:   decoderFunc(tt.decodeFunc),
				validator: validatorFunc(tt.validateFunc),
				finder:    finderFunc(tt.findFunc),
				checker:   checkerFunc(tt.checkFunc),
				tokener:   tokenerFunc(tt.tokenFunc),
			}

			var got Response
			err := auth.Authenticate(nil, &got)

			if tt.wantErr {
				assert.Error(t, err)
			}

			if !tt.wantErr {
				assert.Nil(t, err)
				assert.Equal(t, tt.expect, got)
			}

		})
	}
}

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

type decoderFunc func(r *http.Request, model *Model) error

func (d decoderFunc) Decode(r *http.Request, model *Model) error {
	return d(r, model)
}

type validatorFunc func(model Model) error

func (v validatorFunc) Validate(model Model) error {
	return v(model)
}

type finderFunc func(login string, model *Model) error

func (f finderFunc) Find(login string, model *Model) error {
	return f(login, model)
}

type checkerFunc func(requestPassword, modelPassword string) error

func (c checkerFunc) Check(requestPassword, modelPassword string) error {
	return c(requestPassword, requestExample)
}

type tokenerFunc func(userID int, secretKey string) (string, error)

func (t tokenerFunc) Token(userID int, secretKey string) (string, error) {
	return t(userID, secretKey)
}
