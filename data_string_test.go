package core_utils

import "testing"

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
