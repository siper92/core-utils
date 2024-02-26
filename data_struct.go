package core_utils

import (
	"reflect"
)

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
