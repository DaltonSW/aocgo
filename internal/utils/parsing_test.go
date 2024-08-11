package utils

import (
	"fmt"
	"testing"
	"time"
)

func TestGetDayAndYearFromCWD(t *testing.T) {
	var tests = []struct {
		paren, cur string
		year, day  int
	}{
		{"2016", "1", 2016, 1},
		{"17", "10", 2017, 10},
		{"23", "03", 2023, 3},
		{"'19", "15", 2019, 15},
		{"Year 2020", "Day 09", 2020, 9},
	}

	for _, test := range tests {
		testname := fmt.Sprintf(".../%v/%v", test.paren, test.cur)
		t.Run(testname, func(t *testing.T) {
			outYear, outDay, outErr := GetYearAndDayFromDirInput(test.cur, test.paren)
			if outErr != nil {
				t.Errorf("Got unexpected error: %v", outErr)
			}

			if outDay != test.day {
				t.Errorf("Day mismatch. Got %d, expected %d", outDay, test.day)
			}

			if outYear != test.year {
				t.Errorf("Year mismatch. Got %d, expected %d", outYear, test.year)
			}
		})
	}
}

// TestParseDuration is the unit test for the ParseDuration function
func TestParseDuration(t *testing.T) {
	testCases := []struct {
		input    string
		expected time.Duration
	}{
		{"Please wait one minute before trying again.", time.Minute},
		{"You have 1m 58s left to wait", time.Minute + 58*time.Second},
		{"You need to wait 3 minutes and 30 seconds.", 3*time.Minute + 30*time.Second},
		{"Please try again in 2m", 2 * time.Minute},
		{"Wait for 45 seconds before retrying.", 45 * time.Second},
		{"Wait a second, I'm thinking...", time.Second},
		{"Please give it one second.", time.Second},
		{"2 minutes 45 seconds remaining.", 2*time.Minute + 45*time.Second},
		{"One minute left!", time.Minute},
	}

	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			duration, err := ParseDuration(tc.input)
			if err != nil {
				t.Errorf("Error parsing duration: %v", err)
			}
			if duration != tc.expected {
				t.Errorf("Expected %v, got %v", tc.expected, duration)
			}
		})
	}
}
