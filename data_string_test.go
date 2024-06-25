package core_utils

import (
	"fmt"
	"testing"
)

func Test_SnakeCaseToCamelCase(t *testing.T) {
	tests := []struct {
		name     string
		snake    string
		expected string
	}{
		{
			name:     "empty string",
			snake:    "",
			expected: "",
		},
		{
			name:     "single word",
			snake:    "word",
			expected: "word",
		},
		{
			name:     "two words",
			snake:    "two_words",
			expected: "twoWords",
		},
		{
			name:     "three words",
			snake:    "three_words_here",
			expected: "threeWordsHere",
		},
		{
			name:     "prefix underscore",
			snake:    "_prefix_under_score",
			expected: "prefixUnderScore",
		},
		{
			name:     "suffix underscore",
			snake:    "suffix_under_score_",
			expected: "suffixUnderScore",
		},
		{
			name:     "has upper characters",
			snake:    "has_Upper_cHaracters",
			expected: "hasUpperCHaracters",
		},
		{
			name:     "stars with upper",
			snake:    "Starts_with_upper",
			expected: "startsWithUpper",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual := SnakeCaseToCamelCase(test.snake)
			if actual != test.expected {
				t.Errorf("expected %s, got %s", test.expected, actual)
			}
		})
	}
}

// Mock for fmt.Stringer
type Person struct {
	Name string
}

func (p Person) String() string {
	return fmt.Sprintf("Person Name: %s", p.Name)
}

func TestToStringV2(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected string
	}{
		{"nil pointer", (*int)(nil), ""},
		{"string", "hello", "hello"},
		{"*string", new(string), ""},
		{"int", 42, "42"},
		{"*int", new(int), "0"},
		{"int32", int32(32), "32"},
		{"*int32", new(int32), "0"},
		{"int64", int64(64), "64"},
		{"*int64", new(int64), "0"},
		{"float32", float32(3.14), "3.140000"},
		{"*float32", new(float32), "0.000000"},
		{"float64", 6.28, "6.280000"},
		{"*float64", new(float64), "0.000000"},
		{"bool", true, "true"},
		{"*bool", new(bool), "false"},
		{"bytes", []byte("byte slice"), "byte slice"},
		{"*[]byte", new([]byte), ""},
		{"fmt.Stringer", Person{Name: "Onyx"}, "Person Name: Onyx"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ToString(tt.input)
			if result != tt.expected {
				t.Errorf("ToStringV2(%v) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}
