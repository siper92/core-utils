package core_utils

import (
	"bytes"
	"fmt"
	"github.com/siper92/core-utils/type_utils"
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

func _addLinesNumbers(lines []string, diffs []int) []string {
	formattedLines := make([]string, len(lines))

	linesCount := len(lines)
	padding := len(fmt.Sprintf("%d", linesCount)) + 1

	for i, line := range lines {
		linePadding := padding - len(fmt.Sprintf("%d", i+1))
		isDiff := " "
		if type_utils.InArray(i+1, diffs) {
			isDiff = "*"
		}

		formattedLines[i] = fmt.Sprintf("|%s%d%s|%s", strings.Repeat(" ", linePadding), i+1, isDiff, line)
	}

	return formattedLines
}

func addLinesNumbers[T comparable](expected, actual T) (string, string) {
	paddedExpected := fmt.Sprintf("%v", expected)
	paddedActual := fmt.Sprintf("%v", actual)

	diff := getLineDifferenceInfo(expected, actual)
	if len(diff) > 0 {
		expectedLines := getLines(expected)
		actualLines := getLines(actual)

		paddedExpected = strings.Join(_addLinesNumbers(expectedLines, diff), "\n")
		paddedActual = strings.Join(_addLinesNumbers(actualLines, diff), "\n")
	}

	return paddedExpected, paddedActual
}

func AMatchesBDetailed[T comparable](t *testing.T, expected T, actual T) {
	if expected != actual {
		paddedExpected, paddedActual := addLinesNumbers(expected, actual)

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
