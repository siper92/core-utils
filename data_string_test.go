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
		{"nil *int", (*int)(nil), "0"},
		{"non-nil *int", func() *int { i := 42; return &i }(), "42"},
		{"nil *int32", (*int32)(nil), "0"},
		{"non-nil *int32", func() *int32 { var i int32 = 32; return &i }(), "32"},
		{"nil *int64", (*int64)(nil), "0"},
		{"non-nil *int64", func() *int64 { var i int64 = 64; return &i }(), "64"},
		{"nil *float32", (*float32)(nil), "0"},
		{"non-nil *float32", func() *float32 { f := float32(3.14); return &f }(), "3.140000"},
		{"nil *float64", (*float64)(nil), "0"},
		{"non-nil *float64", func() *float64 { f := 6.28; return &f }(), "6.280000"},
		{"nil *bool", (*bool)(nil), ""},
		{"non-nil *bool", func() *bool { b := true; return &b }(), "true"},
		{"nil *string", (*string)(nil), ""},
		{"non-nil *string", func() *string { s := "hello"; return &s }(), "hello"},
		{"nil *[]byte", (*[]byte)(nil), ""},
		{"non-nil *[]byte", func() *[]byte { b := []byte("byte slice"); return &b }(), "byte slice"},
		{"string", "hello", "hello"},
		{"int", 42, "42"},
		{"int32", int32(32), "32"},
		{"int64", int64(64), "64"},
		{"float32", float32(3.14), "3.140000"},
		{"float64", 6.28, "6.280000"},
		{"bool", true, "true"},
		{"bytes", []byte("byte slice"), "byte slice"},
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
