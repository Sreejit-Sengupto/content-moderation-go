package validator

import "github.com/go-playground/validator/v10"

func Validtor() *validator.Validate {
	var validate *validator.Validate
	validate = validator.New(validator.WithRequiredStructEnabled())
	return validate
}
