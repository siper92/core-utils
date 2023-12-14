package core_utils

import (
	"fmt"
	"testing"
)

func Test_Add_Line_Diffs(t *testing.T) {
	expected := `line 1
line 2
line 3
line 4
line 5`

	actual := `line 1
line 2&
line 3
line 5
line 5`

	diff := getLineDifferenceInfo(expected, actual)
	if len(diff) != 2 {
		t.Errorf("Expected 2 differences, got %d", len(diff))
	}

	paddedExpected, paddedActual := addLinesNumbers(expected, actual)
	fmt.Printf("%+v", paddedExpected)
	fmt.Printf("%+v", paddedActual)
}
