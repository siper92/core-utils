package app_errors

import (
	"fmt"
	core_utils "github.com/siper92/core-utils"
)

func PanicOnError(err error) {
	if err != nil {
		panic(err)
	}
}

func StopOnPanicHandler() func() {
	return func() {
		if r := recover(); r != nil {
			if err, ok := r.(error); ok {
				core_utils.StopOnError(err)
			} else {
				core_utils.StopOnError(
					fmt.Errorf("panic: %v", r),
				)
			}
		}
	}
}
