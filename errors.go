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

	if v, ok := err.(error); ok {
		PrintError(v.Error())
	}

	os.Exit(1)
}

func StopOnPanicHandler() func() {
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

func PrintPanicHandler() func() {
	return func() {
		if r := recover(); r != nil {
			if err, ok := r.(error); ok {
				PrintError(err.Error())
			} else {
				PrintError(
					fmt.Sprintf("panic: %v", r),
				)
			}
		}
	}
}
