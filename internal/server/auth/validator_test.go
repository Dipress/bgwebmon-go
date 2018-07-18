package auth

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ozzoValidatorValidate(t *testing.T) {
	tests := []struct {
		name              string
		model             Model
		wantErr           bool
		validatationError ValidationErrorResponse
	}{
		{
			name: "empty login",
			model: Model{
				Password: "password",
			},
			wantErr: true,
			validatationError: ValidationErrorResponse{
				Status: "error",
				Fields: []validationField{
					validationField{
						Field:   "login",
						Message: "cannot be blank",
					},
				},
			},
		},
		{
			name: "empty password",
			model: Model{
				Login: "login",
			},
			wantErr: true,
			validatationError: ValidationErrorResponse{
				Status: "error",
				Fields: []validationField{
					validationField{
						Field:   "password",
						Message: "cannot be blank",
					},
				},
			},
		},
		{
			name:    "both empty",
			model:   Model{},
			wantErr: true,
			validatationError: ValidationErrorResponse{
				Status: "error",
				Fields: []validationField{
					validationField{
						Field:   "login",
						Message: "cannot be blank",
					},
					validationField{
						Field:   "password",
						Message: "cannot be blank",
					},
				},
			},
		},
		{
			name: "valid",
			model: Model{
				Login:    "login",
				Password: "password",
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			validator := ozzoValidator{}
			err := validator.Validate(tt.model)

			if tt.wantErr {
				assert.NotNil(t, err)
				assert.Equal(t, tt.validatationError, err.(ValidationErrorResponse))
			}

			if !tt.wantErr {
				assert.NoError(t, err)
			}
		})
	}
}
