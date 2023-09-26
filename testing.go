package core_utils

import (
	"fmt"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

type Tester struct {
	t *testing.T
}

func AMatchesB[T comparable](t *testing.T, expected T, actual T) {
	if expected != actual {
		t.Fatal(
			fmt.Sprintf("\nExpected: \n%v\n%s\n%v",
				expected,
				strings.Repeat("=", 60),
				actual,
			),
		)
	}
}

func FatalOnErr(t *testing.T, err error) {
	if err != nil {
		t.Fatal(err)
	}
}

func GetCallerDir() string {
	_, filename, _, _ := runtime.Caller(1)
	return filepath.Dir(filename)
}
