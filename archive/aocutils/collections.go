package aocutils

// Value interface is a Generic wrapper around types that data is commonly used in
type Value interface {
	int | int64 | float64 | rune | string
}

// FindDuplicates will find any elements that appear more than once in the slice.
// Returns a slice containing all duplicates that appear in the input slice.
func FindDuplicates[v Value](elements []v) []v {
	found := map[v]bool{}
	dupes := []v{}
	for _, n := range elements {
		if found[n] {
			dupes = append(dupes, n)
		} else {
			found[n] = true
		}
	}

	return dupes
}

// FlattenList takes a slice matrix of values and will turn it into a 1D slice.
// Returns a 1D slice, constructed by appending each row of the matrix to one another in order.
func FlattenList[v Value](nested [][]v) []v {
	var flat []v
	for _, list := range nested {
		flat = append(flat, list...)
	}

	return flat
}

// NthElements takes a slice of elements and an interval.
// Returns a slice of the elements occuring every `n` spots of the original array.
func NthElements[v Value](elements []v, n int) []v {
	var nth []v
	for i := n - 1; i < len(elements); i += n {
		nth = append(nth, elements[i])
	}
	return nth
}

// NumOccurrences counts the number of times a specific value occurs in a slice.
// Returns the count of occurrences of the target value in the input slice.
func NumOccurrences[v Value](elements []v, target v) int {
	count := 0
	for _, elem := range elements {
		if elem == target {
			count++
		}
	}
	return count
}

// TransposeMatrix takes a 2D slice (matrix) and transposes it (swaps rows and columns).
// Returns the transposed matrix as a new 2D slice.
func TransposeMatrix[v Value](matrix [][]v) [][]v {
	if len(matrix) == 0 {
		return [][]v{}
	}
	rows := len(matrix)
	cols := len(matrix[0])
	transposed := make([][]v, cols)
	for i := range transposed {
		transposed[i] = make([]v, rows)
	}
	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			transposed[j][i] = matrix[i][j]
		}
	}
	return transposed
}

// RemoveDuplicates removes duplicate elements from a slice.
// Returns a new slice with only unique elements, preserving the order of their first occurrence.
func RemoveDuplicates[v Value](elements []v) []v {
	seen := make(map[v]bool)
	unique := []v{}
	for _, elem := range elements {
		if !seen[elem] {
			seen[elem] = true
			unique = append(unique, elem)
		}
	}
	return unique
}

// ReverseSlice reverses the order of elements in a slice.
// Returns a new slice with the elements in reverse order.
func ReverseSlice[v Value](elements []v) []v {
	reversed := make([]v, len(elements))
	for i := range elements {
		reversed[len(elements)-1-i] = elements[i]
	}
	return reversed
}
