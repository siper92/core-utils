package core_utils

import (
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
)

func SimpleStructValidation(s interface{}) error {
	return ValidateWithCustomFunc(s, map[string]validator.Func{})
}

func ValidateWithCustomFunc(s interface{}, validators map[string]validator.Func) error {
	validate := validator.New()
	for field, fn := range validators {
		err := validate.RegisterValidation(field, fn)
		if err != nil {
			ErrorWarning(err)
		}
	}

	err := validate.Struct(s)

	if err != nil {
		var invalidValidationError *validator.InvalidValidationError
		if errors.As(err, &invalidValidationError) {
			ErrorWarning(err)
		}
	}

	return err
}

func HasFieldError(field string, err error) bool {
	if err == nil {
		return false
	}

	if errors.As(err, &validator.ValidationErrors{}) {
		for _, err := range err.(validator.ValidationErrors) {
			if err.Field() == field {
				return true
			}
		}
	}

	return false
}

func ValidateAndWarn(s interface{}) error {
	_errors := SimpleStructValidation(s)
	if errors.As(_errors, &validator.ValidationErrors{}) {
		for _, err := range _errors.(validator.ValidationErrors) {
			FieldErrorWarning(err)
		}
	}

	return _errors
}

func FieldErrorWarning(err validator.FieldError) {
	var errMessage string
	if err.Tag() == "required" {
		errMessage = fmt.Sprintf("%s: %s", err.Tag(), err.StructNamespace())
	} else if err.Tag() == "oneof" {
		errMessage = fmt.Sprintf("%s: %s\n >>> %v alowed: %s", err.Tag(), err.StructNamespace(), err.Value(), err.Param())
	} else {
		errMessage = fmt.Sprintf("%s: %s\n >>> Value: %v", err.Tag(), err.StructNamespace(), err.Value())
	}

	Warning(errMessage)
}
