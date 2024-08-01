package aocgo

func FindDuplicateInts(elements []int) []int {
	found := map[int]bool{}
	dupes := []int{}
	for _, n := range elements {
		if found[n] {
			dupes = append(dupes, n)
		} else {
			found[n] = true
		}
	}

	return dupes
}

func FindDuplicateStrings(elements []string) []string {
	found := map[string]bool{}
	dupes := []string{}
	for _, n := range elements {
		if found[n] {
			dupes = append(dupes, n)
		} else {
			found[n] = true
		}
	}

	return dupes
}

func FlattenIntList(nested [][]int) []int {
	var flat []int
	for _, list := range nested {
		flat = append(flat, list...)
	}

	return flat
}

func FlattenStringList(nested [][]string) []string {
	var flat []string
	for _, list := range nested {
		flat = append(flat, list...)
	}

	return flat
}
