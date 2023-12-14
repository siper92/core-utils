package core_utils

import (
	"bytes"
	"fmt"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

type Tester struct {
	t *testing.T
}

func getLineDifferenceInfo[T comparable](expected T, actual T) string {
	if expected == actual {
		return ""
	}

	expectedStr := fmt.Sprintf("%v", expected)
	actualStr := fmt.Sprintf("%v", actual)

	expectedLines := strings.Split(expectedStr, "\n")
	actualLines := strings.Split(actualStr, "\n")
	if len(expectedLines) == 1 || len(actualLines) == 1 {
		return "Different values"
	}

	result := bytes.NewBufferString("")

	longest := len(expectedLines)
	if len(actualLines) > longest {
		longest = len(actualLines)
		result.WriteString(
			fmt.Sprintf("More lines than expected (%d vs %d)\n", len(expectedLines), len(actualLines)),
		)
	} else if len(expectedLines) > len(actualLines) {
		result.WriteString(
			fmt.Sprintf("Less lines than expected (%d vs %d)\n", len(expectedLines), len(actualLines)),
		)
	}

	var diffSummary []string
	for index, line := range expectedLines {
		if len(actualLines) < index {
			continue
		}

		if line != actualLines[index] {
			diffSummary = append(diffSummary, fmt.Sprintf("%d", index+1))
		}
	}

	if len(diffSummary) > 0 {
		result.WriteString("Diff lines in actual: " + strings.Join(diffSummary, ", "))
	}

	return result.String()
}

func AMatchesB[T comparable](t *testing.T, expected T, actual T) {
	if expected != actual {
		t.Fatal(
			fmt.Sprintf("\nExpected: \n%v\n%s\n%v\n%s",
				expected,
				strings.Repeat("=", 60),
				actual,
				getLineDifferenceInfo(expected, actual),
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
