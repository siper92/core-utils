package core_utils

import (
	"errors"
	"fmt"
)

func NewError(params ...interface{}) error {
	return errors.New(fmt.Sprintf(params[0].(string), params[1:]...))
}
