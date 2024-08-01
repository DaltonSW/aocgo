package aocutils

type Number interface {
	int | int64 | float64
}

func AverageListOfNums[n Number](elements []n) n {
	var total n
	for _, num := range elements {
		total += num
	}

	return total / n(len(elements))
}
