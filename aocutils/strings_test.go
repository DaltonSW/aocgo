package aocutils

import (
	"testing"
)

func TestExtractIntsFromString(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []int
	}{
		{
			name:     "extract integers from string with positive numbers",
			input:    "There are 2 apples, 3 oranges, and 10 bananas.",
			expected: []int{2, 3, 10},
		},
		{
			name:     "extract integers from string with negative numbers",
			input:    "The temperature was -5 degrees, then it dropped to -10.",
			expected: []int{-5, -10},
		},
		{
			name:     "extract integers from string with mixed numbers",
			input:    "Height: 6ft, Depth: -2ft, Width: 3ft",
			expected: []int{6, -2, 3},
		},
		{
			name:     "extract integers from string with no numbers",
			input:    "There are no numbers in this sentence.",
			expected: []int{},
		},
		{
			name:     "extract integers from empty string",
			input:    "",
			expected: []int{},
		},
		{
			name:     "extract integers from string with multiple digit numbers",
			input:    "The codes are 123, 4567, and 890.",
			expected: []int{123, 4567, 890},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ExtractIntsFromString(tt.input)
			if len(result) != len(tt.expected) {
				t.Errorf("expected length %d, got %d", len(tt.expected), len(result))
				return
			}
			for i := range result {
				if result[i] != tt.expected[i] {
					t.Errorf("expected %v, got %v", tt.expected, result)
					break
				}
			}
		})
	}
}
