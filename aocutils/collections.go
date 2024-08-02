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

// TODO: RemoveDuplicates()

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
