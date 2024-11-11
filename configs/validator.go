package configs

import "github.com/go-playground/validator/v10"

type AppValidator struct {
	validator *validator.Validate
}

func (av *AppValidator) Validate(i interface{}) error {
	return av.validator.Struct(i)
}

func NewAppValidator() *AppValidator {
	return &AppValidator{
		validator: validator.New(),
	}
}
