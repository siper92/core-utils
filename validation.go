package core_utils

import (
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"strings"
)

type StructValidatorResult []validator.FieldError

func (s StructValidatorResult) Error() string {
	var messages []string
	for _, err := range s {
		var val validator.FieldError
		if errors.As(err, &val) {
			messages = append(messages, fmt.Sprintf("%s: %s", val.Field(), val.Tag()))
		} else {
			messages = append(messages, fmt.Sprintf("%T: %s", err, err.Error()))
		}
	}

	return fmt.Sprint(strings.Join(messages, "\n"))
}

func (s StructValidatorResult) HasFieldErrorFor(field string) bool {
	for _, err := range s {
		if err.Field() == field {
			return true
		}
	}

	return false
}

func ValidateWithCustomFunc(s interface{}, validators map[string]validator.Func) StructValidatorResult {
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

		var _errors validator.ValidationErrors
		if errors.As(err, &_errors) {
			return StructValidatorResult(err.(validator.ValidationErrors))
		}
	}

	return nil
}

func SimpleStructValidation(s interface{}) StructValidatorResult {
	return ValidateWithCustomFunc(s, map[string]validator.Func{})
}

func ValidateAndWarn(s interface{}) error {
	_errors := SimpleStructValidation(s)
	if len(_errors) > 0 {
		for _, err := range _errors {
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
