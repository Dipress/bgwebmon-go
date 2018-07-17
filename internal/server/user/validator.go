package user

import (
	"github.com/go-ozzo/ozzo-validation"
)

const (
	loginFieldName    = "login"
	passwordFieldName = "password"
)

type validator interface {
	Validate(model Model) error
}

type ozzoValidator struct{}

func (v ozzoValidator) Validate(model Model) error {
	validatationError := ValidationErrorResponse{Status: "error"}

	if err := validation.Validate(model.Login, validation.Required); err != nil {
		validatationError.Fields = append(validatationError.Fields,
			validationField{Field: loginFieldName, Message: err.Error()},
		)
	}

	if err := validation.Validate(model.Password, validation.Required); err != nil {
		validatationError.Fields = append(validatationError.Fields,
			validationField{Field: passwordFieldName, Message: err.Error()},
		)
	}

	if len(validatationError.Fields) > 0 {
		return validatationError
	}
	return nil
}
