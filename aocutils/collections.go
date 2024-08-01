package aocutils

type Value interface {
	int | int64 | float64 | rune | string
}

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

func FlattenList[v Value](nested [][]v) []v {
	var flat []v
	for _, list := range nested {
		flat = append(flat, list...)
	}

	return flat
}

func NthElements[v Value](elements []v, n int) []v {
	var nth []v
	for i := n - 1; i < len(elements); i += n {
		nth = append(nth, elements[i])
	}
	return nth
}
