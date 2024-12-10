package aocutils

import "math"

// Number interface is a Generic wrapper around numeric types that data is commonly used in
type Number interface {
	int | int64 | float64
}

// AverageListOfNums takes a list of Number elements and returns the average of all of them.
func AverageListOfNums[n Number](elements []n) n {
	if len(elements) == 0 {
		return 0
	}

	var total n
	for _, num := range elements {
		total += num
	}

	return total / n(len(elements))
}

// Distance2D calculates the Euclidean distance between two points in 2D space.
// Returns the distance as a float64.
func Distance2D[n Number](x1, y1, x2, y2 n) float64 {
	return math.Sqrt(float64((x2-x1)*(x2-x1) + (y2-y1)*(y2-y1)))
}

// Slope2D calculates the slope between two points in 2D space.
// Returns the slope as a float64
func Slope2D[n Number](x1, y1, x2, y2 n) float64 {
	return float64((y2 - y1) / (x2 - x1))
}

// Distance3D calculates the Euclidean distance between two points in 3D space.
// Returns the distance as a float64.
func Distance3D[n Number](x1, y1, z1, x2, y2, z2 n) float64 {
	return math.Sqrt(float64((x2-x1)*(x2-x1) + (y2-y1)*(y2-y1) + (z2-z1)*(z2-z1)))
}

// ManhattanDistance2D calculates the Manhattan distance between two points in 2D space.
// Returns the distance as an int.
func ManhattanDistance2D(x1, y1, x2, y2 int) int {
	return AbsVal(x2-x1) + AbsVal(y2-y1)
}

// ManhattanDistance3D calculates the Manhattan distance between two points in 3D space.
// Returns the distance as an int.
func ManhattanDistance3D(x1, y1, z1, x2, y2, z2 int) int {
	return AbsVal(x2-x1) + AbsVal(y2-y1) + AbsVal(z2-z1)
}

// AbsVal calculates the absolute value of a number.
func AbsVal[num Number](n num) num {
	if n < 0 {
		return -n
	}
	return n
}
