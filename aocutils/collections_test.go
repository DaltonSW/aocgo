package aocutils

import (
	"reflect"
	"testing"
)

func TestFindDuplicates(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected interface{}
	}{
		{
			name:     "int duplicates",
			input:    []int{1, 2, 3, 4, 2, 3, 1},
			expected: []int{2, 3, 1},
		},
		{
			name:     "int64 duplicates",
			input:    []int64{1, 2, 3, 4, 2, 3, 1},
			expected: []int64{2, 3, 1},
		},
		{
			name:     "float64 duplicates",
			input:    []float64{1.1, 2.2, 3.3, 4.4, 2.2, 3.3, 1.1},
			expected: []float64{2.2, 3.3, 1.1},
		},
		{
			name:     "rune duplicates",
			input:    []rune{'a', 'b', 'c', 'd', 'b', 'c', 'a'},
			expected: []rune{'b', 'c', 'a'},
		},
		{
			name:     "string duplicates",
			input:    []string{"apple", "banana", "cherry", "apple", "banana"},
			expected: []string{"apple", "banana"},
		},
		{
			name:     "no duplicates",
			input:    []int{1, 2, 3, 4, 5},
			expected: []int{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch input := tt.input.(type) {
			case []int:
				result := FindDuplicates(input)
				expected := tt.expected.([]int)
				if !reflect.DeepEqual(result, expected) {
					t.Errorf("expected %v, got %v", expected, result)
				}
			case []int64:
				result := FindDuplicates(input)
				expected := tt.expected.([]int64)
				if !reflect.DeepEqual(result, expected) {
					t.Errorf("expected %v, got %v", expected, result)
				}
			case []float64:
				result := FindDuplicates(input)
				expected := tt.expected.([]float64)
				if !reflect.DeepEqual(result, expected) {
					t.Errorf("expected %v, got %v", expected, result)
				}
			case []rune:
				result := FindDuplicates(input)
				expected := tt.expected.([]rune)
				if !reflect.DeepEqual(result, expected) {
					t.Errorf("expected %v, got %v", expected, result)
				}
			case []string:
				result := FindDuplicates(input)
				expected := tt.expected.([]string)
				if !reflect.DeepEqual(result, expected) {
					t.Errorf("expected %v, got %v", expected, result)
				}
			}
		})
	}
}

func TestFlattenList(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected interface{}
	}{
		{
			name:     "int slice matrix",
			input:    [][]int{{1, 2}, {3, 4}, {5}},
			expected: []int{1, 2, 3, 4, 5},
		},
		{
			name:     "string slice matrix",
			input:    [][]string{{"a", "b"}, {"c", "d"}, {"e"}},
			expected: []string{"a", "b", "c", "d", "e"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch input := tt.input.(type) {
			case [][]int:
				result := FlattenList(input)
				expected := tt.expected.([]int)
				if !reflect.DeepEqual(result, expected) {
					t.Errorf("expected %v, got %v", expected, result)
				}
			case [][]string:
				result := FlattenList(input)
				expected := tt.expected.([]string)
				if !reflect.DeepEqual(result, expected) {
					t.Errorf("expected %v, got %v", expected, result)
				}
			}
		})
	}
}

func TestNthElements(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		n        int
		expected interface{}
	}{
		{
			name:     "int slice every 2nd element",
			input:    []int{1, 2, 3, 4, 5, 6},
			n:        2,
			expected: []int{2, 4, 6},
		},
		{
			name:     "string slice every 3rd element",
			input:    []string{"a", "b", "c", "d", "e", "f", "g"},
			n:        3,
			expected: []string{"c", "f"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch input := tt.input.(type) {
			case []int:
				result := NthElements(input, tt.n)
				expected := tt.expected.([]int)
				if !reflect.DeepEqual(result, expected) {
					t.Errorf("expected %v, got %v", expected, result)
				}
			case []string:
				result := NthElements(input, tt.n)
				expected := tt.expected.([]string)
				if !reflect.DeepEqual(result, expected) {
					t.Errorf("expected %v, got %v", expected, result)
				}
			}
		})
	}
}

func TestNumOccurrences(t *testing.T) {
	tests := []struct {
		input    []int
		target   int
		expected int
	}{
		{[]int{1, 2, 3, 4, 2, 1, 2}, 2, 3},
		{[]int{1, 1, 1, 1}, 1, 4},
		{[]int{5, 6, 7, 8}, 9, 0},
	}

	for _, test := range tests {
		result := NumOccurrences(test.input, test.target)
		if result != test.expected {
			t.Errorf("NumOccurrences(%v, %d) = %d; expected %d", test.input, test.target, result, test.expected)
		}
	}
}

func TestTransposeMatrix(t *testing.T) {
	tests := []struct {
		input    [][]int
		expected [][]int
	}{
		{
			[][]int{
				{1, 2, 3},
				{4, 5, 6},
				{7, 8, 9},
			},
			[][]int{
				{1, 4, 7},
				{2, 5, 8},
				{3, 6, 9},
			},
		},
		{
			[][]int{
				{1, 2},
				{3, 4},
			},
			[][]int{
				{1, 3},
				{2, 4},
			},
		},
	}

	for _, test := range tests {
		result := TransposeMatrix(test.input)
		if !reflect.DeepEqual(result, test.expected) {
			t.Errorf("TransposeMatrix(%v) = %v; expected %v", test.input, result, test.expected)
		}
	}
}

func TestRemoveDuplicates(t *testing.T) {
	tests := []struct {
		input    []int
		expected []int
	}{
		{[]int{1, 2, 3, 4, 2, 1, 2}, []int{1, 2, 3, 4}},
		{[]int{5, 5, 5, 5}, []int{5}},
		{[]int{1, 2, 3, 4}, []int{1, 2, 3, 4}},
	}

	for _, test := range tests {
		result := RemoveDuplicates(test.input)
		if !reflect.DeepEqual(result, test.expected) {
			t.Errorf("RemoveDuplicates(%v) = %v; expected %v", test.input, result, test.expected)
		}
	}
}

func TestReverseSlice(t *testing.T) {
	tests := []struct {
		input    []int
		expected []int
	}{
		{[]int{1, 2, 3, 4}, []int{4, 3, 2, 1}},
		{[]int{5, 6, 7, 8}, []int{8, 7, 6, 5}},
		{[]int{9}, []int{9}},
	}

	for _, test := range tests {
		result := ReverseSlice(test.input)
		if !reflect.DeepEqual(result, test.expected) {
			t.Errorf("ReverseSlice(%v) = %v; expected %v", test.input, result, test.expected)
		}
	}
}
