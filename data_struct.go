package core_utils

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"reflect"
)

func ValidateStruct(s interface{}) error {
	err := validator.New().Struct(s)
	if err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			PrintError("Invalid Validation Error: %s", err.Error())
			return nil
		}

		for _, err := range err.(validator.ValidationErrors) {
			PrintValidationError(err)
		}

		Warning(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>")
		Warning("Config Parsing Error")
	}

	return err
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
