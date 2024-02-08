package core_utils

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"reflect"
)

func ValidateStruct(s interface{}) error {
	errors := ValidateGenericStruct(s)
	for _, err := range errors {
		PrintValidationError(err)
	}

	if len(errors) > 0 {
		return errors.SingleError()
	}

	return nil
}

func PrintValidationError(err validator.FieldError) {
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

func CopyStructToDefault[T interface{}](src, dst T) T {
	r := reflect.ValueOf(src)
	rResult := reflect.ValueOf(&dst)
	for i := 0; i < r.NumField(); i++ {
		field := r.Field(i)
		if !field.IsZero() {
			fieldName := r.Type().Field(i).Name
			_field := reflect.Indirect(rResult).FieldByName(fieldName)
			if _field.CanSet() {
				_field.Set(field)
			}
		}
	}
	return dst
}
