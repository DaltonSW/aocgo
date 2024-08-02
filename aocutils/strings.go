package aocutils

import (
	"regexp"
	"strconv"
)

// ExtractIntsFromString will find all groups of consecutive digits in a string.
// Returns a slice of all integers extracted from the input string.
func ExtractIntsFromString(input string) []int {
	re := regexp.MustCompile(`-?\d+`)
	matches := re.FindAllString(input, -1)
	var numbers []int
	for _, n := range matches {
		num, _ := strconv.Atoi(n)
		numbers = append(numbers, num)
	}

	return numbers
}
