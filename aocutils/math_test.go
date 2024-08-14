package aocutils

import (
	"math"
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

func TestDistance2D(t *testing.T) {
	tests := []struct {
		x1, y1, x2, y2 int
		expected       float64
	}{
		{0, 0, 3, 4, 5},
		{1, 1, 4, 5, 5},
		{-1, -1, -4, -5, 5},
	}

	for _, test := range tests {
		result := Distance2D(test.x1, test.y1, test.x2, test.y2)
		if math.Abs(result-test.expected) > 1e-9 {
			t.Errorf("Distance2D(%d, %d, %d, %d) = %f; expected %f", test.x1, test.y1, test.x2, test.y2, result, test.expected)
		}
	}
}

func TestDistance3D(t *testing.T) {
	tests := []struct {
		x1, y1, z1, x2, y2, z2 int
		expected               float64
	}{
		{0, 0, 0, 1, 1, 1, math.Sqrt(3)},
		{1, 1, 1, 4, 5, 6, math.Sqrt(50)},
		{-1, -1, -1, -4, -5, -6, math.Sqrt(50)},
	}

	for _, test := range tests {
		result := Distance3D(test.x1, test.y1, test.z1, test.x2, test.y2, test.z2)
		if math.Abs(result-test.expected) > 1e-9 {
			t.Errorf("Distance3D(%d, %d, %d, %d, %d, %d) = %f; expected %f", test.x1, test.y1, test.z1, test.x2, test.y2, test.z2, result, test.expected)
		}
	}
}

func TestManhattanDistance2D(t *testing.T) {
	tests := []struct {
		x1, y1, x2, y2 int
		expected       int
	}{
		{0, 0, 3, 4, 7},
		{1, 1, 4, 5, 7},
		{-1, -1, -4, -5, 7},
	}

	for _, test := range tests {
		result := ManhattanDistance2D(test.x1, test.y1, test.x2, test.y2)
		if result != test.expected {
			t.Errorf("ManhattanDistance2D(%d, %d, %d, %d) = %d; expected %d", test.x1, test.y1, test.x2, test.y2, result, test.expected)
		}
	}
}

func TestManhattanDistance3D(t *testing.T) {
	tests := []struct {
		x1, y1, z1, x2, y2, z2 int
		expected               int
	}{
		{0, 0, 0, 1, 1, 1, 3},
		{1, 1, 1, 4, 5, 6, 12},
		{-1, -1, -1, -4, -5, -6, 12},
	}

	for _, test := range tests {
		result := ManhattanDistance3D(test.x1, test.y1, test.z1, test.x2, test.y2, test.z2)
		if result != test.expected {
			t.Errorf("ManhattanDistance3D(%d, %d, %d, %d, %d, %d) = %d; expected %d", test.x1, test.y1, test.z1, test.x2, test.y2, test.z2, result, test.expected)
		}
	}
}
