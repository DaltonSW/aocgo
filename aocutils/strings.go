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

// FindSubstringsOfLength takes a string and a desired length, and returns all unique substrings of that length.
// Returns a slice of substrings.
func FindSubstringsOfLength(s string, length int) []string {
	if length <= 0 || length > len(s) {
		return []string{}
	}

	substrings := make(map[string]bool)
	for i := 0; i <= len(s)-length; i++ {
		substr := s[i : i+length]
		substrings[substr] = true
	}

	uniqueSubstrings := make([]string, 0, len(substrings))
	for substr := range substrings {
		uniqueSubstrings = append(uniqueSubstrings, substr)
	}

	return uniqueSubstrings
}
