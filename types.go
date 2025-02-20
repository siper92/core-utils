package core_utils

import "reflect"

func IsNil(i interface{}) bool {
	if i == nil {
		return true
	}
	value := reflect.ValueOf(i)
	switch value.Kind() {
	case reflect.Ptr, reflect.Map, reflect.Array, reflect.Chan, reflect.Slice:
		return value.IsNil()
	default:
		return false
	}
}
