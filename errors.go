package core_utils

import (
	"fmt"
	"os"
)

func PanicOnError(err error) {
	if err != nil {
		panic(err)
	}
}

func StopOnError(err interface{}) {
	if err == nil {
		return
	}

	switch x := err.(type) {
	case error:
		PrintError(x.Error())
	default:
		PrintError("FatalError[%T]: %v", x, x)
	}

	os.Exit(1)
}

func StopOnPanic() func() {
	return func() {
		if r := recover(); r != nil {
			if err, ok := r.(error); ok {
				StopOnError(err)
			} else {
				StopOnError(
					fmt.Errorf("panic: %v", r),
				)
			}
		}
	}
}

func RecoverPanicPrint() func() {
	return func() {
		if r := recover(); r != nil {
			switch x := r.(type) {
			case error:
				PrintError(x.Error())
			default:
				PrintError("panic[%T]: %v", x, x)
			}
		}
	}
}

func RecoverPanicAsError(err *error) func() error {
	return func() error {
		if r := recover(); r != nil {
			switch x := r.(type) {
			case error:
				*err = x
			default:
				*err = fmt.Errorf("panic[%T]: %v", r, r)
			}

			if IsDebugMode() {
				StopOnError(*err)
			}
		}

		return nil
	}
}
