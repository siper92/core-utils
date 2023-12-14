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

func getLines(v interface{}) []string {
	str := fmt.Sprintf("%v", v)
	return strings.Split(str, "\n")
}
func getLineDifferenceInfo[T comparable](expected T, actual T) []int {
	var diffSummary []int
	if expected == actual {
		return diffSummary
	}

	expectedLines := getLines(expected)
	actualLines := getLines(actual)
	if len(expectedLines) == 1 || len(actualLines) == 1 {
		return diffSummary
	}

	for index, line := range expectedLines {
		if len(actualLines) <= index {
			continue
		}

		if line != actualLines[index] {
			diffSummary = append(diffSummary, index+1)
		}
	}

	return diffSummary
}

func newMismatchError[T comparable](expected T, actual T) string {
	expectedLines := getLines(expected)
	actualLines := getLines(actual)

	details := bytes.NewBufferString("")

	if len(actualLines) > len(expectedLines) {
		details.WriteString(
			fmt.Sprintf("More lines than expected (%d vs %d)\n", len(expectedLines), len(actualLines)),
		)
	} else if len(expectedLines) > len(actualLines) {
		details.WriteString(
			fmt.Sprintf("Less lines than expected (%d vs %d)\n", len(expectedLines), len(actualLines)),
		)
	} else if len(actualLines) == 1 || len(expectedLines) == 1 {
		details.WriteString("Invalid values\n")
	}

	diff := getLineDifferenceInfo(expected, actual)
	if len(diff) > 0 {
		details.WriteString("Lines with differences: ")
		for _, line := range diff {
			details.WriteString(fmt.Sprintf("%d, ", line))
		}
		details.WriteString("\n")
	}

	return details.String()
}

func AMatchesBDetailed[T comparable](t *testing.T, expected T, actual T) {
	if expected != actual {
		paddedExpected := fmt.Sprintf("%v", expected)
		paddedActual := fmt.Sprintf("%v", actual)

		diff := getLineDifferenceInfo(expected, actual)
		if len(diff) > 0 {
			expectedLines := getLines(expected)
			actualLines := getLines(actual)
			formattedExpectedLines := make([]string, len(expectedLines))
			for i, line := range expectedLines {
				formattedExpectedLines[i] = fmt.Sprintf("[ %d ]: %s", i+1, line)
			}
			formattedActualLines := make([]string, len(actualLines))
			for i, line := range actualLines {
				formattedActualLines[i] = fmt.Sprintf("[ %d ]: %s", i+1, line)
			}

			paddedExpected = strings.Join(formattedExpectedLines, "\n")
			paddedActual = strings.Join(formattedActualLines, "\n")
		}
		t.Fatal(
			fmt.Sprintf("\nExpected: \n%s\n%s\n%s\n\n%s",
				paddedExpected,
				strings.Repeat("=", 60),
				paddedActual,
				newMismatchError(expected, actual),
			),
		)
	}
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
