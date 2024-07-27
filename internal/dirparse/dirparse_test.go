package dirparse

import (
	"fmt"
	"testing"
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
			outDay, outYear, outErr := GetDayAndYearFromDirInput(test.cur, test.paren)
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
