package validation

import "github.com/go-playground/validator/v10"

var validate = validator.New(validator.WithRequiredStructEnabled())

func Validate(s any) error {
	return validate.Struct(s)
}
