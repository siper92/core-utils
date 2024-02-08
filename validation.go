package core_utils

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"strings"
)

type StructValidatorResult []validator.FieldError

func (s StructValidatorResult) HasFieldErrorFor(field string) bool {
	for _, err := range s {
		if err.Field() == field {
			return true
		}
	}

	return false
}

func (s StructValidatorResult) SingleError() error {
	var messages []string
	for _, err := range s {
		if val, ok := err.(validator.FieldError); ok {
			messages = append(messages, fmt.Sprintf("%s: %s", val.Field(), val.Tag()))
		}
	}

	return fmt.Errorf(strings.Join(messages, "\n"))
}

func ValidateGenericStruct(s interface{}) StructValidatorResult {
	validate := validator.New()
	err := validate.Struct(s)
	var result StructValidatorResult

	if err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			ErrorWarning(err)
		}

		if errors, ok := err.(validator.ValidationErrors); ok {
			return StructValidatorResult(errors)
		}
	}

	return result
}

func IsValidEmail(email string) bool {
	return ValidateGenericStruct(struct {
		Email string `validate:"required,email"`
	}{
		Email: email,
	}).HasFieldErrorFor("Email") == false
}
