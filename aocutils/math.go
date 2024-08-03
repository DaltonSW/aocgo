package aocutils

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
