package aocutils

import (
	"reflect"
	"testing"
)

func TestAverageListOfNums(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected interface{}
	}{
		{
			name:     "average of int slice",
			input:    []int{1, 2, 3, 4, 5},
			expected: 3,
		},
		{
			name:     "average of int64 slice",
			input:    []int64{1, 2, 3, 4, 5},
			expected: int64(3),
		},
		{
			name:     "average of float64 slice",
			input:    []float64{1.0, 2.0, 3.0, 4.0, 5.0},
			expected: 3.0,
		},
		{
			name:     "average of single element slice",
			input:    []int{42},
			expected: 42,
		},
		{
			name:     "average of empty slice",
			input:    []int{},
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch input := tt.input.(type) {
			case []int:
				var expected int
				if len(input) > 0 {
					expected = tt.expected.(int)
				} else {
					expected = 0 // handle empty slice case for int
				}
				result := AverageListOfNums(input)
				if !reflect.DeepEqual(result, expected) {
					t.Errorf("expected %v, got %v", expected, result)
				}
			case []int64:
				var expected int64
				if len(input) > 0 {
					expected = tt.expected.(int64)
				} else {
					expected = 0 // handle empty slice case for int64
				}
				result := AverageListOfNums(input)
				if !reflect.DeepEqual(result, expected) {
					t.Errorf("expected %v, got %v", expected, result)
				}
			case []float64:
				var expected float64
				if len(input) > 0 {
					expected = tt.expected.(float64)
				} else {
					expected = 0.0 // handle empty slice case for float64
				}
				result := AverageListOfNums(input)
				if !reflect.DeepEqual(result, expected) {
					t.Errorf("expected %v, got %v", expected, result)
				}
			}
		})
	}
}
