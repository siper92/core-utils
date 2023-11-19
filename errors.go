package core_utils

import (
	"errors"
	"fmt"
)

func NewError(params ...interface{}) error {
	return errors.New(fmt.Sprintf(params[0].(string), params[1:]...))
}

func HandleWarning(_ any, err error) {
	if err != nil {
		fmt.Println(err.Error())
	}
}
