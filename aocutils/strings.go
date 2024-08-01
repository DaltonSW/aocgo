package aocutils

import (
	"regexp"
	"strconv"
)

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
